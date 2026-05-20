package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	mu       sync.RWMutex
	buckets  map[string]*bucket
	rate     int           // requests per window
	window   time.Duration // time window
	cleanup  time.Duration // cleanup interval for old buckets
	stopChan chan struct{}
}

type bucket struct {
	tokens    int
	lastReset time.Time
}

// RateLimitConfig holds rate limiter configuration
type RateLimitConfig struct {
	// RequestsPerMinute is the number of requests allowed per minute per IP
	RequestsPerMinute int
	// BurstSize allows temporary bursts above the limit
	BurstSize int
	// SkipPaths are paths that bypass rate limiting (e.g., health checks)
	SkipPaths []string
}

// DefaultRateLimitConfig returns sensible defaults
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerMinute: 60, // 1 request per second on average
		BurstSize:         10, // Allow bursts of 10 requests
		SkipPaths:         []string{"/health", "/health/ready", "/health/live"},
	}
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		buckets:  make(map[string]*bucket),
		rate:     config.RequestsPerMinute,
		window:   time.Minute,
		cleanup:  5 * time.Minute,
		stopChan: make(chan struct{}),
	}

	// Start cleanup goroutine
	go rl.cleanupLoop()

	return rl
}

// cleanupLoop removes old buckets periodically
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			for key, b := range rl.buckets {
				// Remove buckets that haven't been used for 2x the window
				if now.Sub(b.lastReset) > 2*rl.window {
					delete(rl.buckets, key)
				}
			}
			rl.mu.Unlock()
		case <-rl.stopChan:
			return
		}
	}
}

// Stop stops the cleanup goroutine
func (rl *RateLimiter) Stop() {
	close(rl.stopChan)
}

// Allow checks if a request from the given key should be allowed
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, exists := rl.buckets[key]

	if !exists {
		// Create new bucket with full tokens
		rl.buckets[key] = &bucket{
			tokens:    rl.rate - 1, // Use one token for this request
			lastReset: now,
		}
		return true
	}

	// Check if we need to reset the bucket
	if now.Sub(b.lastReset) >= rl.window {
		b.tokens = rl.rate - 1 // Reset and use one token
		b.lastReset = now
		return true
	}

	// Check if we have tokens available
	if b.tokens > 0 {
		b.tokens--
		return true
	}

	return false
}

// Remaining returns the number of remaining requests for a key
func (rl *RateLimiter) Remaining(key string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	b, exists := rl.buckets[key]
	if !exists {
		return rl.rate
	}

	// Check if bucket should be reset
	if time.Since(b.lastReset) >= rl.window {
		return rl.rate
	}

	return b.tokens
}

// ResetTime returns when the rate limit will reset for a key
func (rl *RateLimiter) ResetTime(key string) time.Time {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	b, exists := rl.buckets[key]
	if !exists {
		return time.Now().Add(rl.window)
	}

	return b.lastReset.Add(rl.window)
}

// RateLimitMiddleware wraps the rate limiter for use as middleware
type RateLimitMiddleware struct {
	limiter   *RateLimiter
	config    RateLimitConfig
	skipPaths map[string]bool
}

// NewRateLimitMiddleware creates a new rate limit middleware
func NewRateLimitMiddleware(config RateLimitConfig) *RateLimitMiddleware {
	skipPaths := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipPaths[path] = true
	}

	return &RateLimitMiddleware{
		limiter:   NewRateLimiter(config),
		config:    config,
		skipPaths: skipPaths,
	}
}

// Handler returns the middleware handler
func (m *RateLimitMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip rate limiting for configured paths
		if m.skipPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		// Get client identifier (IP address)
		clientIP := getClientIP(r)

		// Check rate limit
		if !m.limiter.Allow(clientIP) {
			remaining := m.limiter.Remaining(clientIP)
			resetTime := m.limiter.ResetTime(clientIP)

			// Set rate limit headers
			w.Header().Set("X-RateLimit-Limit", string(rune(m.config.RequestsPerMinute)))
			w.Header().Set("X-RateLimit-Remaining", string(rune(remaining)))
			w.Header().Set("X-RateLimit-Reset", resetTime.Format(time.RFC3339))
			w.Header().Set("Retry-After", string(rune(int(time.Until(resetTime).Seconds()))))
			w.Header().Set("Content-Type", "application/json")

			logger.Warn().
				Str("client_ip", clientIP).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Msg("Rate limit exceeded")

			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":       "rate limit exceeded",
				"retry_after": int(time.Until(resetTime).Seconds()),
			})
			return
		}

		// Set rate limit headers for successful requests
		remaining := m.limiter.Remaining(clientIP)
		resetTime := m.limiter.ResetTime(clientIP)

		w.Header().Set("X-RateLimit-Limit", formatInt(m.config.RequestsPerMinute))
		w.Header().Set("X-RateLimit-Remaining", formatInt(remaining))
		w.Header().Set("X-RateLimit-Reset", resetTime.Format(time.RFC3339))

		next.ServeHTTP(w, r)
	})
}

// formatInt converts an int to a string
func formatInt(n int) string {
	if n == 0 {
		return "0"
	}

	// Handle negative numbers
	negative := false
	if n < 0 {
		negative = true
		n = -n
	}

	// Build the string in reverse
	var digits [20]byte
	i := len(digits)
	for n > 0 {
		i--
		digits[i] = byte('0' + n%10)
		n /= 10
	}

	if negative {
		i--
		digits[i] = '-'
	}

	return string(digits[i:])
}

// Stop stops the rate limiter cleanup goroutine
func (m *RateLimitMiddleware) Stop() {
	m.limiter.Stop()
}
