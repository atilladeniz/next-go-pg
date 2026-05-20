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
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atilladeniz/next-go-pg/backend/internal/composition"
	"github.com/atilladeniz/next-go-pg/backend/pkg/config"
	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
	"github.com/atilladeniz/next-go-pg/backend/pkg/metrics"
)

// Build information injected via -ldflags.
var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	cfg := config.Load()

	if cfg.Environment == "production" {
		cfg.MustValidate()
	} else if err := cfg.ValidateWithWarnings(); err != nil {
		log.Printf("Configuration warnings: %v", err)
	}

	logger.Init(logger.Config{
		Level:        cfg.LogLevel,
		Environment:  cfg.Environment,
		ServiceName:  "next-go-pg-api",
		Version:      Version,
		AnonymizeIPs: cfg.Logging.AnonymizeIPs,
		WithCaller:   cfg.Logging.WithCaller,
		LokiURL:      os.Getenv("LOKI_URL"),
	})
	defer logger.Close()

	logger.Info().
		Str("version", Version).
		Str("build_time", BuildTime).
		Str("environment", cfg.Environment).
		Msg("Starting application")

	metrics.Init(Version, cfg.Environment)

	ctx := context.Background()
	app, err := composition.Build(ctx, composition.Inputs{
		Config:    cfg,
		Version:   Version,
		BuildTime: BuildTime,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to build application")
	}

	go func() {
		logger.Info().
			Str("port", cfg.Port).
			Str("frontend_url", cfg.FrontendURL).
			Msg("Server starting")
		if err := app.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server startup failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	app.Shutdown(shutdownCtx)

	logger.Info().Msg("Server exited")
}
