package main

// @title Next-Go-PG API
// @version 1.0
// @description Go Clean Architecture API with Better Auth integration
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"

	"github.com/atilladeniz/next-go-pg/backend/internal/domain"
	"github.com/atilladeniz/next-go-pg/backend/internal/handler"
	"github.com/atilladeniz/next-go-pg/backend/internal/jobs"
	"github.com/atilladeniz/next-go-pg/backend/internal/middleware"
	"github.com/atilladeniz/next-go-pg/backend/internal/repository"
	"github.com/atilladeniz/next-go-pg/backend/internal/sse"
	"github.com/atilladeniz/next-go-pg/backend/pkg/config"
	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
	"github.com/atilladeniz/next-go-pg/backend/pkg/metrics"
	riverPkg "github.com/atilladeniz/next-go-pg/backend/pkg/river"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	sseBroker     *sse.Broker
	riverJobQueue *riverPkg.Client
	exportStore   *jobs.ExportStore
)

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
}

var (
	// Build information (set by build flags)
	Version   = "dev"
	BuildTime = "unknown"
	db        *gorm.DB
	pgxPool   *pgxpool.Pool
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Validate configuration on startup
	// In production: fatal on missing required vars
	// In development: log warnings for common misconfigurations
	if cfg.Environment == "production" {
		cfg.MustValidate()
	} else if err := cfg.ValidateWithWarnings(); err != nil {
		log.Printf("Configuration warnings: %v", err)
	}

	// Initialize structured logger with optional Loki integration
	logger.Init(logger.Config{
		Level:        cfg.LogLevel,
		Environment:  cfg.Environment,
		ServiceName:  "next-go-pg-api",
		Version:      Version,
		AnonymizeIPs: cfg.Logging.AnonymizeIPs,
		WithCaller:   cfg.Logging.WithCaller,
		LokiURL:      os.Getenv("LOKI_URL"), // e.g., http://localhost:3100/loki/api/v1/push
	})
	defer logger.Close()

	logger.Info().
		Str("version", Version).
		Str("build_time", BuildTime).
		Str("environment", cfg.Environment).
		Msg("Starting application")

	// Initialize Prometheus metrics
	metrics.Init(Version, cfg.Environment)

	// Connect to database with retry
	var err error
	db, err = connectToDatabase(cfg)
	if err != nil {
		logger.Warn().
			Err(err).
			Msg("Database connection failed - server will start in degraded mode")
		db = nil // Ensure db is nil for health checks
	} else {
		logger.Info().Msg("Database connected successfully")

		// Run auto-migrations if database is connected
		if err := runAutoMigrations(db); err != nil {
			logger.Warn().
				Err(err).
				Msg("Auto-migration failed - you may need to run migrations manually")
		} else {
			logger.Info().Msg("Database schema is up to date")
		}
	}

	// Setup SSE broker
	sseBroker = sse.NewBroker()
	logger.Info().Msg("SSE broker initialized")

	// Setup repositories early (needed by workers)
	var statsRepo *repository.UserStatsRepository
	if db != nil {
		statsRepo = repository.NewUserStatsRepository(db)
	}

	// Setup River job queue (if database is available)
	var riverCtx context.Context
	var riverCancel context.CancelFunc
	if db != nil {
		riverCtx, riverCancel = context.WithCancel(context.Background())
		defer func() {
			if riverCancel != nil {
				riverCancel()
			}
		}()

		// Create pgx pool for River
		var err error
		pgxPool, err = pgxpool.New(riverCtx, cfg.GetDatabaseURLForPgx())
		if err != nil {
			logger.Warn().Err(err).Msg("Failed to create pgx pool for River - background jobs disabled")
		} else {
			// Run River migrations
			if err := riverPkg.RunMigrations(riverCtx, pgxPool); err != nil {
				logger.Warn().Err(err).Msg("River migrations failed - background jobs may not work")
			}

			// Setup email configuration for job workers
			smtpHost := os.Getenv("SMTP_HOST")
			if smtpHost == "" {
				smtpHost = "127.0.0.1"
			}
			smtpPort := 1025
			if p := os.Getenv("SMTP_PORT"); p != "" {
				if parsed, err := strconv.Atoi(p); err == nil {
					smtpPort = parsed
				}
			}
			smtpFrom := os.Getenv("SMTP_FROM")
			if smtpFrom == "" {
				smtpFrom = "noreply@localhost"
			}
			appURL := os.Getenv("NEXT_PUBLIC_APP_URL")
			if appURL == "" {
				appURL = "http://localhost:3000"
			}

			emailConfig := jobs.NewEmailConfig(smtpHost, smtpPort, smtpFrom, appURL)

			// Create export store for download handling
			exportStore = jobs.NewExportStore()

			// Register workers with all dependencies
			workers := river.NewWorkers()
			jobs.RegisterWorkers(workers, &jobs.WorkerDeps{
				EmailConfig: emailConfig,
				SSEBroker:   sseBroker,
				ExportStore: exportStore,
				StatsRepo:   statsRepo,
			})

			// Create River client
			riverJobQueue, err = riverPkg.NewClient(riverCtx, pgxPool, workers, riverPkg.DefaultConfig())
			if err != nil {
				logger.Warn().Err(err).Msg("Failed to create River client - background jobs disabled")
			} else {
				// Start River
				if err := riverJobQueue.Start(riverCtx); err != nil {
					logger.Error().Err(err).Msg("Failed to start River client")
				} else {
					logger.Info().Msg("River job queue initialized and started")
				}
			}
		}
	}

	// Setup router
	router := mux.NewRouter()

	// Setup middlewares
	loggingMiddleware := middleware.NewLoggingMiddleware()
	corsMiddleware := middleware.NewCORSMiddleware(cfg.FrontendURL)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(middleware.RateLimitConfig{
		RequestsPerMinute: cfg.RateLimit.RequestsPerMinute,
		BurstSize:         cfg.RateLimit.BurstSize,
		SkipPaths:         []string{"/health", "/health/ready", "/health/live", "/metrics"},
	})
	metricsMiddleware := middleware.NewMetricsMiddleware()

	// Apply middlewares (order matters: metrics first, then logging, then rate limit)
	router.Use(metricsMiddleware.Handler)
	router.Use(loggingMiddleware.Handler)
	router.Use(corsMiddleware.Handler)
	router.Use(rateLimitMiddleware.Handler)

	// Health check endpoint with comprehensive checks
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")
	router.HandleFunc("/health/ready", readinessHandler).Methods("GET")
	router.HandleFunc("/health/live", livenessHandler).Methods("GET")

	// Prometheus metrics endpoint
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	// Setup API handlers
	apiHandler := handler.NewAPIHandler(cfg.FrontendURL, sseBroker, statsRepo)

	// Setup combined auth middleware (JWT first, then Better Auth fallback)
	// This allows Go to validate tokens directly without calling Next.js
	combinedAuth := middleware.NewCombinedAuthMiddleware(cfg.FrontendURL)

	// Setup webhook handler
	// Setup webhook handler with optional background job support
	webhookHandler := handler.NewWebhookHandler(db)
	if riverJobQueue != nil {
		webhookHandler = webhookHandler.WithJobEnqueuer(riverJobQueue.Client)
	}

	// API v1 routes
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// Public routes
	apiRouter.HandleFunc("/hello", apiHandler.PublicHello).Methods("GET")

	// Protected routes (require auth) - uses combined JWT + Better Auth
	protectedRouter := apiRouter.PathPrefix("/protected").Subrouter()
	protectedRouter.Use(combinedAuth.RequireAuth)
	protectedRouter.HandleFunc("/hello", apiHandler.ProtectedHello).Methods("GET")

	// User routes (require auth) - uses combined JWT + Better Auth
	apiRouter.Handle("/me", combinedAuth.RequireAuth(http.HandlerFunc(apiHandler.GetCurrentUser))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/stats", combinedAuth.RequireAuth(http.HandlerFunc(apiHandler.GetUserStats))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/stats", combinedAuth.RequireAuth(http.HandlerFunc(apiHandler.UpdateUserStats))).Methods("POST", "OPTIONS")

	// SSE endpoint for real-time updates
	apiRouter.Handle("/events", sseBroker).Methods("GET")

	// Test endpoint to trigger stats update (for demo)
	apiRouter.HandleFunc("/trigger-update", func(w http.ResponseWriter, r *http.Request) {
		sseBroker.Broadcast("stats-updated", `{"trigger":"manual"}`)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "broadcast sent"})
	}).Methods("POST")

	// Webhook routes (internal, protected by secret)
	webhookRouter := apiRouter.PathPrefix("/webhooks").Subrouter()
	webhookRouter.HandleFunc("/session-created", webhookHandler.SessionCreated).Methods("POST")
	webhookRouter.HandleFunc("/send-magic-link", webhookHandler.SendMagicLink).Methods("POST")
	webhookRouter.HandleFunc("/send-verification-email", webhookHandler.SendVerificationEmail).Methods("POST")
	webhookRouter.HandleFunc("/send-2fa-otp", webhookHandler.Send2FAOTP).Methods("POST")
	webhookRouter.HandleFunc("/send-2fa-enabled", webhookHandler.Send2FAEnabledNotification).Methods("POST")
	webhookRouter.HandleFunc("/send-passkey-added", webhookHandler.SendPasskeyAddedNotification).Methods("POST")

	// Export routes (require auth)
	if riverJobQueue != nil && exportStore != nil {
		exportHandler := handler.NewExportHandler(riverJobQueue.Client, exportStore)
		apiRouter.Handle("/export/start", combinedAuth.RequireAuth(http.HandlerFunc(exportHandler.StartExport))).Methods("POST", "OPTIONS")
		apiRouter.HandleFunc("/export/download/{id}", exportHandler.DownloadExport).Methods("GET")
	}

	// Setup HTTP server with timeouts
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		logger.Info().
			Str("port", cfg.Port).
			Str("frontend_url", cfg.FrontendURL).
			Msg("Server starting")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server startup failed")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop River job queue first
	if riverJobQueue != nil {
		logger.Info().Msg("Stopping River job queue...")
		if err := riverJobQueue.Stop(ctx); err != nil {
			logger.Error().Err(err).Msg("River job queue shutdown error")
		}
	}

	// Close pgx pool
	if pgxPool != nil {
		pgxPool.Close()
	}

	if err := server.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server exited")
}

func connectToDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDatabaseURL()

	logger.Info().
		Str("host", cfg.Database.Host).
		Str("port", cfg.Database.Port).
		Str("database", cfg.Database.Name).
		Msg("Connecting to database")

	// Check if this is development mode without database
	if cfg.Environment == "development" && cfg.Database.Password == "" {
		logger.Warn().
			Msg("Development mode: No database password set - server will continue without database")
		return nil, fmt.Errorf("development mode: database not configured")
	}

	// Configure GORM logger to suppress "record not found" messages
	gormConfig := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	}
	if cfg.Environment == "development" {
		// In development, only log errors (not ErrRecordNotFound)
		gormConfig.Logger = gormlogger.Default.LogMode(gormlogger.Error)
	}

	// Retry connection up to 5 times
	for i := 0; i < 5; i++ {
		db, err := gorm.Open(postgres.Open(dsn), gormConfig)
		if err != nil {
			logger.Warn().
				Int("attempt", i+1).
				Err(err).
				Msg("Failed to open database connection")
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		// Get underlying sql.DB for connection pool configuration
		sqlDB, err := db.DB()
		if err != nil {
			logger.Warn().
				Int("attempt", i+1).
				Err(err).
				Msg("Failed to get underlying SQL DB")
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		// Configure connection pool
		sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(cfg.Database.MaxLifetime)

		// Test the connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = sqlDB.PingContext(ctx)
		cancel()

		if err == nil {
			return db, nil
		}

		logger.Warn().
			Int("attempt", i+1).
			Err(err).
			Msg("Database ping failed")
		sqlDBClose, _ := db.DB()
		if sqlDBClose != nil {
			sqlDBClose.Close()
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after 5 attempts")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services:  make(map[string]string),
		Version:   Version,
	}

	// Check database
	if err := checkDatabase(); err != nil {
		status.Status = "degraded"
		status.Services["database"] = fmt.Sprintf("error: %v", err)
		logger.Debug().Err(err).Msg("Database health check failed")
	} else {
		status.Services["database"] = "healthy"
	}

	// Always return 200 for basic health check - let readiness handle critical dependencies
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	// Check if all dependencies are ready
	if err := checkDatabase(); err != nil {
		http.Error(w, fmt.Sprintf("Database not ready: %v", err), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}

func livenessHandler(w http.ResponseWriter, r *http.Request) {
	// Basic liveness check
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alive"))
}

func checkDatabase() error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}

func runAutoMigrations(database *gorm.DB) error {
	if database == nil {
		return fmt.Errorf("database connection is nil")
	}

	logger.Info().Msg("Running GORM auto-migrations")

	// Get all entities from the domain registry
	// New entities only need to be added to domain.AllEntities()
	entities := domain.AllEntities()

	// Run auto-migration for all entities
	for _, entity := range entities {
		if err := database.AutoMigrate(entity); err != nil {
			return fmt.Errorf("failed to auto-migrate entity %T: %w", entity, err)
		}
	}

	// For now, just ensure the connection works
	sqlDB, err := database.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	logger.Info().
		Int("entity_count", len(entities)).
		Msg("GORM auto-migrations completed")
	return nil
}
