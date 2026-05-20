package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	aiapp "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/application"
	ai "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/domain"
	aihttp "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/interfaces/http"
	"github.com/atilladeniz/next-go-pg/backend/internal/platform/middleware"
	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
)

// withUser injects an authenticated user into the request context so the
// handler's middleware lookup finds someone.
func withUser(req *stdhttp.Request, userID string) *stdhttp.Request {
	ctx := context.WithValue(req.Context(), middleware.UserContextKey, &middleware.User{
		ID:    userID,
		Email: userID + "@example.com",
	})
	return req.WithContext(ctx)
}

// fakeStore is a minimal in-memory aiapp.Store for use-case wiring.
type fakeStore struct {
	rows   map[uint]*ai.RepoSummary
	nextID uint
}

func newFakeStore() *fakeStore {
	return &fakeStore{rows: map[uint]*ai.RepoSummary{}, nextID: 1}
}

func (s *fakeStore) Create(_ context.Context, agg *ai.RepoSummary) error {
	agg.ID = s.nextID
	s.nextID++
	s.rows[agg.ID] = agg
	return nil
}

func (s *fakeStore) Save(_ context.Context, agg *ai.RepoSummary) error {
	s.rows[agg.ID] = agg
	return nil
}

func (s *fakeStore) GetByID(_ context.Context, id uint) (*ai.RepoSummary, error) {
	row, ok := s.rows[id]
	if !ok {
		return nil, aiapp.ErrNotFound
	}
	return row, nil
}

type fakeEnqueuer struct {
	runID string
	err   error
}

func (e *fakeEnqueuer) EnqueueSummarizeRepo(_ context.Context, _ aiapp.EnqueueSummarizeRepoInput) (string, error) {
	return e.runID, e.err
}

func TestSummarizeRepo_HappyPath(t *testing.T) {
	t.Parallel()
	store := newFakeStore()
	enq := &fakeEnqueuer{runID: "run-xyz"}
	h := aihttp.NewHandler(&aiapp.SummarizeRepo{Store: store, Enqueuer: enq}, &aiapp.GetRepoSummary{Store: store})

	body := bytes.NewBufferString(`{"repoUrl":"https://github.com/owner/repo"}`)
	req := withUser(httptest.NewRequest(stdhttp.MethodPost, "/api/v1/ai/summarize-repo", body), "user-1")
	w := httptest.NewRecorder()

	h.SummarizeRepo(w, req)

	if w.Code != stdhttp.StatusAccepted {
		t.Fatalf("status = %d, want 202; body=%s", w.Code, w.Body.String())
	}
	var resp aihttp.SummarizeRepoResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.RunID != "run-xyz" {
		t.Errorf("RunID = %q", resp.RunID)
	}
	if resp.SummaryID == 0 {
		t.Errorf("SummaryID should be non-zero")
	}
}

func TestSummarizeRepo_Unauthenticated(t *testing.T) {
	t.Parallel()
	h := aihttp.NewHandler(&aiapp.SummarizeRepo{Store: newFakeStore(), Enqueuer: &fakeEnqueuer{}}, nil)

	req := httptest.NewRequest(stdhttp.MethodPost, "/api/v1/ai/summarize-repo",
		strings.NewReader(`{"repoUrl":"https://github.com/owner/repo"}`))
	w := httptest.NewRecorder()
	h.SummarizeRepo(w, req)

	if w.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
}

func TestSummarizeRepo_DegradedMode(t *testing.T) {
	t.Parallel()
	// summarizeUC nil — token missing path.
	h := aihttp.NewHandler(nil, &aiapp.GetRepoSummary{Store: newFakeStore()})

	body := strings.NewReader(`{"repoUrl":"https://github.com/owner/repo"}`)
	req := withUser(httptest.NewRequest(stdhttp.MethodPost, "/api/v1/ai/summarize-repo", body), "user-1")
	w := httptest.NewRecorder()
	h.SummarizeRepo(w, req)

	if w.Code != stdhttp.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", w.Code)
	}
}

func TestSummarizeRepo_InvalidURL(t *testing.T) {
	t.Parallel()
	store := newFakeStore()
	enq := &fakeEnqueuer{}
	h := aihttp.NewHandler(&aiapp.SummarizeRepo{Store: store, Enqueuer: enq}, nil)

	body := strings.NewReader(`{"repoUrl":"not-a-url"}`)
	req := withUser(httptest.NewRequest(stdhttp.MethodPost, "/api/v1/ai/summarize-repo", body), "user-1")
	w := httptest.NewRecorder()
	h.SummarizeRepo(w, req)

	if w.Code != stdhttp.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestGetRepoSummary_OwnershipReturns404(t *testing.T) {
	t.Parallel()
	store := newFakeStore()
	owner, _ := shared.NewUserID("user-owner")
	url, _ := ai.NewRepoURL("https://github.com/owner/repo")
	agg := ai.NewRepoSummary(owner, url)
	_ = store.Create(context.Background(), agg)

	h := aihttp.NewHandler(nil, &aiapp.GetRepoSummary{Store: store})

	// Wrap router so {id} resolves.
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/ai/summaries/{id}", h.GetRepoSummary).Methods("GET")

	req := withUser(httptest.NewRequest(stdhttp.MethodGet, "/api/v1/ai/summaries/1", nil), "other-user")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != stdhttp.StatusNotFound {
		t.Fatalf("status = %d, want 404 (cross-user must NOT leak existence)", w.Code)
	}
}

func TestGetRepoSummary_HappyPath(t *testing.T) {
	t.Parallel()
	store := newFakeStore()
	owner, _ := shared.NewUserID("user-1")
	url, _ := ai.NewRepoURL("https://github.com/owner/repo")
	agg := ai.NewRepoSummary(owner, url)
	_ = store.Create(context.Background(), agg)

	h := aihttp.NewHandler(nil, &aiapp.GetRepoSummary{Store: store})
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/ai/summaries/{id}", h.GetRepoSummary).Methods("GET")

	req := withUser(httptest.NewRequest(stdhttp.MethodGet, "/api/v1/ai/summaries/1", nil), "user-1")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != stdhttp.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
	var resp aihttp.RepoSummaryResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.RepoURL != "https://github.com/owner/repo" {
		t.Errorf("RepoURL = %q", resp.RepoURL)
	}
}
