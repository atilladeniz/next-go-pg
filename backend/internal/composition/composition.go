// Package composition is the application's composition root. It owns
// the dependency graph across all bounded contexts (stats, auth,
// notifications, exports) plus cross-cutting platform services (SSE
// broker, HTTP middleware). cmd/server only loads config, calls Build,
// runs the server, and shuts it down on signal.
package composition

import (
	"context"
	"encoding/json"
	"errors"
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

	// Bounded context: auth
	authapp "github.com/atilladeniz/next-go-pg/backend/internal/auth/application"
	"github.com/atilladeniz/next-go-pg/backend/internal/auth/infrastructure/betterauth"
	authhttp "github.com/atilladeniz/next-go-pg/backend/internal/auth/interfaces/http"

	// Bounded context: notifications
	notifapp "github.com/atilladeniz/next-go-pg/backend/internal/notifications/application"
	notifemail "github.com/atilladeniz/next-go-pg/backend/internal/notifications/infrastructure/email"
	notifjobs "github.com/atilladeniz/next-go-pg/backend/internal/notifications/infrastructure/jobs"
	notifhttp "github.com/atilladeniz/next-go-pg/backend/internal/notifications/interfaces/http"

	// Bounded context: stats
	statsapp "github.com/atilladeniz/next-go-pg/backend/internal/stats/application"
	statsevents "github.com/atilladeniz/next-go-pg/backend/internal/stats/infrastructure/events"
	statspersist "github.com/atilladeniz/next-go-pg/backend/internal/stats/infrastructure/persistence"
	statshttp "github.com/atilladeniz/next-go-pg/backend/internal/stats/interfaces/http"

	// Bounded context: aiworkflows
	aiapp "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/application"
	aievents "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/infrastructure/events"
	aigit "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/infrastructure/git"
	aillm "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/infrastructure/llm"
	aipersist "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/infrastructure/persistence"
	aiworkflows "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/infrastructure/workflows"
	aihttp "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/interfaces/http"

	hatchet "github.com/hatchet-dev/hatchet/sdks/go"

	// Bounded context: exports
	exportsapp "github.com/atilladeniz/next-go-pg/backend/internal/exports/application"
	exportsinfra "github.com/atilladeniz/next-go-pg/backend/internal/exports/infrastructure"
	exportsjobs "github.com/atilladeniz/next-go-pg/backend/internal/exports/infrastructure/jobs"
	exportshttp "github.com/atilladeniz/next-go-pg/backend/internal/exports/interfaces/http"

	// Shared kernel
	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"

	// Platform (cross-cutting infrastructure)
	"github.com/atilladeniz/next-go-pg/backend/internal/platform/middleware"
	"github.com/atilladeniz/next-go-pg/backend/internal/platform/sse"

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
	sseBroker     *sse.Broker

	// AI workflow worker — running goroutine + cancel. Nil when
	// HATCHET_CLIENT_TOKEN is not set (degraded boot).
	hatchetWorker     *aiworkflows.Worker
	hatchetWorkerStop context.CancelFunc
}

