package logger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
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
	UserEmailKey ContextKey = "user_email"
	UserNameKey  ContextKey = "user_name"
	SessionIDKey ContextKey = "session_id"
	TraceIDKey   ContextKey = "trace_id"
	SpanIDKey    ContextKey = "span_id"
)

// Log categories for filtering
type Category string

const (
	CategoryHTTP     Category = "http"
	CategoryAuth     Category = "auth"
	CategoryDB       Category = "db"
	CategoryBusiness Category = "business"
	CategorySystem   Category = "system"
	CategorySSE      Category = "sse"
	CategoryCache    Category = "cache"
)

// Config holds logger configuration
type Config struct {
	Level        string // debug, info, warn, error
	Environment  string // development, production
	ServiceName  string
	Version      string
	WithCaller   bool   // Include file:line in logs
	AnonymizeIPs bool   // GDPR: Anonymize IP addresses in logs
	LokiURL      string // Loki push URL (e.g., http://localhost:3100/loki/api/v1/push)
}

// Global config reference for runtime checks
var globalConfig Config

// Sensitive headers that should never be logged
var sensitiveHeaders = map[string]bool{
	"authorization":       true,
	"cookie":              true,
	"set-cookie":          true,
	"x-api-key":           true,
	"x-auth-token":        true,
	"x-csrf-token":        true,
	"proxy-authorization": true,
}

// Sensitive field patterns for redaction
var sensitiveFieldPatterns = []string{
	"password", "passwd", "pwd",
	"secret", "token", "key",
	"credential", "auth",
	"credit_card", "creditcard", "card_number",
	"cvv", "cvc", "ssn", "pin",
	"private", "api_key", "apikey",
}

// Global Loki writer reference for cleanup
var lokiWriter *LokiWriter

// Init initializes the global logger with the given configuration
func Init(cfg Config) {
	globalConfig = cfg

	var writers []io.Writer

	// Console output
	var consoleOutput io.Writer = os.Stdout
	if cfg.Environment == "development" {
		consoleOutput = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
			NoColor:    false,
			FormatLevel: func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
			},
			FormatFieldName: func(i interface{}) string {
				return fmt.Sprintf("%s=", i)
			},
		}
	}
	writers = append(writers, consoleOutput)

	// Loki output (if configured)
	if cfg.LokiURL != "" {
		lokiWriter = NewLokiWriter(cfg.LokiURL, cfg.ServiceName)
		// For Loki we need JSON output, not console format
		writers = append(writers, lokiWriter)
	}

	// Combine writers
	output := io.MultiWriter(writers...)

	// Parse log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}

	// Set global level
	zerolog.SetGlobalLevel(level)

	// Configure timestamp format
	zerolog.TimeFieldFormat = time.RFC3339

	// Create base logger
	logCtx := zerolog.New(output).
		With().
		Timestamp().
		Str("service", cfg.ServiceName).
		Str("version", cfg.Version)

	// Add caller info in production for debugging
	if cfg.WithCaller || cfg.Environment == "production" {
		logCtx = logCtx.Caller()
	}

	log = logCtx.Logger()
}

// Close cleans up logger resources (call on shutdown)
func Close() {
	if lokiWriter != nil {
		lokiWriter.Close()
	}
}

// Logger returns the global logger instance
func Logger() *zerolog.Logger {
	return &log
}

// WithContext returns a logger with all context values
func WithContext(ctx context.Context) zerolog.Logger {
	if ctx == nil {
		return log
	}

	l := log.With()

	// Request tracing
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		l = l.Str("request_id", requestID)
	}
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		l = l.Str("trace_id", traceID)
	}
	if spanID, ok := ctx.Value(SpanIDKey).(string); ok && spanID != "" {
		l = l.Str("span_id", spanID)
	}

	// User context
	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		l = l.Str("user_id", userID)
	}
	if userName, ok := ctx.Value(UserNameKey).(string); ok && userName != "" {
		l = l.Str("user_name", userName)
	}
	if sessionID, ok := ctx.Value(SessionIDKey).(string); ok && sessionID != "" {
		l = l.Str("session_id", sessionID)
	}

	return l.Logger()
}

