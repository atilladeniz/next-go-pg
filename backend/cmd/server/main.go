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
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atilladeniz/next-go-pg/backend/internal/domain"
	"github.com/atilladeniz/next-go-pg/backend/internal/handler"
	"github.com/atilladeniz/next-go-pg/backend/internal/middleware"
	"github.com/atilladeniz/next-go-pg/backend/internal/repository"
	"github.com/atilladeniz/next-go-pg/backend/internal/sse"
	"github.com/atilladeniz/next-go-pg/backend/pkg/config"
	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	sseBroker *sse.Broker
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
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize structured logger
	logger.Init(logger.Config{
		Level:        cfg.LogLevel,
		Environment:  cfg.Environment,
		ServiceName:  "next-go-pg-api",
		Version:      Version,
		AnonymizeIPs: cfg.Logging.AnonymizeIPs,
		WithCaller:   cfg.Logging.WithCaller,
	})

	logger.Info().
		Str("version", Version).
		Str("build_time", BuildTime).
		Str("environment", cfg.Environment).
		Msg("Starting application")

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

	// Setup router
	router := mux.NewRouter()

	// Setup middlewares
	loggingMiddleware := middleware.NewLoggingMiddleware()
	corsMiddleware := middleware.NewCORSMiddleware(cfg.FrontendURL)

	// Apply middlewares (order matters: logging first to capture all requests)
	router.Use(loggingMiddleware.Handler)
	router.Use(corsMiddleware.Handler)

	// Health check endpoint with comprehensive checks
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")
	router.HandleFunc("/health/ready", readinessHandler).Methods("GET")
	router.HandleFunc("/health/live", livenessHandler).Methods("GET")

	// Setup repositories
	var statsRepo *repository.UserStatsRepository
	if db != nil {
		statsRepo = repository.NewUserStatsRepository(db)
	}

	// Setup API handlers
	apiHandler := handler.NewAPIHandler(cfg.FrontendURL, sseBroker, statsRepo)
	authMiddleware := apiHandler.GetAuthMiddleware()

	// API v1 routes
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// Public routes
	apiRouter.HandleFunc("/hello", apiHandler.PublicHello).Methods("GET")

	// Protected routes (require auth)
	protectedRouter := apiRouter.PathPrefix("/protected").Subrouter()
	protectedRouter.Use(authMiddleware.RequireAuth)
	protectedRouter.HandleFunc("/hello", apiHandler.ProtectedHello).Methods("GET")

	// User routes (require auth)
	apiRouter.Handle("/me", authMiddleware.RequireAuth(http.HandlerFunc(apiHandler.GetCurrentUser))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/stats", authMiddleware.RequireAuth(http.HandlerFunc(apiHandler.GetUserStats))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/stats", authMiddleware.RequireAuth(http.HandlerFunc(apiHandler.UpdateUserStats))).Methods("POST", "OPTIONS")

	// SSE endpoint for real-time updates
	apiRouter.Handle("/events", sseBroker).Methods("GET")

	// Test endpoint to trigger stats update (for demo)
	apiRouter.HandleFunc("/trigger-update", func(w http.ResponseWriter, r *http.Request) {
		sseBroker.Broadcast("stats-updated", `{"trigger":"manual"}`)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "broadcast sent"})
	}).Methods("POST")

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

	// Retry connection up to 5 times
	for i := 0; i < 5; i++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
