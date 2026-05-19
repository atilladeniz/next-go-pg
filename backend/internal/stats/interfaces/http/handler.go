// Package http is the stats bounded context's inbound interface
// adapter. HTTP endpoints translate request/response shapes and call
// application use cases — they own no business logic.
package http

import (
	"encoding/json"
	"net/http"

	"github.com/atilladeniz/next-go-pg/backend/internal/platform/middleware"
	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
	statsapp "github.com/atilladeniz/next-go-pg/backend/internal/stats/application"
	stats "github.com/atilladeniz/next-go-pg/backend/internal/stats/domain"
)

// Handler exposes the stats context's HTTP endpoints. It depends on
// use cases; never on the persistence adapter.
type Handler struct {
	getStats      *statsapp.GetUserStats
	incrementStat *statsapp.IncrementStatField
}

func NewHandler(getStats *statsapp.GetUserStats, incrementStat *statsapp.IncrementStatField) *Handler {
	return &Handler{getStats: getStats, incrementStat: incrementStat}
}

// UserStatsResponse represents user statistics on the wire.
type UserStatsResponse struct {
	UserID        string `json:"userId"`
	ProjectCount  int    `json:"projectCount"`
	ActivityToday int    `json:"activityToday"`
	Notifications int    `json:"notifications"`
	LastLogin     string `json:"lastLogin"`
	MemberSince   string `json:"memberSince"`
}

// ErrorResponse is the stats context's error envelope.
type ErrorResponse struct {
	Error string `json:"error" example:"unauthorized"`
}

// UpdateStatRequest carries a counter mutation.
type UpdateStatRequest struct {
	Field string `json:"field" example:"projects"`
	Delta int    `json:"delta" example:"1"`
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
func (h *Handler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Degraded mode (no DB): return seed defaults so the dashboard renders.
	if h.getStats == nil {
		writeJSON(w, UserStatsResponse{
			UserID:        user.ID,
			ProjectCount:  3,
			ActivityToday: 10,
			Notifications: 2,
			LastLogin:     "2025-12-03T22:00:00Z",
			MemberSince:   "2025-11-15T10:00:00Z",
		})
		return
	}

	userID, err := shared.NewUserID(user.ID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid user id")
		return
	}

	agg, err := h.getStats.Execute(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get stats")
		return
	}

	writeJSON(w, toResponse(agg))
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
func (h *Handler) UpdateUserStats(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req UpdateStatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if h.incrementStat == nil {
		writeError(w, http.StatusServiceUnavailable, "database not available")
		return
	}

	userID, err := shared.NewUserID(user.ID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid user id")
		return
	}

	field, err := stats.NewStatField(req.Field)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	agg, err := h.incrementStat.Execute(r.Context(), userID, field, req.Delta)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update stats")
		return
	}

	writeJSON(w, toResponse(agg))
}

func toResponse(s *stats.UserStats) UserStatsResponse {
	return UserStatsResponse{
		UserID:        string(s.UserID),
		ProjectCount:  s.ProjectCount,
		ActivityToday: s.ActivityToday,
		Notifications: s.Notifications,
		LastLogin:     s.LastLogin.Format("2006-01-02T15:04:05Z"),
		MemberSince:   s.MemberSince.Format("2006-01-02T15:04:05Z"),
	}
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	writeJSON(w, ErrorResponse{Error: message})
}