// WithCategory returns a logger with a category field
func WithCategory(category Category) *zerolog.Event {
	return log.Info().Str("category", string(category))
}

// WithCategoryCtx returns a context-aware logger with category
func WithCategoryCtx(ctx context.Context, category Category) zerolog.Logger {
	return WithContext(ctx).With().Str("category", string(category)).Logger()
}

// Convenience methods for structured logging

func Info() *zerolog.Event  { return log.Info() }
func Debug() *zerolog.Event { return log.Debug() }
func Warn() *zerolog.Event  { return log.Warn() }
func Error() *zerolog.Event { return log.Error() }
func Fatal() *zerolog.Event { return log.Fatal() }

// WithError adds an error to the log event
func WithError(err error) *zerolog.Event {
	return log.Error().Err(err)
}

// ============================================================================
// HTTP Logging
// ============================================================================

// HTTPRequestLog contains all HTTP request information
type HTTPRequestLog struct {
	Method       string
	Path         string
	Query        string
	StatusCode   int
	Latency      time.Duration
	ClientIP     string
	UserAgent    string
	Referer      string
	RequestSize  int64
	ResponseSize int64
}

// HTTPRequest logs HTTP request with full details
func HTTPRequest(ctx context.Context, req HTTPRequestLog) {
	l := WithCategoryCtx(ctx, CategoryHTTP)

	// Determine log level based on status code
	var event *zerolog.Event
	switch {
	case req.StatusCode >= 500:
		event = l.Error()
	case req.StatusCode >= 400:
		event = l.Warn()
	case req.StatusCode >= 300:
		event = l.Info()
	default:
		event = l.Info()
	}

	event.
		Str("method", req.Method).
		Str("path", req.Path).
		Int("status", req.StatusCode).
		Dur("latency", req.Latency).
		Str("client_ip", req.ClientIP).
		Str("user_agent", truncate(req.UserAgent, 100)).
		Int64("request_size", req.RequestSize).
		Int64("response_size", req.ResponseSize)

	if req.Query != "" {
		event.Str("query", truncate(req.Query, 200))
	}

	event.Msg("HTTP request")
}

// HTTPRequestSimple is a simpler version for backwards compatibility
func HTTPRequestSimple(ctx context.Context, method, path string, statusCode int, latency time.Duration, clientIP string) {
	HTTPRequest(ctx, HTTPRequestLog{
		Method:     method,
		Path:       path,
		StatusCode: statusCode,
		Latency:    latency,
		ClientIP:   clientIP,
	})
}

// ============================================================================
// Database Logging
// ============================================================================

// DBQuery logs a database query
func DBQuery(ctx context.Context, operation, table string, latency time.Duration, rowsAffected int64, err error) {
	l := WithCategoryCtx(ctx, CategoryDB)

	event := l.Debug()
	if err != nil {
		event = l.Error().Err(err)
	} else if latency > 100*time.Millisecond {
		event = l.Warn() // Slow query warning
	}

	event.
		Str("operation", operation).
		Str("table", table).
		Dur("latency", latency).
		Int64("rows_affected", rowsAffected)

	if latency > 100*time.Millisecond {
		event.Bool("slow_query", true)
	}

	event.Msg("Database query")
}

// DBConnection logs database connection events
func DBConnection(event string, details map[string]any) {
	l := log.Info().
		Str("category", string(CategoryDB)).
		Str("event", event)

	if details != nil {
		l.Interface("details", details)
	}

	l.Msg("Database connection")
}

// ============================================================================
// Authentication Logging
// ============================================================================

// AuthEvent represents an authentication event
type AuthEvent struct {
	Action    string // login, logout, register, password_reset, token_refresh
	UserID    string
	Email     string
	Success   bool
	Reason    string // failure reason if not successful
	IP        string
	UserAgent string
}

