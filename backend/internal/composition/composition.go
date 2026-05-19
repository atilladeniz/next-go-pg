// Package composition is the application's composition root. It owns
// the dependency graph — DB connections, persistence + SSE adapters,
// use cases, handlers, the HTTP router, and background workers — and
// returns a ready-to-run *App. cmd/server only loads config, calls
// Build, runs the server, and shuts it down on signal.
package composition

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riverqueue/river"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/atilladeniz/next-go-pg/backend/internal/application"
	"github.com/atilladeniz/next-go-pg/backend/internal/handler"
	"github.com/atilladeniz/next-go-pg/backend/internal/infrastructure/persistence"
	"github.com/atilladeniz/next-go-pg/backend/internal/infrastructure/sse"
	"github.com/atilladeniz/next-go-pg/backend/internal/jobs"
	"github.com/atilladeniz/next-go-pg/backend/internal/middleware"
	"github.com/atilladeniz/next-go-pg/backend/pkg/config"
	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
	riverPkg "github.com/atilladeniz/next-go-pg/backend/pkg/river"
)

// Inputs to the composition root.
type Inputs struct {
	Config    *config.Config
	Version   string
	BuildTime string
}

// App is the fully wired application graph.
type App struct {
	HTTPServer *http.Server

	db            *gorm.DB
	pgxPool       *pgxpool.Pool
	riverJobQueue *riverPkg.Client
}

// Build assembles the dependency graph. It never returns a hard error
// for missing DB or River — those degrade to a server that still
// serves the public/static endpoints.
func Build(ctx context.Context, in Inputs) (*App, error) {
	cfg := in.Config

	app := &App{}

	// Database (optional in dev). Failures degrade to nil; the app
	// still boots so health endpoints can report the state.
	db, err := connectToDatabase(cfg)
	if err != nil {
		logger.Warn().Err(err).Msg("Database connection failed - server will start in degraded mode")
	} else {
		logger.Info().Msg("Database connected successfully")
		if err := runAutoMigrations(db); err != nil {
			logger.Warn().Err(err).Msg("Auto-migration failed - you may need to run migrations manually")
		} else {
			logger.Info().Msg("Database schema is up to date")
		}
		app.db = db
	}

	// SSE broker is independent of the DB.
	sseBroker := sse.NewBroker()
	logger.Info().Msg("SSE broker initialized")

	// Persistence + use cases (only with DB).
	var statsRepo application.StatsRepository
	var userDirectory application.UserDirectory
	var getStatsUC *application.GetUserStats
	var incrementStatUC *application.IncrementStatField
	if db != nil {
		statsRepo = persistence.NewUserStatsRepository(db)
		userDirectory = persistence.NewUserDirectoryRepository(db)
		getStatsUC = &application.GetUserStats{Repo: statsRepo}
		incrementStatUC = &application.IncrementStatField{Repo: statsRepo, Events: sseBroker}
	}

	// River (background jobs) — only with DB and a healthy pgx pool.
	var exportStore *jobs.ExportStore
	if db != nil {
		if pool, err := pgxpool.New(ctx, cfg.GetDatabaseURLForPgx()); err != nil {
			logger.Warn().Err(err).Msg("Failed to create pgx pool for River - background jobs disabled")
		} else {
			app.pgxPool = pool
			if err := riverPkg.RunMigrations(ctx, pool); err != nil {
				logger.Warn().Err(err).Msg("River migrations failed - background jobs may not work")
			}
			exportStore = jobs.NewExportStore()

			workers := river.NewWorkers()
			jobs.RegisterWorkers(workers, &jobs.WorkerDeps{
				EmailConfig: emailConfigFromEnv(),
				Events:      sseBroker,
				ExportStore: exportStore,
				StatsRepo:   statsRepo,
			})

			client, err := riverPkg.NewClient(ctx, pool, workers, riverPkg.DefaultConfig())
			if err != nil {
				logger.Warn().Err(err).Msg("Failed to create River client - background jobs disabled")
			} else {
				if err := client.Start(ctx); err != nil {
					logger.Error().Err(err).Msg("Failed to start River client")
				} else {
					logger.Info().Msg("River job queue initialized and started")
					app.riverJobQueue = client
				}
			}
		}
	}

	// HTTP layer.
	apiHandler := handler.NewAPIHandler(getStatsUC, incrementStatUC)
	webhookHandler := handler.NewWebhookHandler(userDirectory)
	if app.riverJobQueue != nil {
		webhookHandler = webhookHandler.WithJobEnqueuer(app.riverJobQueue.Client)
	}

	combinedAuth := middleware.NewCombinedAuthMiddleware(cfg.FrontendURL)

	router := buildRouter(routerDeps{
		cfg:            cfg,
		version:        in.Version,
		db:             db,
		sseBroker:      sseBroker,
		apiHandler:     apiHandler,
		webhookHandler: webhookHandler,
		combinedAuth:   combinedAuth,
		river:          app.riverJobQueue,
		exportStore:    exportStore,
	})

	app.HTTPServer = &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return app, nil
}

