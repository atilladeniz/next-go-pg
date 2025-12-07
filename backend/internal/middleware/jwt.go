package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the claims in the Go JWT token
type JWTClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
	SID   string `json:"sid"` // Session ID
}

// JWTContextKey is the context key for JWT claims
const JWTContextKey contextKey = "jwt_claims"

// JWTMiddleware handles JWT validation for Go backend
type JWTMiddleware struct {
	secret     []byte
	cookieName string
}

// NewJWTMiddleware creates a new JWT middleware instance
func NewJWTMiddleware() *JWTMiddleware {
	secret := os.Getenv("GO_JWT_SECRET")
	env := os.Getenv("ENVIRONMENT")

	if secret == "" {
		if env == "production" {
			logger.Fatal().Msg("GO_JWT_SECRET must be set in production environment")
		}
		secret = "dev-secret-change-in-production"
		logger.Warn().Msg("GO_JWT_SECRET not set, using insecure default (development only)")
	}

	return &JWTMiddleware{
		secret:     []byte(secret),
		cookieName: "go-auth-token",
	}
}

// RequireJWT validates the JWT and adds claims to context
// Use this for endpoints that require authentication via JWT
func (m *JWTMiddleware) RequireJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip OPTIONS requests
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := m.validateToken(r)
		if err != nil {
			logger.Debug().Err(err).Msg("JWT validation failed")
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), JWTContextKey, claims)

		// Also populate User context for compatibility with existing handlers
		user := &User{
			ID:    claims.Subject,
			Email: claims.Email,
			Name:  claims.Name,
		}
		ctx = context.WithValue(ctx, UserContextKey, user)

		// Add session info
		session := &Session{
			ID:     claims.SID,
			UserID: claims.Subject,
		}
		ctx = context.WithValue(ctx, SessionContextKey, session)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalJWT adds claims to context if token is present, but doesn't require it
func (m *JWTMiddleware) OptionalJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := m.validateToken(r)
		if err == nil && claims != nil {
			ctx := context.WithValue(r.Context(), JWTContextKey, claims)

			user := &User{
				ID:    claims.Subject,
				Email: claims.Email,
				Name:  claims.Name,
			}
			ctx = context.WithValue(ctx, UserContextKey, user)

			session := &Session{
				ID:     claims.SID,
				UserID: claims.Subject,
			}
			ctx = context.WithValue(ctx, SessionContextKey, session)

			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

func (m *JWTMiddleware) validateToken(r *http.Request) (*JWTClaims, error) {
	tokenString := ""

	// First, check for cookie
	if cookie, err := r.Cookie(m.cookieName); err == nil {
		tokenString = cookie.Value
	}

	// Fall back to Authorization header
	if tokenString == "" {
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if tokenString == "" {
		return nil, jwt.ErrTokenMalformed
	}

	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

// GetJWTClaimsFromContext retrieves JWT claims from context
func GetJWTClaimsFromContext(ctx context.Context) *JWTClaims {
	if claims, ok := ctx.Value(JWTContextKey).(*JWTClaims); ok {
		return claims
	}
	return nil
}
