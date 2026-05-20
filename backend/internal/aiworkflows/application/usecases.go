package application

import (
	"context"
	"fmt"

	ai "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/domain"
	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
)

// SummarizeRepo validates the request, persists a freshly-pending
// RepoSummary aggregate, and enqueues a workflow run. The aggregate is
// created BEFORE enqueue so the workflow's first step can load it by
// SummaryID. If enqueue fails, the use case marks the aggregate as
// failed so the row does not linger in `pending` forever.
type SummarizeRepo struct {
	Store    Store
	Enqueuer HatchetEnqueuer
}

// SummarizeRepoInput is the wire-level request. The use case is
// responsible for validating RepoURL — handlers MUST NOT pre-validate.
type SummarizeRepoInput struct {
	UserID  shared.UserID
	RepoURL string
}

// SummarizeRepoOutput is returned to the HTTP layer; the RunID is the
// Hatchet workflow run ID and the SummaryID is the aggregate ID used
// for subsequent reads.
type SummarizeRepoOutput struct {
	SummaryID uint
	RunID     string
}

func (uc SummarizeRepo) Execute(ctx context.Context, in SummarizeRepoInput) (SummarizeRepoOutput, error) {
	url, err := ai.NewRepoURL(in.RepoURL)
	if err != nil {
		return SummarizeRepoOutput{}, fmt.Errorf("invalid repo url: %w", err)
	}
	agg := ai.NewRepoSummary(in.UserID, url)
	if err := uc.Store.Create(ctx, agg); err != nil {
		return SummarizeRepoOutput{}, fmt.Errorf("store create: %w", err)
	}

	runID, err := uc.Enqueuer.EnqueueSummarizeRepo(ctx, EnqueueSummarizeRepoInput{
		SummaryID: agg.ID,
		UserID:    agg.UserID,
		RepoURL:   agg.RepoURL,
	})
	if err != nil {
		// Best-effort: mark the row failed so it doesn't sit in `pending`.
		// We deliberately ignore Save errors here — the original error is
		// more useful to the caller.
		if markErr := agg.MarkFailed("workflow enqueue failed: "+err.Error(), nowFn()); markErr == nil {
			_ = uc.Store.Save(ctx, agg)
		}
		return SummarizeRepoOutput{}, fmt.Errorf("enqueue workflow: %w", err)
	}

	return SummarizeRepoOutput{SummaryID: agg.ID, RunID: runID}, nil
}

// GetRepoSummary loads a RepoSummary aggregate for the requesting user.
// Returns ErrNotFound for both missing rows AND cross-user reads so the
// HTTP layer maps cleanly to a 404 without leaking existence.
type GetRepoSummary struct {
	Store Store
}

type GetRepoSummaryInput struct {
	UserID    shared.UserID
	SummaryID uint
}

func (uc GetRepoSummary) Execute(ctx context.Context, in GetRepoSummaryInput) (*ai.RepoSummary, error) {
	agg, err := uc.Store.GetByID(ctx, in.SummaryID)
	if err != nil {
		return nil, err
	}
	if agg.UserID != in.UserID {
		return nil, ErrNotFound
	}
	return agg, nil
}
