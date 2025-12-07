package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/atilladeniz/next-go-pg/backend/internal/jobs"
	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
)

// ExportHandler handles data export endpoints
type ExportHandler struct {
	jobEnqueuer jobs.JobEnqueuer
	exportStore *jobs.ExportStore
}

// NewExportHandler creates a new export handler
func NewExportHandler(enqueuer jobs.JobEnqueuer, store *jobs.ExportStore) *ExportHandler {
	return &ExportHandler{
		jobEnqueuer: enqueuer,
		exportStore: store,
	}
}

// StartExportRequest represents the request to start an export
type StartExportRequest struct {
	Format   string `json:"format"`   // "csv" or "json"
	DataType string `json:"dataType"` // "stats", "activity", or "all"
}

// StartExportResponse represents the response after starting an export
type StartExportResponse struct {
	JobID   string `json:"jobId"`
	Message string `json:"message"`
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
func (h *ExportHandler) StartExport(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	userID := r.Context().Value("user_id")
	if userID == nil {
		userID = "anonymous"
	}

	var req StartExportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate format
	var format jobs.ExportFormat
	switch req.Format {
	case "csv":
		format = jobs.ExportFormatCSV
	case "json":
		format = jobs.ExportFormatJSON
	default:
		respondError(w, http.StatusBadRequest, "invalid format: must be 'csv' or 'json'")
		return
	}

	// Validate data type
	switch req.DataType {
	case "stats", "activity", "all":
		// valid
	default:
		respondError(w, http.StatusBadRequest, "invalid dataType: must be 'stats', 'activity', or 'all'")
		return
	}

	// Generate job ID
	jobID := uuid.New().String()

	// Check if job queue is available
	if h.jobEnqueuer == nil {
		respondError(w, http.StatusServiceUnavailable, "export service is currently unavailable")
		return
	}

	// Enqueue the export job
	if err := jobs.EnqueueDataExport(context.Background(), h.jobEnqueuer, jobID, userID.(string), format, req.DataType); err != nil {
		logger.Error().Err(err).Str("job_id", jobID).Msg("Failed to enqueue export job")
		respondError(w, http.StatusInternalServerError, "failed to start export")
		return
	}

	logger.Info().
		Str("job_id", jobID).
		Str("user_id", userID.(string)).
		Str("format", req.Format).
		Str("data_type", req.DataType).
		Msg("Export job enqueued")

	respondJSON(w, StartExportResponse{
		JobID:   jobID,
		Message: "Export gestartet",
	})
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
func (h *ExportHandler) DownloadExport(w http.ResponseWriter, r *http.Request) {
	// Get download ID from URL
	downloadID := r.PathValue("id")
	if downloadID == "" {
		// Fallback for older routers
		downloadID = r.URL.Query().Get("id")
	}

	if downloadID == "" {
		respondError(w, http.StatusBadRequest, "missing download ID")
		return
	}

	// Get export from store
	result, ok := h.exportStore.Get(downloadID)
	if !ok {
		respondError(w, http.StatusNotFound, "export not found or expired")
		return
	}

	// Set headers for download
	w.Header().Set("Content-Type", result.ContentType)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+result.FileName+"\"")
	w.Header().Set("Content-Length", string(rune(len(result.Data))))

	// Write data
	w.Write(result.Data)

	// Optionally delete after download
	// h.exportStore.Delete(downloadID)

	logger.Info().
		Str("download_id", downloadID).
		Str("file_name", result.FileName).
		Msg("Export downloaded")
}
