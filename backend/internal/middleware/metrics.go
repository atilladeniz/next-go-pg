package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/atilladeniz/next-go-pg/backend/pkg/metrics"
	"github.com/gorilla/mux"
)

// MetricsMiddleware records HTTP metrics for Prometheus
type MetricsMiddleware struct{}

// NewMetricsMiddleware creates a new metrics middleware
func NewMetricsMiddleware() *MetricsMiddleware {
	return &MetricsMiddleware{}
}

// Handler wraps an HTTP handler with metrics collection
func (m *MetricsMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip metrics endpoint to avoid recursion
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		// Track in-flight requests
		metrics.HTTPRequestsInFlight.Inc()
		defer metrics.HTTPRequestsInFlight.Dec()

		// Get route pattern for consistent labeling
		path := getRoutePath(r)

		// Measure request duration
		start := time.Now()

		// Use existing responseWriter from logging middleware
		wrapped := newResponseWriter(w)

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(wrapped.statusCode)

		// Record metrics
		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		metrics.HTTPRequestDuration.WithLabelValues(r.Method, path).Observe(duration)

		// Track rate limit hits
		if wrapped.statusCode == http.StatusTooManyRequests {
			metrics.RateLimitHits.Inc()
		}
	})
}

// getRoutePath returns the route pattern instead of actual path
// This prevents high cardinality from path parameters
func getRoutePath(r *http.Request) string {
	route := mux.CurrentRoute(r)
	if route != nil {
		if path, err := route.GetPathTemplate(); err == nil {
			return path
		}
	}
	// Fallback: normalize common dynamic paths
	return normalizePath(r.URL.Path)
}

// normalizePath normalizes paths to reduce cardinality
func normalizePath(path string) string {
	// Keep known static paths as-is
	staticPaths := map[string]bool{
		"/health":           true,
		"/health/ready":     true,
		"/health/live":      true,
		"/metrics":          true,
		"/api/v1/hello":     true,
		"/api/v1/me":        true,
		"/api/v1/stats":     true,
		"/api/v1/events":    true,
		"/api/v1/protected": true,
	}

	if staticPaths[path] {
		return path
	}

	// For unknown paths, use a generic label
	return "/other"
}
