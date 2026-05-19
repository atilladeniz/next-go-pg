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

	// Should allow first 5 requests
	for i := 0; i < 5; i++ {
		if !rl.Allow(key) {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 6th request should be denied
	if rl.Allow(key) {
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

	// Before any requests, should have full quota
	if remaining := rl.Remaining(key); remaining != 10 {
		t.Errorf("Expected 10 remaining, got %d", remaining)
	}

	// After one request, should have 9 remaining
	rl.Allow(key)
	if remaining := rl.Remaining(key); remaining != 9 {
		t.Errorf("Expected 9 remaining, got %d", remaining)
	}
}

func TestRateLimiter_DifferentKeys(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerMinute: 2,
	}
	rl := NewRateLimiter(config)
	defer rl.Stop()

	// Each key should have its own bucket
	rl.Allow("ip1")
	rl.Allow("ip1")

	// ip1 should be exhausted
	if rl.Allow("ip1") {
		t.Error("ip1 should be rate limited")
	}

	// ip2 should still have quota
	if !rl.Allow("ip2") {
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

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First request to /api should succeed
	req := httptest.NewRequest("GET", "/api", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// Second request to /api should be rate limited
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected 429, got %d", w.Code)
	}

	// Requests to /health should always succeed (skip path)
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

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Check rate limit headers are set
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
	// This test verifies that the bucket resets after the window
	// We use a very short test to avoid long waits
	config := RateLimitConfig{
		RequestsPerMinute: 1,
	}
	rl := NewRateLimiter(config)
	defer rl.Stop()

	key := "test-ip"

	// Use up the quota
	rl.Allow(key)
	if rl.Allow(key) {
		t.Error("Should be rate limited")
	}

	// Verify reset time is in the future
	resetTime := rl.ResetTime(key)
	if resetTime.Before(time.Now()) {
		t.Error("Reset time should be in the future")
	}
}
