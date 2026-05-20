package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-jwt-secret-for-unit-tests"

// createTestToken creates a valid JWT token for testing
func createTestToken(t *testing.T, claims *JWTClaims, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	require.NoError(t, err)
	return tokenString
}

// TestNewJWTMiddleware tests middleware initialization
func TestNewJWTMiddleware(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		secret      string
		shouldPanic bool
	}{
		{
			name:        "development without secret - uses default",
			environment: "development",
			secret:      "",
			shouldPanic: false,
		},
		{
			name:        "development with secret",
			environment: "development",
			secret:      "my-secret",
			shouldPanic: false,
		},
		{
			name:        "staging with secret",
			environment: "staging",
			secret:      "my-secret",
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("ENVIRONMENT", tt.environment)
			if tt.secret != "" {
				os.Setenv("GO_JWT_SECRET", tt.secret)
			} else {
				os.Unsetenv("GO_JWT_SECRET")
			}
			defer func() {
				os.Unsetenv("ENVIRONMENT")
				os.Unsetenv("GO_JWT_SECRET")
			}()

			if tt.shouldPanic {
				assert.Panics(t, func() {
					NewJWTMiddleware()
				})
			} else {
				middleware := NewJWTMiddleware()
				assert.NotNil(t, middleware)
				assert.Equal(t, "go-auth-token", middleware.cookieName)
			}
		})
	}
}

// TestRequireJWT tests the RequireJWT middleware
func TestRequireJWT(t *testing.T) {
	os.Setenv("GO_JWT_SECRET", testSecret)
	defer os.Unsetenv("GO_JWT_SECRET")

	middleware := NewJWTMiddleware()

	// Create valid claims
	validClaims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Email: "test@example.com",
		Name:  "Test User",
		SID:   "session-456",
	}

	tests := []struct {
		name           string
		setupRequest   func(*http.Request)
		expectedStatus int
		expectUser     bool
	}{
		{
			name: "valid token in cookie",
			setupRequest: func(r *http.Request) {
				token := createTestToken(t, validClaims, testSecret)
				r.AddCookie(&http.Cookie{
					Name:  "go-auth-token",
					Value: token,
				})
			},
			expectedStatus: http.StatusOK,
			expectUser:     true,
		},
		{
			name: "valid token in Authorization header",
			setupRequest: func(r *http.Request) {
				token := createTestToken(t, validClaims, testSecret)
				r.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusOK,
			expectUser:     true,
		},
		{
			name: "missing token",
			setupRequest: func(r *http.Request) {
				// No token
			},
			expectedStatus: http.StatusUnauthorized,
			expectUser:     false,
		},
		{
			name: "invalid token signature",
			setupRequest: func(r *http.Request) {
				token := createTestToken(t, validClaims, "wrong-secret")
				r.AddCookie(&http.Cookie{
					Name:  "go-auth-token",
					Value: token,
				})
			},
			expectedStatus: http.StatusUnauthorized,
			expectUser:     false,
		},
		{
			name: "expired token",
			setupRequest: func(r *http.Request) {
				expiredClaims := &JWTClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						Subject:   "user-123",
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
					},
					Email: "test@example.com",
					SID:   "session-456",
				}
				token := createTestToken(t, expiredClaims, testSecret)
				r.AddCookie(&http.Cookie{
					Name:  "go-auth-token",
					Value: token,
				})
			},
			expectedStatus: http.StatusUnauthorized,
			expectUser:     false,
		},
		{
			name: "malformed token",
			setupRequest: func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:  "go-auth-token",
					Value: "not.a.valid.token",
				})
			},
			expectedStatus: http.StatusUnauthorized,
			expectUser:     false,
		},
		{
			name: "OPTIONS request - skip auth",
			setupRequest: func(r *http.Request) {
				// OPTIONS without token should pass
			},
			expectedStatus: http.StatusOK,
			expectUser:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedUser *User
			var capturedSession *Session

			// Create test handler that captures context values
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedUser = GetUserFromContext(r.Context())
				capturedSession = GetSessionFromContext(r.Context())
				w.WriteHeader(http.StatusOK)
			})

			// Wrap with middleware
			handler := middleware.RequireJWT(testHandler)

			// Create request
			method := http.MethodGet
			if tt.name == "OPTIONS request - skip auth" {
				method = http.MethodOptions
			}
			req := httptest.NewRequest(method, "/test", nil)
			tt.setupRequest(req)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectUser {
				assert.NotNil(t, capturedUser)
				assert.Equal(t, "user-123", capturedUser.ID)
				assert.Equal(t, "test@example.com", capturedUser.Email)
				assert.Equal(t, "Test User", capturedUser.Name)

				assert.NotNil(t, capturedSession)
				assert.Equal(t, "session-456", capturedSession.ID)
				assert.Equal(t, "user-123", capturedSession.UserID)
			}
		})
	}
}

