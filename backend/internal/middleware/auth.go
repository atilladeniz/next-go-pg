package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type User struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"emailVerified"`
	Image         string    `json:"image"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type SessionResponse struct {
	Session *Session `json:"session"`
	User    *User    `json:"user"`
}

type contextKey string

const (
	UserContextKey    contextKey = "user"
	SessionContextKey contextKey = "session"
)

type AuthMiddleware struct {
	betterAuthURL string
	httpClient    *http.Client
}

func NewAuthMiddleware(betterAuthURL string) *AuthMiddleware {
	return &AuthMiddleware{
		betterAuthURL: strings.TrimSuffix(betterAuthURL, "/"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// RequireAuth middleware validates the session with Better Auth
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for preflight OPTIONS requests (handled by CORS middleware)
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		session, user, err := m.validateSession(r)
		if err != nil {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// Add session and user to context
		ctx := context.WithValue(r.Context(), SessionContextKey, session)
		ctx = context.WithValue(ctx, UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth middleware adds user to context if authenticated, but doesn't require it
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, user, err := m.validateSession(r)
		if err == nil && session != nil {
			ctx := context.WithValue(r.Context(), SessionContextKey, session)
			ctx = context.WithValue(ctx, UserContextKey, user)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) validateSession(r *http.Request) (*Session, *User, error) {
	// Forward cookies to Better Auth
	req, err := http.NewRequest("GET", m.betterAuthURL+"/api/auth/get-session", nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Copy cookies from original request
	for _, cookie := range r.Cookies() {
		req.AddCookie(cookie)
	}

	// Also check Authorization header
	if auth := r.Header.Get("Authorization"); auth != "" {
		req.Header.Set("Authorization", auth)
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to validate session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("invalid session: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response: %w", err)
	}

	var sessionResp SessionResponse
	if err := json.Unmarshal(body, &sessionResp); err != nil {
		return nil, nil, fmt.Errorf("failed to parse session: %w", err)
	}

	if sessionResp.Session == nil {
		return nil, nil, fmt.Errorf("no session found")
	}

	return sessionResp.Session, sessionResp.User, nil
}

// GetUserFromContext retrieves the user from context
func GetUserFromContext(ctx context.Context) *User {
	if user, ok := ctx.Value(UserContextKey).(*User); ok {
		return user
	}
	return nil
}

// GetSessionFromContext retrieves the session from context
func GetSessionFromContext(ctx context.Context) *Session {
	if session, ok := ctx.Value(SessionContextKey).(*Session); ok {
		return session
	}
	return nil
}