// Auth logs authentication events
func Auth(ctx context.Context, evt AuthEvent) {
	l := WithCategoryCtx(ctx, CategoryAuth)

	var event *zerolog.Event
	if evt.Success {
		event = l.Info()
	} else {
		event = l.Warn()
	}

	event.
		Str("action", evt.Action).
		Bool("success", evt.Success)

	if evt.UserID != "" {
		event.Str("user_id", evt.UserID)
	}
	if evt.Email != "" {
		event.Str("email", maskEmail(evt.Email))
	}
	if evt.IP != "" {
		event.Str("ip", evt.IP)
	}
	if !evt.Success && evt.Reason != "" {
		event.Str("reason", evt.Reason)
	}

	event.Msg("Auth event")
}

// AuthSuccess is a convenience method for successful auth
func AuthSuccess(ctx context.Context, action, userID, email string) {
	Auth(ctx, AuthEvent{
		Action:  action,
		UserID:  userID,
		Email:   email,
		Success: true,
	})
}

// AuthFailure is a convenience method for failed auth
func AuthFailure(ctx context.Context, action, email, reason string) {
	Auth(ctx, AuthEvent{
		Action:  action,
		Email:   email,
		Success: false,
		Reason:  reason,
	})
}

// ============================================================================
// Business Event Logging (Audit Trail)
// ============================================================================

// BusinessEvent logs business domain events for audit trails
type BusinessEventLog struct {
	EventType string         // user.created, order.placed, payment.processed
	Entity    string         // user, order, payment
	EntityID  string         // The ID of the affected entity
	Action    string         // create, update, delete, view
	Changes   map[string]any // What changed (old -> new values)
	Metadata  map[string]any // Additional context
}

// Business logs a business event
func Business(ctx context.Context, evt BusinessEventLog) {
	l := WithCategoryCtx(ctx, CategoryBusiness)

	event := l.Info().
		Str("event_type", evt.EventType).
		Str("entity", evt.Entity).
		Str("entity_id", evt.EntityID).
		Str("action", evt.Action)

	if evt.Changes != nil {
		event.Interface("changes", evt.Changes)
	}
	if evt.Metadata != nil {
		event.Interface("metadata", evt.Metadata)
	}

	event.Msg("Business event")
}

// BusinessEvent is a simpler version for backwards compatibility
func BusinessEvent(ctx context.Context, eventType, entity, entityID string, details map[string]any) {
	Business(ctx, BusinessEventLog{
		EventType: eventType,
		Entity:    entity,
		EntityID:  entityID,
		Metadata:  details,
	})
}

// ============================================================================
// System Logging
// ============================================================================

// SystemEvent logs system-level events (startup, shutdown, config changes)
func SystemEvent(event string, details map[string]any) {
	l := log.Info().
		Str("category", string(CategorySystem)).
		Str("event", event)

	if details != nil {
		l.Interface("details", details)
	}

	l.Msg("System event")
}

// SystemError logs system-level errors
func SystemError(err error, event string, details map[string]any) {
	l := log.Error().
		Str("category", string(CategorySystem)).
		Str("event", event).
		Err(err)

	if details != nil {
		l.Interface("details", details)
	}

	l.Msg("System error")
}

// ============================================================================
// SSE Logging
// ============================================================================

// SSEEvent logs Server-Sent Events
func SSEEvent(ctx context.Context, eventType, channel string, clientCount int) {
	l := WithCategoryCtx(ctx, CategorySSE)
	l.Debug().
		Str("event_type", eventType).
		Str("channel", channel).
		Int("client_count", clientCount).
		Msg("SSE event")
}

// SSEConnection logs SSE connection events
func SSEConnection(ctx context.Context, action string, clientID string) {
	l := WithCategoryCtx(ctx, CategorySSE)
	l.Info().
		Str("action", action).
		Str("client_id", clientID).
		Msg("SSE connection")
}

// ============================================================================
// Performance Logging
// ============================================================================

