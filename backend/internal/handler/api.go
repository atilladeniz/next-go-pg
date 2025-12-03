package handler

import (
	"encoding/json"
	"net/http"

	"github.com/atilla/gocatest/backend/internal/middleware"
)

// @title GocaTest API
// @version 1.0
// @description Go Clean Architecture API with Better Auth integration
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

type APIHandler struct {
	authMiddleware *middleware.AuthMiddleware
}

func NewAPIHandler(betterAuthURL string) *APIHandler {
	return &APIHandler{
		authMiddleware: middleware.NewAuthMiddleware(betterAuthURL),
	}
}

// UserResponse represents the current user
type UserResponse struct {
	ID    string `json:"id" example:"user_123"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

// MessageResponse represents a simple message
type MessageResponse struct {
	Message string `json:"message" example:"Hello World"`
}

// ErrorResponse represents an error
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
func (h *APIHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}

// PublicHello godoc
// @Summary Public hello endpoint
// @Description Returns a hello message, no auth required
// @Tags public
// @Accept json
// @Produce json
// @Success 200 {object} MessageResponse
// @Router /hello [get]
func (h *APIHandler) PublicHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MessageResponse{Message: "Hello from Go API!"})
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
func (h *APIHandler) ProtectedHello(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MessageResponse{
		Message: "Hello " + user.Name + "! You are authenticated.",
	})
}

// GetAuthMiddleware returns the auth middleware for use in routes
func (h *APIHandler) GetAuthMiddleware() *middleware.AuthMiddleware {
	return h.authMiddleware
}