// Shutdown stops River, closes the pgx pool, and shuts down the HTTP
// server. Errors are logged but do not abort subsequent steps.
func (a *App) Shutdown(ctx context.Context) {
	if a.riverJobQueue != nil {
		logger.Info().Msg("Stopping River job queue...")
		if err := a.riverJobQueue.Stop(ctx); err != nil {
			logger.Error().Err(err).Msg("River job queue shutdown error")
		}
	}
	if a.pgxPool != nil {
		a.pgxPool.Close()
	}
	if a.HTTPServer != nil {
		if err := a.HTTPServer.Shutdown(ctx); err != nil {
			logger.Error().Err(err).Msg("Server forced to shutdown")
		}
	}
}

// --- internals -----------------------------------------------------

type routerDeps struct {
	cfg            *config.Config
	version        string
	db             *gorm.DB
	sseBroker      *sse.Broker
	apiHandler     *handler.APIHandler
	webhookHandler *handler.WebhookHandler
	combinedAuth   *middleware.CombinedAuthMiddleware
	river          *riverPkg.Client
	exportStore    *jobs.ExportStore
}

func buildRouter(d routerDeps) http.Handler {
	router := mux.NewRouter()

	loggingMW := middleware.NewLoggingMiddleware()
	corsMW := middleware.NewCORSMiddleware(d.cfg.FrontendURL)
	rateLimitMW := middleware.NewRateLimitMiddleware(middleware.RateLimitConfig{
		RequestsPerMinute: d.cfg.RateLimit.RequestsPerMinute,
		BurstSize:         d.cfg.RateLimit.BurstSize,
		SkipPaths:         []string{"/health", "/health/ready", "/health/live", "/metrics"},
	})
	metricsMW := middleware.NewMetricsMiddleware()

	router.Use(metricsMW.Handler)
	router.Use(loggingMW.Handler)
	router.Use(corsMW.Handler)
	router.Use(rateLimitMW.Handler)

	health := &healthEndpoints{db: d.db, version: d.version}
	router.HandleFunc("/health", health.health).Methods("GET")
	router.HandleFunc("/health/ready", health.ready).Methods("GET")
	router.HandleFunc("/health/live", health.live).Methods("GET")
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/hello", d.apiHandler.PublicHello).Methods("GET")

	protectedRouter := apiRouter.PathPrefix("/protected").Subrouter()
	protectedRouter.Use(d.combinedAuth.RequireAuth)
	protectedRouter.HandleFunc("/hello", d.apiHandler.ProtectedHello).Methods("GET")

	apiRouter.Handle("/me", d.combinedAuth.RequireAuth(http.HandlerFunc(d.apiHandler.GetCurrentUser))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/stats", d.combinedAuth.RequireAuth(http.HandlerFunc(d.apiHandler.GetUserStats))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/stats", d.combinedAuth.RequireAuth(http.HandlerFunc(d.apiHandler.UpdateUserStats))).Methods("POST", "OPTIONS")

	apiRouter.Handle("/events", d.sseBroker).Methods("GET")
	apiRouter.HandleFunc("/trigger-update", func(w http.ResponseWriter, r *http.Request) {
		d.sseBroker.Broadcast("stats-updated", `{"trigger":"manual"}`)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "broadcast sent"})
	}).Methods("POST")

	webhookRouter := apiRouter.PathPrefix("/webhooks").Subrouter()
	webhookRouter.HandleFunc("/session-created", d.webhookHandler.SessionCreated).Methods("POST")
	webhookRouter.HandleFunc("/send-magic-link", d.webhookHandler.SendMagicLink).Methods("POST")
	webhookRouter.HandleFunc("/send-verification-email", d.webhookHandler.SendVerificationEmail).Methods("POST")
	webhookRouter.HandleFunc("/send-2fa-otp", d.webhookHandler.Send2FAOTP).Methods("POST")
	webhookRouter.HandleFunc("/send-2fa-enabled", d.webhookHandler.Send2FAEnabledNotification).Methods("POST")
	webhookRouter.HandleFunc("/send-passkey-added", d.webhookHandler.SendPasskeyAddedNotification).Methods("POST")

	if d.river != nil && d.exportStore != nil {
		exportHandler := handler.NewExportHandler(d.river.Client, d.exportStore)
		apiRouter.Handle("/export/start", d.combinedAuth.RequireAuth(http.HandlerFunc(exportHandler.StartExport))).Methods("POST", "OPTIONS")
		apiRouter.HandleFunc("/export/download/{id}", exportHandler.DownloadExport).Methods("GET")
	}

	return router
}

