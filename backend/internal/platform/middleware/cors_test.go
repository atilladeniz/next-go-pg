package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCORSMiddleware tests the CORS middleware
func TestCORSMiddleware(t *testing.T) {
	frontendURL := "http://localhost:3000"
	corsMiddleware := NewCORSMiddleware(frontendURL)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := corsMiddleware.Handler(testHandler)

	tests := []struct {
		name                string
		origin              string
		method              string
		expectedStatus      int
		expectOriginHeader  bool
		expectCredentials   bool
		expectMethodsHeader bool
	}{
		{
			name:                "allowed origin",
			origin:              "http://localhost:3000",
			method:              http.MethodGet,
			expectedStatus:      http.StatusOK,
			expectOriginHeader:  true,
			expectCredentials:   true,
			expectMethodsHeader: true,
		},
		{
			name:                "disallowed origin",
			origin:              "http://evil.com",
			method:              http.MethodGet,
			expectedStatus:      http.StatusOK,
			expectOriginHeader:  false,
			expectCredentials:   true, // Still set, but no origin
			expectMethodsHeader: true,
		},
		{
			name:                "preflight request",
			origin:              "http://localhost:3000",
			method:              http.MethodOptions,
			expectedStatus:      http.StatusNoContent,
			expectOriginHeader:  true,
			expectCredentials:   true,
			expectMethodsHeader: true,
		},
		{
			name:                "no origin header",
			origin:              "",
			method:              http.MethodGet,
			expectedStatus:      http.StatusOK,
			expectOriginHeader:  false,
			expectCredentials:   true,
			expectMethodsHeader: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectOriginHeader {
				assert.Equal(t, tt.origin, rr.Header().Get("Access-Control-Allow-Origin"))
			} else {
				assert.Empty(t, rr.Header().Get("Access-Control-Allow-Origin"))
			}

			if tt.expectCredentials {
				assert.Equal(t, "true", rr.Header().Get("Access-Control-Allow-Credentials"))
			}

			if tt.expectMethodsHeader {
				assert.Contains(t, rr.Header().Get("Access-Control-Allow-Methods"), "GET")
				assert.Contains(t, rr.Header().Get("Access-Control-Allow-Methods"), "POST")
			}
		})
	}
}

// TestCORSAllowedHeaders tests that webhook headers are allowed
func TestCORSAllowedHeaders(t *testing.T) {
	corsMiddleware := NewCORSMiddleware("http://localhost:3000")

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := corsMiddleware.Handler(testHandler)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/webhooks/send-magic-link", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type, X-Webhook-Secret")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	allowedHeaders := rr.Header().Get("Access-Control-Allow-Headers")

	// Verify X-Webhook-Secret is in allowed headers
	assert.Contains(t, allowedHeaders, "X-Webhook-Secret")
	assert.Contains(t, allowedHeaders, "Content-Type")
	assert.Contains(t, allowedHeaders, "Authorization")
	assert.Contains(t, allowedHeaders, "Cookie")
}

// TestCORSWildcard tests wildcard origin support
func TestCORSWildcard(t *testing.T) {
	// Create middleware with wildcard origin
	corsMiddleware := &CORSMiddleware{
		allowedOrigins: []string{"*"},
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := corsMiddleware.Handler(testHandler)

	tests := []struct {
		name   string
		origin string
	}{
		{"any origin 1", "http://example.com"},
		{"any origin 2", "http://another.example.com"},
		{"localhost", "http://localhost:3000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
			req.Header.Set("Origin", tt.origin)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// Wildcard should allow any origin
			assert.Equal(t, tt.origin, rr.Header().Get("Access-Control-Allow-Origin"))
		})
	}
}

// TestDefaultCORSConfig tests the default configuration
func TestDefaultCORSConfig(t *testing.T) {
	config := DefaultCORSConfig()

	assert.Equal(t, []string{"http://localhost:3000"}, config.AllowedOrigins)
	assert.Contains(t, config.AllowedMethods, "GET")
	assert.Contains(t, config.AllowedMethods, "POST")
	assert.Contains(t, config.AllowedMethods, "PUT")
	assert.Contains(t, config.AllowedMethods, "PATCH")
	assert.Contains(t, config.AllowedMethods, "DELETE")
	assert.Contains(t, config.AllowedMethods, "OPTIONS")
	assert.Contains(t, config.AllowedHeaders, "Content-Type")
	assert.Contains(t, config.AllowedHeaders, "Authorization")
	assert.Contains(t, config.AllowedHeaders, "X-Webhook-Secret")
	assert.True(t, config.AllowCredentials)
}

// TestCORSFunction tests the configurable CORS middleware
func TestCORSFunction(t *testing.T) {
	config := CORSConfig{
		AllowedOrigins:   []string{"http://custom.example.com", "http://another.example.com"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-Type", "X-Custom-Header"},
		AllowCredentials: true,
	}

	corsHandler := CORS(config)
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := corsHandler(testHandler)

	t.Run("first allowed origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.Header.Set("Origin", "http://custom.example.com")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, "http://custom.example.com", rr.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("second allowed origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.Header.Set("Origin", "http://another.example.com")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, "http://another.example.com", rr.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("disallowed origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.Header.Set("Origin", "http://evil.com")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Empty(t, rr.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("custom allowed headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/api/test", nil)
		req.Header.Set("Origin", "http://custom.example.com")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Contains(t, rr.Header().Get("Access-Control-Allow-Headers"), "X-Custom-Header")
	})

	t.Run("credentials disabled", func(t *testing.T) {
		noCredConfig := CORSConfig{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET"},
			AllowedHeaders:   []string{"Content-Type"},
			AllowCredentials: false,
		}
		noCredHandler := CORS(noCredConfig)(testHandler)

		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.Header.Set("Origin", "http://example.com")

		rr := httptest.NewRecorder()
		noCredHandler.ServeHTTP(rr, req)

		assert.Empty(t, rr.Header().Get("Access-Control-Allow-Credentials"))
	})
}

// TestCORSPreflightForWebhooks tests preflight for webhook endpoints
func TestCORSPreflightForWebhooks(t *testing.T) {
	corsMiddleware := NewCORSMiddleware("http://localhost:3000")

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Should not reach here for OPTIONS"))
	})

	handler := corsMiddleware.Handler(testHandler)

	webhookEndpoints := []string{
		"/api/v1/webhooks/session-created",
		"/api/v1/webhooks/send-magic-link",
		"/api/v1/webhooks/send-verification-email",
		"/api/v1/webhooks/send-2fa-otp",
		"/api/v1/webhooks/send-2fa-enabled",
		"/api/v1/webhooks/send-passkey-added",
	}

	for _, endpoint := range webhookEndpoints {
		t.Run(endpoint, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodOptions, endpoint, nil)
			req.Header.Set("Origin", "http://localhost:3000")
			req.Header.Set("Access-Control-Request-Method", "POST")
			req.Header.Set("Access-Control-Request-Headers", "Content-Type, X-Webhook-Secret")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// Preflight should return 204
			assert.Equal(t, http.StatusNoContent, rr.Code)

			// Response body should be empty for preflight
			assert.Empty(t, rr.Body.String())

			// Headers should be set correctly
			assert.Equal(t, "http://localhost:3000", rr.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, "true", rr.Header().Get("Access-Control-Allow-Credentials"))
			assert.Contains(t, rr.Header().Get("Access-Control-Allow-Methods"), "POST")
			assert.Contains(t, rr.Header().Get("Access-Control-Allow-Headers"), "X-Webhook-Secret")
		})
	}
}