// Timed returns a function that logs the duration of an operation
func Timed(ctx context.Context, operation string) func() {
	start := time.Now()
	return func() {
		l := WithContext(ctx)
		l.Debug().
			Str("operation", operation).
			Dur("duration", time.Since(start)).
			Msg("Operation completed")
	}
}

// TimedWithThreshold logs only if duration exceeds threshold
func TimedWithThreshold(ctx context.Context, operation string, threshold time.Duration) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start)
		if duration > threshold {
			l := WithContext(ctx)
			l.Warn().
				Str("operation", operation).
				Dur("duration", duration).
				Dur("threshold", threshold).
				Msg("Slow operation")
		}
	}
}

// ============================================================================
// Helper functions - Privacy & Security
// ============================================================================

// maskEmail masks email for privacy in logs (GDPR compliant)
func maskEmail(email string) string {
	if len(email) < 5 {
		return "***@***"
	}
	atIdx := strings.Index(email, "@")
	if atIdx <= 0 {
		return "***@***"
	}

	// Get domain
	domain := email[atIdx:]
	// Mask local part, show first char only
	if atIdx <= 2 {
		return email[:1] + "***" + domain
	}
	return email[:2] + "***" + domain
}

// AnonymizeIP anonymizes an IP address for GDPR compliance
// IPv4: 192.168.1.100 -> 192.168.1.0
// IPv6: 2001:db8::1 -> 2001:db8::0
func AnonymizeIP(ip string) string {
	if !globalConfig.AnonymizeIPs {
		return ip
	}

	// IPv4
	if strings.Contains(ip, ".") && !strings.Contains(ip, ":") {
		parts := strings.Split(ip, ".")
		if len(parts) == 4 {
			return parts[0] + "." + parts[1] + "." + parts[2] + ".0"
		}
	}

	// IPv6
	if strings.Contains(ip, ":") {
		// Simple approach: zero out last 80 bits (keep /48)
		parts := strings.Split(ip, ":")
		if len(parts) >= 3 {
			return parts[0] + ":" + parts[1] + ":" + parts[2] + "::0"
		}
	}

	return ip
}

// RedactSensitive checks if a field name is sensitive and returns redacted value
func RedactSensitive(fieldName string, value string) string {
	fieldLower := strings.ToLower(fieldName)
	for _, pattern := range sensitiveFieldPatterns {
		if strings.Contains(fieldLower, pattern) {
			return "[REDACTED]"
		}
	}
	return value
}

// IsSensitiveHeader checks if a header should not be logged
func IsSensitiveHeader(headerName string) bool {
	return sensitiveHeaders[strings.ToLower(headerName)]
}

