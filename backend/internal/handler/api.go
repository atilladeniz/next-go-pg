package handler

import (
	"encoding/json"
	"net/http"

	"github.com/atilladeniz/next-go-pg/backend/internal/middleware"
	"github.com/atilladeniz/next-go-pg/backend/internal/repository"
	"github.com/atilladeniz/next-go-pg/backend/internal/sse"
)

type APIHandler struct {
	authMiddleware  *middleware.AuthMiddleware
	sseBroker       *sse.Broker
	statsRepository *repository.UserStatsRepository
}

func NewAPIHandler(betterAuthURL string, broker *sse.Broker, statsRepo *repository.UserStatsRepository) *APIHandler {
	return &APIHandler{
		authMiddleware:  middleware.NewAuthMiddleware(betterAuthURL),
		sseBroker:       broker,
		statsRepository: statsRepo,
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

// UserStatsResponse represents user statistics
type UserStatsResponse struct {
	UserID        string `json:"userId"`
	ProjectCount  int    `json:"projectCount"`
	ActivityToday int    `json:"activityToday"`
	Notifications int    `json:"notifications"`
	LastLogin     string `json:"lastLogin"`
	MemberSince   string `json:"memberSince"`
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

// GetUserStats godoc
// @Summary Get user statistics
// @Description Get statistics for the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} UserStatsResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /stats [get]
func (h *APIHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	// Return default stats if database is not available
	if h.statsRepository == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(UserStatsResponse{
			UserID:        user.ID,
			ProjectCount:  3,
			ActivityToday: 10,
			Notifications: 2,
			LastLogin:     "2025-12-03T22:00:00Z",
			MemberSince:   "2025-11-15T10:00:00Z",
		})
		return
	}

	stats, err := h.statsRepository.GetOrCreate(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to get stats"})
		return
	}

	response := UserStatsResponse{
		UserID:        stats.UserID,
		ProjectCount:  stats.ProjectCount,
		ActivityToday: stats.ActivityToday,
		Notifications: stats.Notifications,
		LastLogin:     stats.LastLogin.Format("2006-01-02T15:04:05Z"),
		MemberSince:   stats.MemberSince.Format("2006-01-02T15:04:05Z"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateStatRequest for modifying stats
type UpdateStatRequest struct {
	Field string `json:"field" example:"projects"` // "projects", "activity", "notifications"
	Delta int    `json:"delta" example:"1"`        // +1 or -1
}

// UpdateUserStats godoc
// @Summary Update user statistics
// @Description Modify a stat value and broadcast the change via SSE
// @Tags users
// @Accept json
// @Produce json
// @Param request body UpdateStatRequest true "Stat update request"
// @Success 200 {object} UserStatsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /stats [post]
func (h *APIHandler) UpdateUserStats(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "unauthorized"})
		return
	}

	var req UpdateStatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request"})
		return
	}

	// Return error if database is not available
	if h.statsRepository == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "database not available"})
		return
	}

	stats, err := h.statsRepository.IncrementField(user.ID, req.Field, req.Delta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to update stats"})
		return
	}

	// Broadcast update to all connected clients
	if h.sseBroker != nil {
		h.sseBroker.Broadcast("stats-updated", `{"field":"`+req.Field+`"}`)
	}

	response := UserStatsResponse{
		UserID:        stats.UserID,
		ProjectCount:  stats.ProjectCount,
		ActivityToday: stats.ActivityToday,
		Notifications: stats.Notifications,
		LastLogin:     stats.LastLogin.Format("2006-01-02T15:04:05Z"),
		MemberSince:   stats.MemberSince.Format("2006-01-02T15:04:05Z"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAuthMiddleware returns the auth middleware for use in routes
func (h *APIHandler) GetAuthMiddleware() *middleware.AuthMiddleware {
	return h.authMiddleware
}
