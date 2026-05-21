package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
)

// RateLimiter implements a token bucket rate limiter. Each bucket
// remembers the limit it was created with so that mixing tiers (e.g.
// anonymous vs authenticated callers) on the same instance keeps
// refilling at the correct rate after the window rolls over.
type RateLimiter struct {
	mu       sync.RWMutex
	buckets  map[string]*bucket
	window   time.Duration
	cleanup  time.Duration
	stopChan chan struct{}
}

type bucket struct {
	tokens    int
	limit     int
	lastReset time.Time
}

// RateLimitConfig holds rate limiter configuration.
//
// The middleware applies two tiers:
//   - Anonymous callers (no session cookie / bearer token) get
//     RequestsPerMinute per IP. This is the strict tier that protects
//     the API surface from direct hammering.
//   - Authenticated callers get AuthRequestsPerMinute per session.
//     A logged-in user refreshing the dashboard fans out into many
//     parallel queries (SSR prefetch + React Query refetch + SSE
//     reconnect) — the strict tier would trip on a single rapid
//     refresh, which is bad UX. The higher tier keeps real abuse
//     bounded while making honest refreshes invisible.
type RateLimitConfig struct {
	RequestsPerMinute     int
	AuthRequestsPerMinute int
	BurstSize             int
	SkipPaths             []string
}

// DefaultRateLimitConfig returns sensible defaults.
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerMinute:     60,   // anonymous: 1 req/s on average
		AuthRequestsPerMinute: 1200, // authenticated: 20 req/s — plenty for rapid refresh + SSR fanout
		BurstSize:             10,
		SkipPaths:             []string{"/health", "/health/ready", "/health/live"},
	}
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		buckets:  make(map[string]*bucket),
		window:   time.Minute,
		cleanup:  5 * time.Minute,
		stopChan: make(chan struct{}),
	}
	go rl.cleanupLoop()
	return rl
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			for key, b := range rl.buckets {
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

// Stop stops the cleanup goroutine.
func (rl *RateLimiter) Stop() {
	close(rl.stopChan)
}

// Allow checks if a request keyed by `key` is allowed under `limit`
// requests per window. The first call for a key creates a bucket sized
// to `limit`; subsequent calls in the same window use whatever limit
// the bucket was originally created with (so a session that flips
// tiers mid-window keeps its existing budget until the window rolls).
func (rl *RateLimiter) Allow(key string, limit int) bool {
	if limit <= 0 {
		return true
	}
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, exists := rl.buckets[key]

	if !exists {
		rl.buckets[key] = &bucket{tokens: limit - 1, limit: limit, lastReset: now}
		return true
	}

	if now.Sub(b.lastReset) >= rl.window {
		b.tokens = limit - 1
		b.limit = limit
		b.lastReset = now
		return true
	}

	if b.tokens > 0 {
		b.tokens--
		return true
	}

	return false
}

// Remaining returns the number of remaining requests for a key.
// Returns `defaultLimit` if the key has no bucket yet (so the response
// headers stay sensible on the very first request of a session).
func (rl *RateLimiter) Remaining(key string, defaultLimit int) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	b, exists := rl.buckets[key]
	if !exists {
		return defaultLimit
	}
	if time.Since(b.lastReset) >= rl.window {
		return b.limit
	}
	return b.tokens
}

// ResetTime returns when the rate limit will reset for a key.
func (rl *RateLimiter) ResetTime(key string) time.Time {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	b, exists := rl.buckets[key]
	if !exists {
		return time.Now().Add(rl.window)
	}
	return b.lastReset.Add(rl.window)
}

// RateLimitMiddleware wraps the rate limiter for use as middleware.
type RateLimitMiddleware struct {
	limiter   *RateLimiter
	config    RateLimitConfig
	skipPaths map[string]bool
}

