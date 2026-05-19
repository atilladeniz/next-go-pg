package middleware

import (
	"net/http"

	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
)

// CombinedAuthMiddleware tries JWT first, then falls back to Better Auth session validation
type CombinedAuthMiddleware struct {
	jwt  *JWTMiddleware
	auth *AuthMiddleware
}

// NewCombinedAuthMiddleware creates a combined auth middleware
func NewCombinedAuthMiddleware(betterAuthURL string) *CombinedAuthMiddleware {
	return &CombinedAuthMiddleware{
		jwt:  NewJWTMiddleware(),
		auth: NewAuthMiddleware(betterAuthURL),
	}
}

// RequireAuth validates authentication - tries JWT first, then Better Auth
// This provides the best of both worlds:
// - JWT: Fast, local validation (no network call)
// - Better Auth: Fallback for clients without JWT cookie
func (m *CombinedAuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip OPTIONS requests
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// Try JWT validation first (fast path)
		claims, err := m.jwt.validateToken(r)
		if err == nil && claims != nil {
			logger.Debug().
				Str("user_id", claims.Subject).
				Msg("Authenticated via JWT")

			// Add to context (same as RequireJWT)
			ctx := r.Context()
			ctx = withJWTContext(ctx, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Fall back to Better Auth session validation (slow path)
		session, user, err := m.auth.validateSession(r)
		if err != nil {
			logger.Debug().
				Err(err).
				Msg("Both JWT and Better Auth validation failed")
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}

		logger.Debug().
			Str("user_id", user.ID).
			Msg("Authenticated via Better Auth")

		// Add session and user to context
		ctx := r.Context()
		ctx = withSessionContext(ctx, session)
		ctx = withUserContext(ctx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth adds user to context if authenticated, but doesn't require it
func (m *CombinedAuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Try JWT first
		claims, err := m.jwt.validateToken(r)
		if err == nil && claims != nil {
			ctx = withJWTContext(ctx, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Try Better Auth
		session, user, err := m.auth.validateSession(r)
		if err == nil && session != nil {
			ctx = withSessionContext(ctx, session)
			ctx = withUserContext(ctx, user)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// JWTOnly returns a middleware that only uses JWT validation
func (m *CombinedAuthMiddleware) JWTOnly() *JWTMiddleware {
	return m.jwt
}

// BetterAuthOnly returns a middleware that only uses Better Auth validation
func (m *CombinedAuthMiddleware) BetterAuthOnly() *AuthMiddleware {
	return m.auth
}
