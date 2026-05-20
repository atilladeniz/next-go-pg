// Package http is the aiworkflows context's inbound interface adapter.
// HTTP endpoints translate request/response shapes and call application
// use cases — they own no business logic and never touch Hatchet, GORM,
// or Ollama directly.
package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	aiapp "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/application"
	ai "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/domain"
	"github.com/atilladeniz/next-go-pg/backend/internal/platform/middleware"
	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
)

// Handler exposes the aiworkflows context's HTTP endpoints.
type Handler struct {
	summarizeRepo  *aiapp.SummarizeRepo
	getRepoSummary *aiapp.GetRepoSummary
}

// NewHandler returns a Handler. Either use case may be nil; in that case
// the corresponding endpoints respond with 503 Service Unavailable so
// the dev stack still boots even without Hatchet wired.
func NewHandler(summarize *aiapp.SummarizeRepo, get *aiapp.GetRepoSummary) *Handler {
	return &Handler{summarizeRepo: summarize, getRepoSummary: get}
}

// SummarizeRepoRequest is the wire-level request body.
type SummarizeRepoRequest struct {
	RepoURL string `json:"repoUrl" example:"https://github.com/owner/repo"`
}

// SummarizeRepoResponse is the 202 body returned to the caller.
type SummarizeRepoResponse struct {
	SummaryID uint   `json:"summaryId" example:"42"`
	RunID     string `json:"runId" example:"a1b2c3d4-..."`
	Status    string `json:"status" example:"pending"`
}

// FileSummaryDTO mirrors the persisted per-file summary.
type FileSummaryDTO struct {
	Filename string `json:"filename"`
	Summary  string `json:"summary"`
}

// RepoSummaryResponse is the 200 body for GET /ai/summaries/{id}.
type RepoSummaryResponse struct {
	ID          uint             `json:"id"`
	RepoURL     string           `json:"repoUrl"`
	Status      string           `json:"status"`
	Files       []FileSummaryDTO `json:"files"`
	Summary     string           `json:"summary"`
	FailReason  string           `json:"failReason,omitempty"`
	StartedAt   string           `json:"startedAt,omitempty"`
	CompletedAt string           `json:"completedAt,omitempty"`
}

// ErrorResponse is the aiworkflows error envelope.
type ErrorResponse struct {
	Error string `json:"error" example:"invalid repo url"`
}

// SummarizeRepo godoc
// @Summary  Trigger a repository summarization workflow
// @Description Enqueues a Hatchet workflow that clones the repository, summarises individual files via Ollama, and produces a repo-level summary.
// @Tags     ai
// @Accept   json
// @Produce  json
// @Param    request body SummarizeRepoRequest true "Repo URL to summarize"
// @Success  202 {object} SummarizeRepoResponse
// @Failure  400 {object} ErrorResponse
// @Failure  401 {object} ErrorResponse
// @Failure  503 {object} ErrorResponse
// @Security BearerAuth
// @Router   /ai/summarize-repo [post]
func (h *Handler) SummarizeRepo(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if h.summarizeRepo == nil {
		writeError(w, http.StatusServiceUnavailable, "ai workflows not configured")
		return
	}

	var req SummarizeRepoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	uid, err := shared.NewUserID(user.ID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid user id")
		return
	}

	out, err := h.summarizeRepo.Execute(r.Context(), aiapp.SummarizeRepoInput{
		UserID:  uid,
		RepoURL: req.RepoURL,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
	writeJSON(w, SummarizeRepoResponse{
		SummaryID: out.SummaryID,
		RunID:     out.RunID,
		Status:    string(ai.StatusPending),
	})
}

// GetRepoSummary godoc
// @Summary  Get a repository summarization result
// @Description Returns the stored summary for a run owned by the authenticated user. Cross-user reads return 404 to avoid leaking existence.
// @Tags     ai
// @Produce  json
// @Param    id path integer true "Summary ID"
// @Success  200 {object} RepoSummaryResponse
// @Failure  401 {object} ErrorResponse
// @Failure  404 {object} ErrorResponse
// @Failure  503 {object} ErrorResponse
// @Security BearerAuth
// @Router   /ai/summaries/{id} [get]
func (h *Handler) GetRepoSummary(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if h.getRepoSummary == nil {
		writeError(w, http.StatusServiceUnavailable, "ai workflows not configured")
		return
	}

	vars := mux.Vars(r)
	idRaw := vars["id"]
	id64, err := strconv.ParseUint(idRaw, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	uid, err := shared.NewUserID(user.ID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid user id")
		return
	}

	agg, err := h.getRepoSummary.Execute(r.Context(), aiapp.GetRepoSummaryInput{
		UserID:    uid,
		SummaryID: uint(id64),
	})
	if err != nil {
		if errors.Is(err, aiapp.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to load summary")
		return
	}

	writeJSON(w, toResponse(agg))
}

func toResponse(s *ai.RepoSummary) RepoSummaryResponse {
	files := make([]FileSummaryDTO, 0, len(s.Files))
	for _, f := range s.Files {
		files = append(files, FileSummaryDTO{Filename: f.Filename(), Summary: f.Summary()})
	}
	resp := RepoSummaryResponse{
		ID:         s.ID,
		RepoURL:    s.RepoURL.String(),
		Status:     s.Status.String(),
		Files:      files,
		Summary:    s.Summary,
		FailReason: s.FailReason,
	}
	if !s.StartedAt.IsZero() {
		resp.StartedAt = s.StartedAt.UTC().Format("2006-01-02T15:04:05Z")
	}
	if !s.CompletedAt.IsZero() {
		resp.CompletedAt = s.CompletedAt.UTC().Format("2006-01-02T15:04:05Z")
	}
	return resp
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	writeJSON(w, ErrorResponse{Error: message})
}