// NewRateLimitMiddleware creates a new rate limit middleware.
func NewRateLimitMiddleware(config RateLimitConfig) *RateLimitMiddleware {
	if config.AuthRequestsPerMinute <= 0 {
		// Backward-compatible fallback: callers that don't set the auth tier
		// keep the old single-tier behavior.
		config.AuthRequestsPerMinute = config.RequestsPerMinute
	}
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

// classifyRequest picks the rate-limit bucket key and effective limit
// for a request. The check is intentionally cheap — we don't validate
// the session here (the auth middleware downstream does that), we just
// look for the *presence* of a credential. The downside is that a
// stolen cookie buys the attacker the higher tier; the upside is that
// honest logged-in users get a tier that doesn't trip on rapid
// refresh. The auth middleware still rejects invalid credentials, so
// rate-limit-only protection is not load-bearing for security.
func (m *RateLimitMiddleware) classifyRequest(r *http.Request) (key string, limit int, authed bool) {
	if token := bearerToken(r); token != "" {
		return "auth:" + hashCred(token), m.config.AuthRequestsPerMinute, true
	}
	if cookie := sessionCookieValue(r); cookie != "" {
		return "auth:" + hashCred(cookie), m.config.AuthRequestsPerMinute, true
	}
	return "ip:" + getClientIP(r), m.config.RequestsPerMinute, false
}

func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if h == "" {
		return ""
	}
	const prefix = "Bearer "
	if len(h) <= len(prefix) || !strings.EqualFold(h[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(h[len(prefix):])
}

// sessionCookieValue returns the value of the Better Auth session
// cookie if any of its variants are present. Better Auth defaults to
// `better-auth.session_token` and adds a `__Secure-` prefix when the
// cookie is set over HTTPS — we accept both.
func sessionCookieValue(r *http.Request) string {
	for _, c := range r.Cookies() {
		name := c.Name
		if strings.Contains(name, "session_token") || strings.HasPrefix(name, "better-auth.") || strings.HasPrefix(name, "__Secure-better-auth.") {
			if c.Value != "" {
				return c.Value
			}
		}
	}
	return ""
}

func hashCred(v string) string {
	sum := sha256.Sum256([]byte(v))
	return hex.EncodeToString(sum[:8])
}

// Handler returns the middleware handler.
func (m *RateLimitMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.skipPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		key, limit, authed := m.classifyRequest(r)

		if !m.limiter.Allow(key, limit) {
			remaining := m.limiter.Remaining(key, limit)
			resetTime := m.limiter.ResetTime(key)

			w.Header().Set("X-RateLimit-Limit", formatInt(limit))
			w.Header().Set("X-RateLimit-Remaining", formatInt(remaining))
			w.Header().Set("X-RateLimit-Reset", resetTime.Format(time.RFC3339))
			w.Header().Set("Retry-After", formatInt(int(time.Until(resetTime).Seconds())))
			w.Header().Set("Content-Type", "application/json")

			logger.Warn().
				Str("client_ip", getClientIP(r)).
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Bool("authenticated", authed).
				Int("limit", limit).
				Msg("Rate limit exceeded")

			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"error":       "rate limit exceeded",
				"retry_after": int(time.Until(resetTime).Seconds()),
			})
			return
		}

		remaining := m.limiter.Remaining(key, limit)
		resetTime := m.limiter.ResetTime(key)

		w.Header().Set("X-RateLimit-Limit", formatInt(limit))
		w.Header().Set("X-RateLimit-Remaining", formatInt(remaining))
		w.Header().Set("X-RateLimit-Reset", resetTime.Format(time.RFC3339))

		next.ServeHTTP(w, r)
	})
}

// formatInt converts an int to a string (no fmt to avoid an alloc).
func formatInt(n int) string {
	if n == 0 {
		return "0"
	}

	negative := false
	if n < 0 {
		negative = true
		n = -n
	}

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

// Stop stops the rate limiter cleanup goroutine.
func (m *RateLimitMiddleware) Stop() {
	m.limiter.Stop()
}
