// Package http is the auth bounded context's inbound interface.
// Today: read-only endpoints that surface the authenticated user.
package http

import (
	"encoding/json"
	"net/http"

	"github.com/atilladeniz/next-go-pg/backend/internal/platform/middleware"
)

// Handler exposes identity endpoints. The auth context does not own
// the authentication flow — that's Better Auth on the frontend — only
// the read-side of "who is the current user".
type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

type UserResponse struct {
	ID    string `json:"id" example:"user_123"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

type MessageResponse struct {
	Message string `json:"message" example:"Hello World"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"unauthorized"`
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description Get the currently authenticated user's information
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /me [get]
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	writeJSON(w, UserResponse{ID: user.ID, Name: user.Name, Email: user.Email})
}

// PublicHello godoc
// @Summary Public hello endpoint
// @Description Returns a hello message, no auth required
// @Tags public
// @Accept json
// @Produce json
// @Success 200 {object} MessageResponse
// @Router /hello [get]
func (h *Handler) PublicHello(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, MessageResponse{Message: "Hello from Go API!"})
}

// ProtectedHello godoc
// @Summary Protected hello endpoint
// @Description Returns a personalized hello message for authenticated users
// @Tags protected
// @Accept json
// @Produce json
// @Success 200 {object} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /protected/hello [get]
func (h *Handler) ProtectedHello(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	writeJSON(w, MessageResponse{Message: "Hello " + user.Name + "! You are authenticated."})
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	writeJSON(w, ErrorResponse{Error: message})
}
