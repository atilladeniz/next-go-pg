package application_test

import (
	"context"
	"errors"
	"testing"

	aiapp "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/application"
	ai "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/domain"
	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
)

// fakeStore is an in-memory Store implementation for unit tests.
type fakeStore struct {
	rows        map[uint]*ai.RepoSummary
	nextID      uint
	createErr   error
	getErr      error
	saveErr     error
	createCalls int
	saveCalls   int
}

func newFakeStore() *fakeStore {
	return &fakeStore{rows: map[uint]*ai.RepoSummary{}, nextID: 1}
}

func (s *fakeStore) Create(_ context.Context, agg *ai.RepoSummary) error {
	s.createCalls++
	if s.createErr != nil {
		return s.createErr
	}
	agg.ID = s.nextID
	s.nextID++
	s.rows[agg.ID] = agg
	return nil
}

func (s *fakeStore) Save(_ context.Context, agg *ai.RepoSummary) error {
	s.saveCalls++
	if s.saveErr != nil {
		return s.saveErr
	}
	s.rows[agg.ID] = agg
	return nil
}

func (s *fakeStore) GetByID(_ context.Context, id uint) (*ai.RepoSummary, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	row, ok := s.rows[id]
	if !ok {
		return nil, aiapp.ErrNotFound
	}
	return row, nil
}

