package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerMinute: 5,
		BurstSize:         2,
	}
	rl := NewRateLimiter(config)
	defer rl.Stop()

	key := "test-ip"

	for i := 0; i < 5; i++ {
		if !rl.Allow(key, 5) {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	if rl.Allow(key, 5) {
		t.Error("6th request should be denied")
	}
}

func TestRateLimiter_Remaining(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerMinute: 10,
	}
	rl := NewRateLimiter(config)
	defer rl.Stop()

	key := "test-ip"

	if remaining := rl.Remaining(key, 10); remaining != 10 {
		t.Errorf("Expected 10 remaining, got %d", remaining)
	}

	rl.Allow(key, 10)
	if remaining := rl.Remaining(key, 10); remaining != 9 {
		t.Errorf("Expected 9 remaining, got %d", remaining)
	}
}

func TestRateLimiter_DifferentKeys(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerMinute: 2,
	}
	rl := NewRateLimiter(config)
	defer rl.Stop()

	rl.Allow("ip1", 2)
	rl.Allow("ip1", 2)

	if rl.Allow("ip1", 2) {
		t.Error("ip1 should be rate limited")
	}

	if !rl.Allow("ip2", 2) {
		t.Error("ip2 should not be rate limited")
	}
}

func TestRateLimitMiddleware_SkipPaths(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerMinute: 1,
		SkipPaths:         []string{"/health"},
	}
	middleware := NewRateLimitMiddleware(config)
	defer middleware.Stop()

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected 429, got %d", w.Code)
	}

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/health", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Health check %d should not be rate limited, got %d", i, w.Code)
		}
	}
}

func TestRateLimitMiddleware_Headers(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerMinute: 10,
	}
	middleware := NewRateLimitMiddleware(config)
	defer middleware.Stop()

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Header().Get("X-RateLimit-Limit") == "" {
		t.Error("X-RateLimit-Limit header should be set")
	}
	if w.Header().Get("X-RateLimit-Remaining") == "" {
		t.Error("X-RateLimit-Remaining header should be set")
	}
	if w.Header().Get("X-RateLimit-Reset") == "" {
		t.Error("X-RateLimit-Reset header should be set")
	}
}

// TestRateLimitMiddleware_AuthTier verifies that a request carrying a
// Better Auth session cookie uses the higher AuthRequestsPerMinute
// tier and gets its own bucket (keyed by cookie, not IP). The
// anonymous tier on the same IP is exhausted after `RequestsPerMinute`
// requests, but the authenticated session keeps going.
func TestRateLimitMiddleware_AuthTier(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerMinute:     2,
		AuthRequestsPerMinute: 50,
	}
	mw := NewRateLimitMiddleware(config)
	defer mw.Stop()

	handler := mw.Handler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	authedReq := func() *http.Request {
		r := httptest.NewRequest("GET", "/api/v1/me", nil)
		r.RemoteAddr = "192.168.1.1:12345"
		r.AddCookie(&http.Cookie{Name: "better-auth.session_token", Value: "abc123"})
		return r
	}

	// Authenticated user blows past the anonymous limit on the same IP.
	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, authedReq())
		if w.Code != http.StatusOK {
			t.Fatalf("auth request %d expected 200, got %d", i+1, w.Code)
		}
	}

	// Anonymous traffic from the *same IP* still falls under the strict
	// tier, because the auth bucket is keyed by cookie, not IP.
	anonReq := func() *http.Request {
		r := httptest.NewRequest("GET", "/api/v1/hello", nil)
		r.RemoteAddr = "192.168.1.1:12345"
		return r
	}
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, anonReq())
		if w.Code != http.StatusOK {
			t.Fatalf("anon request %d expected 200, got %d", i+1, w.Code)
		}
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, anonReq())
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("3rd anonymous request from same IP should be 429, got %d", w.Code)
	}
}

// TestRateLimitMiddleware_AuthTierBearer covers the Bearer-token path
// alongside the cookie path.
func TestRateLimitMiddleware_AuthTierBearer(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerMinute:     1,
		AuthRequestsPerMinute: 5,
	}
	mw := NewRateLimitMiddleware(config)
	defer mw.Stop()

	handler := mw.Handler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := func() *http.Request {
		r := httptest.NewRequest("GET", "/api", nil)
		r.RemoteAddr = "10.0.0.5:9999"
		r.Header.Set("Authorization", "Bearer some-jwt-token")
		return r
	}

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req())
		if w.Code != http.StatusOK {
			t.Fatalf("bearer request %d expected 200, got %d", i+1, w.Code)
		}
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req())
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("6th bearer request should be 429, got %d", w.Code)
	}
}

func TestFormatInt(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{10, "10"},
		{123, "123"},
		{-1, "-1"},
		{-123, "-123"},
	}

	for _, tt := range tests {
		result := formatInt(tt.input)
		if result != tt.expected {
			t.Errorf("formatInt(%d) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestRateLimiter_Reset(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerMinute: 1,
	}
	rl := NewRateLimiter(config)
	defer rl.Stop()

	key := "test-ip"

	rl.Allow(key, 1)
	if rl.Allow(key, 1) {
		t.Error("Should be rate limited")
	}

	resetTime := rl.ResetTime(key)
	if resetTime.Before(time.Now()) {
		t.Error("Reset time should be in the future")
	}
}
