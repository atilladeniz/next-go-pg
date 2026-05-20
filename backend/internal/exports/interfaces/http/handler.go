// Package http is the exports bounded context's inbound interface.
package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	exportsapp "github.com/atilladeniz/next-go-pg/backend/internal/exports/application"
	"github.com/atilladeniz/next-go-pg/backend/internal/platform/middleware"
	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
)

// Handler exposes the data-export endpoints.
type Handler struct {
	jobEnqueuer exportsapp.JobEnqueuer
	store       exportsapp.Store
}

func NewHandler(enqueuer exportsapp.JobEnqueuer, store exportsapp.Store) *Handler {
	return &Handler{jobEnqueuer: enqueuer, store: store}
}

// StartExportRequest is the create-export payload.
type StartExportRequest struct {
	Format   string `json:"format"`
	DataType string `json:"dataType"`
}

// StartExportResponse returns the queued job's id.
type StartExportResponse struct {
	JobID   string `json:"jobId"`
	Message string `json:"message"`
}

// ErrorResponse is this context's error envelope.
type ErrorResponse struct {
	Error string `json:"error" example:"unauthorized"`
}

// StartExport godoc
// @Summary Start a data export job
// @Description Starts a background job to export user data
// @Tags export
// @Accept json
// @Produce json
// @Param request body StartExportRequest true "Export parameters"
// @Success 200 {object} StartExportResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /export/start [post]
func (h *Handler) StartExport(w http.ResponseWriter, r *http.Request) {
	var userID string
	if user := middleware.GetUserFromContext(r.Context()); user != nil {
		userID = user.ID
	} else {
		userID = "anonymous"
	}

	var req StartExportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	switch req.Format {
	case "csv", "json":
		// valid
	default:
		respondError(w, http.StatusBadRequest, "invalid format: must be 'csv' or 'json'")
		return
	}

	switch req.DataType {
	case "stats", "activity", "all":
		// valid
	default:
		respondError(w, http.StatusBadRequest, "invalid dataType: must be 'stats', 'activity', or 'all'")
		return
	}

	jobID := uuid.New().String()

	if h.jobEnqueuer == nil {
		respondError(w, http.StatusServiceUnavailable, "export service is currently unavailable")
		return
	}

	if err := h.jobEnqueuer.EnqueueDataExport(context.Background(), jobID, userID, req.Format, req.DataType); err != nil {
		logger.Error().Err(err).Str("job_id", jobID).Msg("Failed to enqueue export job")
		respondError(w, http.StatusInternalServerError, "failed to start export")
		return
	}

	logger.Info().
		Str("job_id", jobID).
		Str("user_id", userID).
		Str("format", req.Format).
		Str("data_type", req.DataType).
		Msg("Export job enqueued")

	respondJSON(w, StartExportResponse{JobID: jobID, Message: "Export gestartet"})
}

// DownloadExport godoc
// @Summary Download a completed export
// @Description Downloads the exported file by download ID
// @Tags export
// @Produce application/octet-stream
// @Param id path string true "Download ID"
// @Success 200 {file} binary
// @Failure 404 {object} ErrorResponse
// @Router /export/download/{id} [get]
func (h *Handler) DownloadExport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	downloadID := vars["id"]
	if downloadID == "" {
		respondError(w, http.StatusBadRequest, "missing download ID")
		return
	}

	result, ok := h.store.Get(downloadID)
	if !ok {
		respondError(w, http.StatusNotFound, "export not found or expired")
		return
	}

	w.Header().Set("Content-Type", result.ContentType)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+result.FileName+"\"")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(result.Data)))
	_, _ = w.Write(result.Data)

	logger.Info().Str("download_id", downloadID).Str("file_name", result.FileName).Msg("Export downloaded")
}

func respondJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	respondJSON(w, ErrorResponse{Error: message})
}