func (s *fakeStore) ListByUserID(_ context.Context, userID shared.UserID, limit int) ([]*ai.RepoSummary, error) {
	out := make([]*ai.RepoSummary, 0)
	for _, row := range s.rows {
		if row.UserID == userID {
			out = append(out, row)
		}
	}
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

// fakeEnqueuer is a stub HatchetEnqueuer.
type fakeEnqueuer struct {
	runID string
	err   error
	calls int
	last  aiapp.EnqueueSummarizeRepoInput
}

func (e *fakeEnqueuer) EnqueueSummarizeRepo(_ context.Context, in aiapp.EnqueueSummarizeRepoInput) (string, error) {
	e.calls++
	e.last = in
	return e.runID, e.err
}

func uid(t *testing.T, s string) shared.UserID {
	t.Helper()
	u, err := shared.NewUserID(s)
	if err != nil {
		t.Fatalf("NewUserID: %v", err)
	}
	return u
}

func TestSummarizeRepo_HappyPath(t *testing.T) {
	t.Parallel()
	store := newFakeStore()
	enq := &fakeEnqueuer{runID: "run-123"}

	uc := aiapp.SummarizeRepo{Store: store, Enqueuer: enq}
	out, err := uc.Execute(context.Background(), aiapp.SummarizeRepoInput{
		UserID:  uid(t, "user-1"),
		RepoURL: "https://github.com/owner/repo",
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if out.RunID != "run-123" {
		t.Errorf("RunID = %q, want run-123", out.RunID)
	}
	if out.SummaryID == 0 {
		t.Errorf("SummaryID should be non-zero")
	}
	if store.createCalls != 1 {
		t.Errorf("Create calls = %d, want 1", store.createCalls)
	}
	if enq.calls != 1 {
		t.Errorf("Enqueue calls = %d, want 1", enq.calls)
	}
	if enq.last.SummaryID != out.SummaryID {
		t.Errorf("Enqueue SummaryID mismatch: got %d, want %d", enq.last.SummaryID, out.SummaryID)
	}
	if enq.last.RepoURL != "https://github.com/owner/repo" {
		t.Errorf("Enqueue RepoURL = %q", enq.last.RepoURL)
	}
}

func TestSummarizeRepo_InvalidURL(t *testing.T) {
	t.Parallel()
	cases := []string{
		"",
		"not-a-url",
		"ssh://git@github.com/owner/repo",
		"https://github.com/owner/repo;rm -rf /",
	}
	for _, raw := range cases {
		t.Run(raw, func(t *testing.T) {
			store := newFakeStore()
			enq := &fakeEnqueuer{}
			uc := aiapp.SummarizeRepo{Store: store, Enqueuer: enq}
			_, err := uc.Execute(context.Background(), aiapp.SummarizeRepoInput{
				UserID:  uid(t, "user-1"),
				RepoURL: raw,
			})
			if err == nil {
				t.Errorf("expected error for %q, got nil", raw)
			}
			if store.createCalls != 0 {
				t.Errorf("Create should not be called on invalid URL, got %d", store.createCalls)
			}
			if enq.calls != 0 {
				t.Errorf("Enqueue should not be called on invalid URL, got %d", enq.calls)
			}
		})
	}
}

func TestSummarizeRepo_StoreCreateError(t *testing.T) {
	t.Parallel()
	store := newFakeStore()
	store.createErr = errors.New("db down")
	enq := &fakeEnqueuer{runID: "ignored"}
	uc := aiapp.SummarizeRepo{Store: store, Enqueuer: enq}

	_, err := uc.Execute(context.Background(), aiapp.SummarizeRepoInput{
		UserID:  uid(t, "user-1"),
		RepoURL: "https://github.com/owner/repo",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if enq.calls != 0 {
		t.Errorf("Enqueue should not be called if Create failed")
	}
}

func TestSummarizeRepo_EnqueueErrorMarksFailed(t *testing.T) {
	t.Parallel()
	store := newFakeStore()
	enq := &fakeEnqueuer{err: errors.New("hatchet unavailable")}
	uc := aiapp.SummarizeRepo{Store: store, Enqueuer: enq}

	_, err := uc.Execute(context.Background(), aiapp.SummarizeRepoInput{
		UserID:  uid(t, "user-1"),
		RepoURL: "https://github.com/owner/repo",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	// One row was created and then Save was called to mark it failed.
	if len(store.rows) != 1 {
		t.Errorf("expected 1 row in store, got %d", len(store.rows))
	}
	if store.saveCalls != 1 {
		t.Errorf("Save calls = %d, want 1 (mark-failed path)", store.saveCalls)
	}
	for _, row := range store.rows {
		if row.Status != ai.StatusFailed {
			t.Errorf("row status = %s, want failed", row.Status)
		}
	}
}

func TestGetRepoSummary_HappyPath(t *testing.T) {
	t.Parallel()
	store := newFakeStore()
	user := uid(t, "user-1")
	url, _ := ai.NewRepoURL("https://github.com/owner/repo")
	agg := ai.NewRepoSummary(user, url)
	_ = store.Create(context.Background(), agg)

	uc := aiapp.GetRepoSummary{Store: store}
	got, err := uc.Execute(context.Background(), aiapp.GetRepoSummaryInput{
		UserID:    user,
		SummaryID: agg.ID,
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if got.ID != agg.ID {
		t.Errorf("ID mismatch: got %d, want %d", got.ID, agg.ID)
	}
}

func TestGetRepoSummary_OwnershipMismatchReturns404(t *testing.T) {
	t.Parallel()
	store := newFakeStore()
	owner := uid(t, "user-1")
	url, _ := ai.NewRepoURL("https://github.com/owner/repo")
	agg := ai.NewRepoSummary(owner, url)
	_ = store.Create(context.Background(), agg)

	uc := aiapp.GetRepoSummary{Store: store}
	_, err := uc.Execute(context.Background(), aiapp.GetRepoSummaryInput{
		UserID:    uid(t, "different-user"),
		SummaryID: agg.ID,
	})
	if !errors.Is(err, aiapp.ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestGetRepoSummary_MissingReturns404(t *testing.T) {
	t.Parallel()
	store := newFakeStore()
	uc := aiapp.GetRepoSummary{Store: store}

	_, err := uc.Execute(context.Background(), aiapp.GetRepoSummaryInput{
		UserID:    uid(t, "user-1"),
		SummaryID: 9999,
	})
	if !errors.Is(err, aiapp.ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}
