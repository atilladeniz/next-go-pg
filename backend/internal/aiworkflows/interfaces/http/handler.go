// Package http is the aiworkflows context's inbound interface adapter.
// HTTP endpoints translate request/response shapes and call application
// use cases — they own no business logic and never touch Hatchet, GORM,
// or the LLM client directly.
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
	listSummaries  *aiapp.ListUserSummaries
	deleteSummary  *aiapp.DeleteUserSummary
}

// NewHandler returns a Handler. Any use case may be nil; in that case
// the corresponding endpoints respond with 503 Service Unavailable so
// the dev stack still boots even without Hatchet wired.
func NewHandler(summarize *aiapp.SummarizeRepo, get *aiapp.GetRepoSummary, list *aiapp.ListUserSummaries, del *aiapp.DeleteUserSummary) *Handler {
	return &Handler{summarizeRepo: summarize, getRepoSummary: get, listSummaries: list, deleteSummary: del}
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
	ID            uint             `json:"id"`
	RepoURL       string           `json:"repoUrl"`
	Status        string           `json:"status"`
	Files         []FileSummaryDTO `json:"files"`
	Summary       string           `json:"summary"`
	FailReason    string           `json:"failReason,omitempty"`
	StartedAt     string           `json:"startedAt,omitempty"`
	CompletedAt   string           `json:"completedAt,omitempty"`
	StepDurations map[string]int64 `json:"stepDurations,omitempty"`
}

// RepoSummaryListItem is the compact projection returned by GET /ai/summaries.
type RepoSummaryListItem struct {
	ID        uint   `json:"id"`
	RepoURL   string `json:"repoUrl"`
	Status    string `json:"status"`
	FileCount int    `json:"fileCount"`
	StartedAt string `json:"startedAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// RepoSummaryListResponse wraps the list so we can add pagination later
// without a breaking JSON change.
type RepoSummaryListResponse struct {
	Items []RepoSummaryListItem `json:"items"`
}

// ErrorResponse is the aiworkflows error envelope.
type ErrorResponse struct {
	Error string `json:"error" example:"invalid repo url"`
}

// SummarizeRepo godoc
// @Summary  Trigger a repository summarization workflow
// @Description Enqueues a Hatchet workflow that clones the repository, summarises individual files via the configured LLM provider (OpenRouter), and produces a repo-level summary.
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

	writeJSONStatus(w, http.StatusAccepted, SummarizeRepoResponse{
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

// ListRepoSummaries godoc
// @Summary  List the user's recent repository summaries
// @Description Returns up to 50 of the authenticated user's runs, newest first. Used by the AI page to show a history of past runs.
// @Tags     ai
// @Produce  json
// @Success  200 {object} RepoSummaryListResponse
// @Failure  401 {object} ErrorResponse
// @Failure  503 {object} ErrorResponse
// @Security BearerAuth
// @Router   /ai/summaries [get]
func (h *Handler) ListRepoSummaries(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if h.listSummaries == nil {
		writeError(w, http.StatusServiceUnavailable, "ai workflows not configured")
		return
	}
	uid, err := shared.NewUserID(user.ID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid user id")
		return
	}
	rows, err := h.listSummaries.Execute(r.Context(), uid, 20)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load summaries")
		return
	}
	items := make([]RepoSummaryListItem, 0, len(rows))
	for _, row := range rows {
		item := RepoSummaryListItem{
			ID:        row.ID,
			RepoURL:   row.RepoURL.String(),
			Status:    row.Status.String(),
			FileCount: len(row.Files),
		}
		if !row.StartedAt.IsZero() {
			item.StartedAt = row.StartedAt.UTC().Format("2006-01-02T15:04:05Z")
		}
		if !row.UpdatedAt.IsZero() {
			item.UpdatedAt = row.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z")
		}
		items = append(items, item)
	}
	writeJSON(w, RepoSummaryListResponse{Items: items})
}

// DeleteRepoSummary godoc
// @Summary  Delete a repository summary run
// @Description Removes a run owned by the authenticated user. Returns 404 for missing rows AND cross-user deletes to avoid leaking existence.
// @Tags     ai
// @Produce  json
// @Param    id path integer true "Summary ID"
// @Success  204 "No Content"
// @Failure  401 {object} ErrorResponse
// @Failure  404 {object} ErrorResponse
// @Failure  503 {object} ErrorResponse
// @Security BearerAuth
// @Router   /ai/summaries/{id} [delete]
func (h *Handler) DeleteRepoSummary(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if h.deleteSummary == nil {
		writeError(w, http.StatusServiceUnavailable, "ai workflows not configured")
		return
	}

	vars := mux.Vars(r)
	id64, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	uid, err := shared.NewUserID(user.ID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid user id")
		return
	}

	err = h.deleteSummary.Execute(r.Context(), aiapp.DeleteUserSummaryInput{
		UserID:    uid,
		SummaryID: uint(id64),
	})
	if err != nil {
		if errors.Is(err, aiapp.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete summary")
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
	writeJSONStatus(w, http.StatusOK, payload)
}

// writeJSONStatus sets Content-Type BEFORE WriteHeader. Reversing the
// order silently strips the Content-Type header (Go's http.ResponseWriter
// freezes headers at WriteHeader time), which then makes the frontend
// `customFetch` fall through to the text branch and return `{message: …}`
// instead of the parsed JSON body. Symptom: 202 succeeds but the
// returned `data` object has no `summaryId` field.
func writeJSONStatus(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSONStatus(w, status, ErrorResponse{Error: message})
}
