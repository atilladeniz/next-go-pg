package middleware

import (
	"bufio"
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
	"github.com/google/uuid"
)

// responseWriter wraps http.ResponseWriter to capture status code and size
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int64
	wroteHeader  bool
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.statusCode = code
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += int64(n)
	return n, err
}

// Hijack implements http.Hijacker for WebSocket support
func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

// Flush implements http.Flusher for SSE support
func (rw *responseWriter) Flush() {
	if flusher, ok := rw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// LoggingMiddleware logs HTTP requests with request ID tracing
type LoggingMiddleware struct {
	// SkipPaths are paths that should not be logged (e.g., health checks)
	SkipPaths map[string]bool
	// SkipHealthLogs skips logging for health check endpoints
	SkipHealthLogs bool
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{
		SkipPaths: map[string]bool{
			"/health":       true,
			"/health/ready": true,
			"/health/live":  true,
			"/favicon.ico":  true,
		},
		SkipHealthLogs: true,
	}
}

// Handler returns the middleware handler
func (m *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Skip logging for certain paths
		if m.SkipHealthLogs && m.SkipPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		// Generate or extract request ID
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add request ID to response header
		w.Header().Set("X-Request-ID", requestID)

		// Add request ID to context
		ctx := logger.AddRequestIDToContext(r.Context(), requestID)

		// Check for distributed tracing headers
		if traceID := r.Header.Get("X-Trace-ID"); traceID != "" {
			spanID := r.Header.Get("X-Span-ID")
			ctx = logger.AddTraceToContext(ctx, traceID, spanID)
		}

		// Wrap response writer to capture status code and size
		wrapped := newResponseWriter(w)

		// Process request
		next.ServeHTTP(wrapped, r.WithContext(ctx))

		// Calculate latency
		latency := time.Since(start)

		// Log the request with full details (IP anonymized if configured)
		logger.HTTPRequest(ctx, logger.HTTPRequestLog{
			Method:       r.Method,
			Path:         r.URL.Path,
			Query:        r.URL.RawQuery,
			StatusCode:   wrapped.statusCode,
			Latency:      latency,
			ClientIP:     logger.AnonymizeIP(getClientIP(r)),
			UserAgent:    r.UserAgent(),
			Referer:      r.Referer(),
			RequestSize:  r.ContentLength,
			ResponseSize: wrapped.bytesWritten,
		})
	})
}

// getClientIP extracts the real client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for reverse proxies like nginx, cloudflare)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// Check X-Real-IP header (nginx)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Check CF-Connecting-IP header (Cloudflare)
	if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
		return cfIP
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// RequestIDMiddleware is a simpler middleware that only adds request ID
// Use this if you want to separate concerns
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		w.Header().Set("X-Request-ID", requestID)
		ctx := logger.AddRequestIDToContext(r.Context(), requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserContextMiddleware adds user information to the context after authentication
// This should be used after the auth middleware
func UserContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get user from context (set by auth middleware)
		if user := GetUserFromContext(ctx); user != nil {
			ctx = logger.AddUserToContext(ctx, user.ID, user.Name, user.Email, "")
		}

		// Get session from context
		if session := GetSessionFromContext(ctx); session != nil {
			ctx = context.WithValue(ctx, logger.SessionIDKey, session.ID)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