// TestOptionalJWT tests the OptionalJWT middleware
func TestOptionalJWT(t *testing.T) {
	os.Setenv("GO_JWT_SECRET", testSecret)
	defer os.Unsetenv("GO_JWT_SECRET")

	middleware := NewJWTMiddleware()

	validClaims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Email: "test@example.com",
		Name:  "Test User",
		SID:   "session-456",
	}

	tests := []struct {
		name           string
		setupRequest   func(*http.Request)
		expectedStatus int
		expectUser     bool
	}{
		{
			name: "valid token - adds user to context",
			setupRequest: func(r *http.Request) {
				token := createTestToken(t, validClaims, testSecret)
				r.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusOK,
			expectUser:     true,
		},
		{
			name: "no token - request continues without user",
			setupRequest: func(r *http.Request) {
				// No token
			},
			expectedStatus: http.StatusOK,
			expectUser:     false,
		},
		{
			name: "invalid token - request continues without user",
			setupRequest: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer invalid-token")
			},
			expectedStatus: http.StatusOK,
			expectUser:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedUser *User

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedUser = GetUserFromContext(r.Context())
				w.WriteHeader(http.StatusOK)
			})

			handler := middleware.OptionalJWT(testHandler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			tt.setupRequest(req)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectUser {
				assert.NotNil(t, capturedUser)
				assert.Equal(t, "user-123", capturedUser.ID)
			} else {
				assert.Nil(t, capturedUser)
			}
		})
	}
}

// TestGetJWTClaimsFromContext tests context retrieval
func TestGetJWTClaimsFromContext(t *testing.T) {
	t.Run("claims present in context", func(t *testing.T) {
		claims := &JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{Subject: "user-123"},
			Email:            "test@example.com",
			SID:              "session-456",
		}
		ctx := context.WithValue(context.Background(), JWTContextKey, claims)

		result := GetJWTClaimsFromContext(ctx)
		assert.NotNil(t, result)
		assert.Equal(t, "user-123", result.Subject)
		assert.Equal(t, "test@example.com", result.Email)
	})

	t.Run("no claims in context", func(t *testing.T) {
		ctx := context.Background()
		result := GetJWTClaimsFromContext(ctx)
		assert.Nil(t, result)
	})

	t.Run("wrong type in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), JWTContextKey, "not-claims")
		result := GetJWTClaimsFromContext(ctx)
		assert.Nil(t, result)
	})
}

// TestJWTValidateToken tests token validation directly
func TestJWTValidateToken(t *testing.T) {
	os.Setenv("GO_JWT_SECRET", testSecret)
	defer os.Unsetenv("GO_JWT_SECRET")

	middleware := NewJWTMiddleware()

	t.Run("cookie takes precedence over header", func(t *testing.T) {
		cookieClaims := &JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   "cookie-user",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
			Email: "cookie@example.com",
			SID:   "cookie-session",
		}
		headerClaims := &JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   "header-user",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
			Email: "header@example.com",
			SID:   "header-session",
		}

		cookieToken := createTestToken(t, cookieClaims, testSecret)
		headerToken := createTestToken(t, headerClaims, testSecret)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{Name: "go-auth-token", Value: cookieToken})
		req.Header.Set("Authorization", "Bearer "+headerToken)

		claims, err := middleware.validateToken(req)
		require.NoError(t, err)
		assert.Equal(t, "cookie-user", claims.Subject)
	})

	t.Run("header used when no cookie", func(t *testing.T) {
		headerClaims := &JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   "header-user",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
			Email: "header@example.com",
			SID:   "header-session",
		}

		headerToken := createTestToken(t, headerClaims, testSecret)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+headerToken)

		claims, err := middleware.validateToken(req)
		require.NoError(t, err)
		assert.Equal(t, "header-user", claims.Subject)
	})

	t.Run("invalid signing method rejected", func(t *testing.T) {
		// Create token with different signing method (RS256 instead of HS256)
		token := jwt.NewWithClaims(jwt.SigningMethodNone, &JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{Subject: "user-123"},
		})
		tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		_, err := middleware.validateToken(req)
		assert.Error(t, err)
	})
}

// TestJWTClaims tests the claims structure
func TestJWTClaims(t *testing.T) {
	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-id-123",
			Issuer:    "test-issuer",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Email: "user@example.com",
		Name:  "Test User",
		SID:   "session-id-456",
	}

	// Create and sign token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(testSecret))
	require.NoError(t, err)

	// Parse token back
	parsedToken, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(testSecret), nil
	})
	require.NoError(t, err)

	parsedClaims, ok := parsedToken.Claims.(*JWTClaims)
	require.True(t, ok)

	assert.Equal(t, "user-id-123", parsedClaims.Subject)
	assert.Equal(t, "user@example.com", parsedClaims.Email)
	assert.Equal(t, "Test User", parsedClaims.Name)
	assert.Equal(t, "session-id-456", parsedClaims.SID)
}

// TestAuthorizationHeaderParsing tests Bearer token extraction
func TestAuthorizationHeaderParsing(t *testing.T) {
	os.Setenv("GO_JWT_SECRET", testSecret)
	defer os.Unsetenv("GO_JWT_SECRET")

	middleware := NewJWTMiddleware()

	validClaims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Email: "test@example.com",
		SID:   "session-456",
	}
	validToken := createTestToken(t, validClaims, testSecret)

	tests := []struct {
		name          string
		header        string
		expectSuccess bool
	}{
		{
			name:          "valid Bearer token",
			header:        "Bearer " + validToken,
			expectSuccess: true,
		},
		{
			name:          "missing Bearer prefix",
			header:        validToken,
			expectSuccess: false,
		},
		{
			name:          "lowercase bearer",
			header:        "bearer " + validToken,
			expectSuccess: false,
		},
		{
			name:          "empty header",
			header:        "",
			expectSuccess: false,
		},
		{
			name:          "only Bearer keyword",
			header:        "Bearer ",
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}

			claims, err := middleware.validateToken(req)
			if tt.expectSuccess {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