// Build assembles the dependency graph.
func Build(ctx context.Context, in Inputs) (*App, error) {
	cfg := in.Config
	app := &App{}

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

	// Platform: SSE broker.
	sseBroker := sse.NewBroker()
	app.sseBroker = sseBroker
	logger.Info().Msg("SSE broker initialized")

	// Notifications context — email sender is always constructed.
	emailSender := notifemail.NewSender(emailConfigFromEnv())

	// Stats context.
	var statsRepo statsapp.Repository
	var getStatsUC *statsapp.GetUserStats
	var incrementStatUC *statsapp.IncrementStatField
	if db != nil {
		statsRepo = statspersist.NewRepository(db)
		statsPublisher := statsevents.NewPublisher(sseBroker)
		getStatsUC = &statsapp.GetUserStats{Repo: statsRepo}
		incrementStatUC = &statsapp.IncrementStatField{Repo: statsRepo, Events: statsPublisher}
	}

	// Auth context.
	var userDirectory authapp.UserDirectory
	if db != nil {
		userDirectory = betterauth.NewDirectory(db)
	}

	// Exports context — ACL over stats (composition-level adapter).
	var statsReader exportsapp.StatsReader
	if statsRepo != nil {
		statsReader = &statsToExportsReader{repo: statsRepo}
	}
	exportStore := exportsinfra.NewMemoryStore()

	// River queue — wires per-context workers.
	var notifEnqueuer notifapp.JobEnqueuer
	var exportsEnqueuer exportsapp.JobEnqueuer
	if db != nil {
		if pool, err := pgxpool.New(ctx, cfg.GetDatabaseURLForPgx()); err != nil {
			logger.Warn().Err(err).Msg("Failed to create pgx pool for River - background jobs disabled")
		} else {
			app.pgxPool = pool
			if err := riverPkg.RunMigrations(ctx, pool); err != nil {
				logger.Warn().Err(err).Msg("River migrations failed - background jobs may not work")
			}

			workers := river.NewWorkers()
			notifjobs.Register(workers, emailSender)
			exportsjobs.Register(workers, sseBroker, exportStore, statsReader)

			client, err := riverPkg.NewClient(ctx, pool, workers, riverPkg.DefaultConfig())
			if err != nil {
				logger.Warn().Err(err).Msg("Failed to create River client - background jobs disabled")
			} else {
				if err := client.Start(ctx); err != nil {
					logger.Error().Err(err).Msg("Failed to start River client")
				} else {
					logger.Info().Msg("River job queue initialized and started")
					app.riverJobQueue = client
					notifEnqueuer = notifjobs.NewEnqueuer(client.Client)
					exportsEnqueuer = exportsjobs.NewEnqueuer(client.Client)
				}
			}
		}
	}

	// AI workflows context — gated on HATCHET_CLIENT_TOKEN. Without it
	// we skip the Hatchet wiring entirely so `just dev` still boots
	// when the AI compose profile is down.
	var aiHandler *aihttp.Handler
	if db != nil {
		aiHandler = buildAIWorkflowsHandler(ctx, app, db, sseBroker)
	}

	// HTTP layer — per-context handlers.
	authHandler := authhttp.NewHandler()
	statsHandler := statshttp.NewHandler(getStatsUC, incrementStatUC)

	// ACL: notifications declares a local UserDirectory port. The
	// composition root adapts auth's UserDirectory to it so the two
	// contexts stay decoupled.
	var notifUsers notifapp.UserDirectory
	if userDirectory != nil {
		notifUsers = &authToNotificationsDirectory{users: userDirectory}
	}
	webhookHandler := notifhttp.NewHandler(notifUsers, emailSender)
	if notifEnqueuer != nil {
		webhookHandler = webhookHandler.WithJobEnqueuer(notifEnqueuer)
	}
	exportHandler := exportshttp.NewHandler(exportsEnqueuer, exportStore)

	combinedAuth := middleware.NewCombinedAuthMiddleware(cfg.FrontendURL)

	router := buildRouter(routerDeps{
		cfg:            cfg,
		version:        in.Version,
		db:             db,
		sseBroker:      sseBroker,
		authHandler:    authHandler,
		statsHandler:   statsHandler,
		webhookHandler: webhookHandler,
		exportHandler:  exportHandler,
		aiHandler:      aiHandler,
		combinedAuth:   combinedAuth,
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

// Shutdown stops the HTTP server, the Hatchet worker, River, the SSE
// broker and closes the pgx pool. Order matters: HTTP first so no new
// SSE connections arrive, Hatchet worker so no new tasks are claimed,
// then the broker drains existing clients, then the rest.
func (a *App) Shutdown(ctx context.Context) {
	if a.HTTPServer != nil {
		if err := a.HTTPServer.Shutdown(ctx); err != nil {
			logger.Error().Err(err).Msg("Server forced to shutdown")
		}
	}
	if a.hatchetWorkerStop != nil {
		logger.Info().Msg("Stopping Hatchet worker...")
		a.hatchetWorkerStop()
	}
	if a.sseBroker != nil {
		if err := a.sseBroker.Shutdown(ctx); err != nil {
			logger.Error().Err(err).Msg("SSE broker shutdown error")
		}
	}
	if a.riverJobQueue != nil {
		logger.Info().Msg("Stopping River job queue...")
		if err := a.riverJobQueue.Stop(ctx); err != nil {
			logger.Error().Err(err).Msg("River job queue shutdown error")
		}
	}
	if a.pgxPool != nil {
		a.pgxPool.Close()
	}
}

// buildAIWorkflowsHandler wires the aiworkflows bounded context end to
// end. It is gated on HATCHET_CLIENT_TOKEN: without a token we cannot
// dial hatchet-lite, so we skip the wiring and return a handler that
// responds with 503 (Service Unavailable). The store and use cases
// still work in degraded mode so GET /ai/summaries/{id} can answer for
// rows enqueued before a restart.
func buildAIWorkflowsHandler(ctx context.Context, app *App, db *gorm.DB, broker *sse.Broker) *aihttp.Handler {
	repo := aipersist.NewRepository(db)
	getUC := &aiapp.GetRepoSummary{Store: repo}

	token := os.Getenv("HATCHET_CLIENT_TOKEN")
	if token == "" {
		logger.Warn().Msg("HATCHET_CLIENT_TOKEN unset — AI workflows disabled (degraded boot). GET /ai/summaries/{id} still serves existing rows.")
		return aihttp.NewHandler(nil, getUC)
	}

	client, err := hatchet.NewClient()
	if err != nil {
		logger.Warn().Err(err).Msg("Hatchet client init failed — AI workflows disabled")
		return aihttp.NewHandler(nil, getUC)
	}

	llmCfg := aillm.Config{
		URL:   os.Getenv("OLLAMA_URL"),
		Model: os.Getenv("OLLAMA_MODEL"),
	}
	if raw := os.Getenv("OLLAMA_TIMEOUT"); raw != "" {
		if d, err := time.ParseDuration(raw); err == nil {
			llmCfg.Timeout = d
		} else {
			logger.Warn().Str("value", raw).Msg("OLLAMA_TIMEOUT not parseable as duration, falling back to default")
		}
	}
	maxFiles := 25
	if raw := os.Getenv("AI_MAX_FILES"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			maxFiles = n
		}
	}
	deps := aiworkflows.Deps{
		Cloner:   aigit.NewCloner("", 50*1024*1024),
		LLM:      aillm.NewClient(llmCfg),
		Store:    repo,
		Progress: aievents.NewPublisher(broker),
		MaxFiles: maxFiles,
		MaxBytes: 64 * 1024,
	}

	worker, err := aiworkflows.NewWorker(client, deps, "ai-workflows-worker")
	if err != nil {
		logger.Warn().Err(err).Msg("Hatchet worker init failed — AI workflows disabled")
		return aihttp.NewHandler(nil, getUC)
	}

	workerCtx, cancel := context.WithCancel(ctx)
	app.hatchetWorker = worker
	app.hatchetWorkerStop = cancel
	go func() {
		logger.Info().Str("worker", "ai-workflows-worker").Msg("Starting Hatchet worker")
		if err := worker.Start(workerCtx); err != nil && !errors.Is(err, context.Canceled) {
			logger.Error().Err(err).Msg("Hatchet worker exited with error")
		}
	}()

	enqueuer := aiworkflows.NewEnqueuer(client)
	summarizeUC := &aiapp.SummarizeRepo{Store: repo, Enqueuer: enqueuer}

	logger.Info().Msg("AI workflows context wired: Hatchet + Ollama")
	return aihttp.NewHandler(summarizeUC, getUC)
}

// statsToExportsReader is the anti-corruption layer between the stats
// and exports bounded contexts. Exports declares the shape it needs
// (StatsSnapshot); composition implements it against stats's port.
type statsToExportsReader struct {
	repo statsapp.Repository
}

func (r *statsToExportsReader) Read(ctx context.Context, userID string) (exportsapp.StatsSnapshot, error) {
	uid, err := shared.NewUserID(userID)
	if err != nil {
		return exportsapp.StatsSnapshot{}, err
	}
	s, err := r.repo.GetOrCreate(ctx, uid)
	if err != nil {
		return exportsapp.StatsSnapshot{}, err
	}
	return exportsapp.StatsSnapshot{
		Projects:      s.ProjectCount,
		Activity:      s.ActivityToday,
		Notifications: s.Notifications,
	}, nil
}

// authToNotificationsDirectory is the anti-corruption layer between
// the auth and notifications bounded contexts. Notifications declares
// the shape it needs (UserSnapshot); composition implements it against
// auth's UserDirectory port.
type authToNotificationsDirectory struct {
	users authapp.UserDirectory
}

func (a *authToNotificationsDirectory) UserByID(ctx context.Context, userID shared.UserID) (notifapp.UserSnapshot, error) {
	u, err := a.users.UserByID(ctx, userID)
	if err != nil {
		return notifapp.UserSnapshot{}, err
	}
	return notifapp.UserSnapshot{Email: u.Email, Name: u.Name}, nil
}

func (a *authToNotificationsDirectory) HasKnownDevice(ctx context.Context, userID shared.UserID, userAgent, ipAddress, excludeSessionID string) (bool, error) {
	return a.users.HasKnownDevice(ctx, userID, userAgent, ipAddress, excludeSessionID)
}

// --- internals -----------------------------------------------------

type routerDeps struct {
	cfg            *config.Config
	version        string
	db             *gorm.DB
	sseBroker      *sse.Broker
	authHandler    *authhttp.Handler
	statsHandler   *statshttp.Handler
	webhookHandler *notifhttp.Handler
	exportHandler  *exportshttp.Handler
	aiHandler      *aihttp.Handler
	combinedAuth   *middleware.CombinedAuthMiddleware
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
	apiRouter.HandleFunc("/hello", d.authHandler.PublicHello).Methods("GET")

	protectedRouter := apiRouter.PathPrefix("/protected").Subrouter()
	protectedRouter.Use(d.combinedAuth.RequireAuth)
	protectedRouter.HandleFunc("/hello", d.authHandler.ProtectedHello).Methods("GET")

	apiRouter.Handle("/me", d.combinedAuth.RequireAuth(http.HandlerFunc(d.authHandler.GetCurrentUser))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/stats", d.combinedAuth.RequireAuth(http.HandlerFunc(d.statsHandler.GetUserStats))).Methods("GET", "OPTIONS")
	apiRouter.Handle("/stats", d.combinedAuth.RequireAuth(http.HandlerFunc(d.statsHandler.UpdateUserStats))).Methods("POST", "OPTIONS")

	apiRouter.Handle("/events", d.sseBroker).Methods("GET")
	apiRouter.HandleFunc("/trigger-update", func(w http.ResponseWriter, _ *http.Request) {
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

	if d.exportHandler != nil {
		apiRouter.Handle("/export/start", d.combinedAuth.RequireAuth(http.HandlerFunc(d.exportHandler.StartExport))).Methods("POST", "OPTIONS")
		apiRouter.HandleFunc("/export/download/{id}", d.exportHandler.DownloadExport).Methods("GET")
	}

	if d.aiHandler != nil {
		apiRouter.Handle("/ai/summarize-repo", d.combinedAuth.RequireAuth(http.HandlerFunc(d.aiHandler.SummarizeRepo))).Methods("POST", "OPTIONS")
		apiRouter.Handle("/ai/summaries/{id}", d.combinedAuth.RequireAuth(http.HandlerFunc(d.aiHandler.GetRepoSummary))).Methods("GET", "OPTIONS")
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

	gormConfig := &gorm.Config{Logger: gormlogger.Default.LogMode(gormlogger.Silent)}
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

	// Collect entities from every context that owns persistence.
	entities := []any{}
	entities = append(entities, statspersist.Entities()...)
	entities = append(entities, aipersist.Entities()...)

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

func emailConfigFromEnv() notifemail.Config {
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
	return notifemail.Config{
		SMTPHost: smtpHost,
		SMTPPort: smtpPort,
		SMTPFrom: smtpFrom,
		AppURL:   appURL,
	}
}
