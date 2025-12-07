package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTPRequestsTotal counts total HTTP requests
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestDuration measures request latency
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	// HTTPRequestsInFlight tracks concurrent requests
	HTTPRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
	)

	// DatabaseQueryDuration measures database query latency
	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Database query latency in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
		},
		[]string{"operation"},
	)

	// DatabaseConnectionsActive tracks active database connections
	DatabaseConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_active",
			Help: "Number of active database connections",
		},
	)

	// DatabaseConnectionsIdle tracks idle database connections
	DatabaseConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_idle",
			Help: "Number of idle database connections",
		},
	)

	// AuthAttemptsTotal counts authentication attempts
	AuthAttemptsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_attempts_total",
			Help: "Total number of authentication attempts",
		},
		[]string{"type", "status"},
	)

	// SSEConnectionsActive tracks active SSE connections
	SSEConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "sse_connections_active",
			Help: "Number of active SSE connections",
		},
	)

	// SSEMessagesSent counts SSE messages sent
	SSEMessagesSent = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "sse_messages_sent_total",
			Help: "Total number of SSE messages sent",
		},
	)

	// RateLimitHits counts rate limit hits
	RateLimitHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "rate_limit_hits_total",
			Help: "Total number of rate limit hits (429 responses)",
		},
	)

	// AppInfo provides application metadata
	AppInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_info",
			Help: "Application information",
		},
		[]string{"version", "environment"},
	)
)

// Init initializes application info metrics
func Init(version, environment string) {
	AppInfo.WithLabelValues(version, environment).Set(1)
}
