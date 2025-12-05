package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Global logger instance
var log zerolog.Logger

// ContextKey for request-scoped values
type ContextKey string

const (
	RequestIDKey ContextKey = "request_id"
	UserIDKey    ContextKey = "user_id"
)

// Config holds logger configuration
type Config struct {
	Level       string // debug, info, warn, error
	Environment string // development, production
	ServiceName string
	Version     string
}

// Init initializes the global logger with the given configuration
func Init(cfg Config) {
	var output io.Writer = os.Stdout

	// Pretty printing for development
	if cfg.Environment == "development" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}
	}

	// Parse log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}

	// Set global level
	zerolog.SetGlobalLevel(level)

	// Create logger with base fields
	log = zerolog.New(output).
		With().
		Timestamp().
		Str("service", cfg.ServiceName).
		Str("version", cfg.Version).
		Logger()
}

// Logger returns the global logger instance
func Logger() *zerolog.Logger {
	return &log
}

// WithContext returns a logger with context values (request_id, user_id)
func WithContext(ctx context.Context) zerolog.Logger {
	l := log.With()

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		l = l.Str("request_id", requestID)
	}

	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		l = l.Str("user_id", userID)
	}

	return l.Logger()
}

// Convenience methods for structured logging

// Info logs an info level message
func Info() *zerolog.Event {
	return log.Info()
}

// Debug logs a debug level message
func Debug() *zerolog.Event {
	return log.Debug()
}

// Warn logs a warning level message
func Warn() *zerolog.Event {
	return log.Warn()
}

// Error logs an error level message
func Error() *zerolog.Event {
	return log.Error()
}

// Fatal logs a fatal level message and exits
func Fatal() *zerolog.Event {
	return log.Fatal()
}

// WithError adds an error to the log event
func WithError(err error) *zerolog.Event {
	return log.Error().Err(err)
}

// HTTP request logging helpers

// HTTPRequest logs details of an HTTP request
func HTTPRequest(ctx context.Context, method, path string, statusCode int, latency time.Duration, clientIP string) {
	l := WithContext(ctx)

	event := l.Info()
	if statusCode >= 500 {
		event = l.Error()
	} else if statusCode >= 400 {
		event = l.Warn()
	}

	event.
		Str("method", method).
		Str("path", path).
		Int("status", statusCode).
		Dur("latency", latency).
		Str("client_ip", clientIP).
		Msg("HTTP request")
}

// Database logging helpers

// DBQuery logs a database query
func DBQuery(ctx context.Context, query string, latency time.Duration, err error) {
	l := WithContext(ctx)

	event := l.Debug()
	if err != nil {
		event = l.Error().Err(err)
	}

	event.
		Str("query", query).
		Dur("latency", latency).
		Msg("Database query")
}

// Auth logging helpers

// AuthSuccess logs successful authentication
func AuthSuccess(ctx context.Context, userID, email string) {
	l := WithContext(ctx)
	l.Info().
		Str("user_id", userID).
		Str("email", maskEmail(email)).
		Msg("Authentication successful")
}

// AuthFailure logs failed authentication
func AuthFailure(ctx context.Context, email, reason string) {
	l := WithContext(ctx)
	l.Warn().
		Str("email", maskEmail(email)).
		Str("reason", reason).
		Msg("Authentication failed")
}

// Business event logging

// BusinessEvent logs a business domain event (for audit trails)
func BusinessEvent(ctx context.Context, eventType, entity string, entityID string, details map[string]interface{}) {
	l := WithContext(ctx)

	event := l.Info().
		Str("event_type", eventType).
		Str("entity", entity).
		Str("entity_id", entityID)

	if details != nil {
		event = event.Interface("details", details)
	}

	event.Msg("Business event")
}

// Helper functions

// maskEmail masks email for privacy in logs
func maskEmail(email string) string {
	if len(email) < 5 {
		return "***"
	}
	// Show first 2 chars and domain
	atIdx := -1
	for i, c := range email {
		if c == '@' {
			atIdx = i
			break
		}
	}
	if atIdx <= 2 {
		return email[:1] + "***" + email[atIdx:]
	}
	return email[:2] + "***" + email[atIdx:]
}