type healthEndpoints struct {
	db      *gorm.DB
	version string
}

type healthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
}

func (h *healthEndpoints) health(w http.ResponseWriter, _ *http.Request) {
	status := healthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services:  make(map[string]string),
		Version:   h.version,
	}
	if err := h.pingDB(); err != nil {
		status.Status = "degraded"
		status.Services["database"] = fmt.Sprintf("error: %v", err)
	} else {
		status.Services["database"] = "healthy"
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}

func (h *healthEndpoints) ready(w http.ResponseWriter, _ *http.Request) {
	if err := h.pingDB(); err != nil {
		http.Error(w, fmt.Sprintf("Database not ready: %v", err), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Ready"))
}

func (h *healthEndpoints) live(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Alive"))
}

func (h *healthEndpoints) pingDB() error {
	if h.db == nil {
		return fmt.Errorf("database connection is nil")
	}
	sqlDB, err := h.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql DB: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return sqlDB.PingContext(ctx)
}

func connectToDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDatabaseURL()

	logger.Info().
		Str("host", cfg.Database.Host).
		Str("port", cfg.Database.Port).
		Str("database", cfg.Database.Name).
		Msg("Connecting to database")

	if cfg.Environment == "development" && cfg.Database.Password == "" {
		logger.Warn().Msg("Development mode: No database password set - server will continue without database")
		return nil, fmt.Errorf("development mode: database not configured")
	}

	gormConfig := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	}
	if cfg.Environment == "development" {
		gormConfig.Logger = gormlogger.Default.LogMode(gormlogger.Error)
	}

	for attempt := range 5 {
		db, err := gorm.Open(postgres.Open(dsn), gormConfig)
		if err != nil {
			logger.Warn().Int("attempt", attempt+1).Err(err).Msg("Failed to open database connection")
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}

		sqlDB, err := db.DB()
		if err != nil {
			logger.Warn().Int("attempt", attempt+1).Err(err).Msg("Failed to get underlying SQL DB")
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}

		sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(cfg.Database.MaxLifetime)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = sqlDB.PingContext(ctx)
		cancel()

		if err == nil {
			return db, nil
		}

		logger.Warn().Int("attempt", attempt+1).Err(err).Msg("Database ping failed")
		if closer, _ := db.DB(); closer != nil {
			_ = closer.Close()
		}
		time.Sleep(time.Duration(attempt+1) * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after 5 attempts")
}

func runAutoMigrations(database *gorm.DB) error {
	if database == nil {
		return fmt.Errorf("database connection is nil")
	}
	logger.Info().Msg("Running GORM auto-migrations")

	entities := persistence.AllEntities()
	for _, entity := range entities {
		if err := database.AutoMigrate(entity); err != nil {
			return fmt.Errorf("failed to auto-migrate entity %T: %w", entity, err)
		}
	}

	sqlDB, err := database.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	logger.Info().Int("entity_count", len(entities)).Msg("GORM auto-migrations completed")
	return nil
}

func emailConfigFromEnv() *jobs.EmailConfig {
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
	return jobs.NewEmailConfig(smtpHost, smtpPort, smtpFrom, appURL)
}