// RedactMap redacts sensitive fields from a map
func RedactMap(data map[string]any) map[string]any {
	if data == nil {
		return nil
	}

	result := make(map[string]any, len(data))
	for k, v := range data {
		fieldLower := strings.ToLower(k)
		isSensitive := false
		for _, pattern := range sensitiveFieldPatterns {
			if strings.Contains(fieldLower, pattern) {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			result[k] = "[REDACTED]"
		} else if nestedMap, ok := v.(map[string]any); ok {
			result[k] = RedactMap(nestedMap)
		} else {
			result[k] = v
		}
	}
	return result
}

// truncate truncates a string to maxLen
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// GetCallerInfo returns file:line of the caller
func GetCallerInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	// Get just the filename, not full path
	parts := strings.Split(file, "/")
	if len(parts) > 0 {
		file = parts[len(parts)-1]
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// ============================================================================
// Context Helpers
// ============================================================================

// AddUserToContext adds user information to context for logging
func AddUserToContext(ctx context.Context, userID, userName, email, sessionID string) context.Context {
	ctx = context.WithValue(ctx, UserIDKey, userID)
	ctx = context.WithValue(ctx, UserNameKey, userName)
	ctx = context.WithValue(ctx, UserEmailKey, email)
	ctx = context.WithValue(ctx, SessionIDKey, sessionID)
	return ctx
}

// AddRequestIDToContext adds request ID to context
func AddRequestIDToContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// AddTraceToContext adds distributed tracing IDs to context
func AddTraceToContext(ctx context.Context, traceID, spanID string) context.Context {
	ctx = context.WithValue(ctx, TraceIDKey, traceID)
	ctx = context.WithValue(ctx, SpanIDKey, spanID)
	return ctx
}

// ============================================================================
// Loki Writer - Send logs directly to Loki
// ============================================================================

// LokiWriter sends logs to Loki via HTTP
type LokiWriter struct {
	url         string
	serviceName string
	client      *http.Client
	batch       []lokiEntry
	batchMu     sync.Mutex
	batchSize   int
	flushTicker *time.Ticker
	done        chan struct{}
}

type lokiEntry struct {
	timestamp time.Time
	line      string
	level     string
}

// NewLokiWriter creates a new Loki writer
func NewLokiWriter(url, serviceName string) *LokiWriter {
	lw := &LokiWriter{
		url:         url,
		serviceName: serviceName,
		client:      &http.Client{Timeout: 5 * time.Second},
		batch:       make([]lokiEntry, 0, 100),
		batchSize:   50,
		flushTicker: time.NewTicker(2 * time.Second),
		done:        make(chan struct{}),
	}

	// Background flusher
	go lw.backgroundFlush()

	return lw
}

// Write implements io.Writer for zerolog
func (lw *LokiWriter) Write(p []byte) (n int, err error) {
	// Parse level from JSON
	level := "info"
	line := string(p)
	if strings.Contains(line, `"level":"error"`) {
		level = "error"
	} else if strings.Contains(line, `"level":"warn"`) {
		level = "warn"
	} else if strings.Contains(line, `"level":"debug"`) {
		level = "debug"
	}

	lw.batchMu.Lock()
	lw.batch = append(lw.batch, lokiEntry{
		timestamp: time.Now(),
		line:      strings.TrimSpace(line),
		level:     level,
	})
	shouldFlush := len(lw.batch) >= lw.batchSize
	lw.batchMu.Unlock()

	if shouldFlush {
		go lw.flush()
	}

	return len(p), nil
}

func (lw *LokiWriter) backgroundFlush() {
	for {
		select {
		case <-lw.flushTicker.C:
			lw.flush()
		case <-lw.done:
			lw.flush() // Final flush
			return
		}
	}
}

func (lw *LokiWriter) flush() {
	lw.batchMu.Lock()
	if len(lw.batch) == 0 {
		lw.batchMu.Unlock()
		return
	}
	entries := lw.batch
	lw.batch = make([]lokiEntry, 0, 100)
	lw.batchMu.Unlock()

	// Group by level for separate streams
	streams := make(map[string][]lokiEntry)
	for _, e := range entries {
		streams[e.level] = append(streams[e.level], e)
	}

	// Build Loki payload
	var streamArr []string
	for level, levelEntries := range streams {
		var values []string
		for _, e := range levelEntries {
			ts := fmt.Sprintf("%d", e.timestamp.UnixNano())
			// Escape JSON string
			escaped := strings.ReplaceAll(e.line, `\`, `\\`)
			escaped = strings.ReplaceAll(escaped, `"`, `\"`)
			values = append(values, fmt.Sprintf(`["%s", "%s"]`, ts, escaped))
		}
		stream := fmt.Sprintf(`{"stream":{"service":"%s","level":"%s","job":"backend"},"values":[%s]}`,
			lw.serviceName, level, strings.Join(values, ","))
		streamArr = append(streamArr, stream)
	}

	payload := fmt.Sprintf(`{"streams":[%s]}`, strings.Join(streamArr, ","))

	req, err := http.NewRequest("POST", lw.url, bytes.NewBufferString(payload))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := lw.client.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()
}

// Close stops the background flusher
func (lw *LokiWriter) Close() {
	lw.flushTicker.Stop()
	close(lw.done)
}
