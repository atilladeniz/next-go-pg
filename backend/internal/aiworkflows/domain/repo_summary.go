package domain

import (
	"fmt"
	"time"

	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
)

// RepoSummary is the aggregate root for a single repository-summarization
// run. It owns the lifecycle transitions and records domain events on
// mutation. The aggregate carries no persistence concerns — the GORM
// model and mapper live in infrastructure/persistence/.
type RepoSummary struct {
	shared.AggregateBase

	ID          uint
	UserID      shared.UserID
	RepoURL     RepoURL
	Status      Status
	Files       []FileSummary
	Summary     string
	FailReason  string
	StartedAt   time.Time
	CompletedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var _ shared.AggregateRoot = (*RepoSummary)(nil)

// NewRepoSummary is the factory for a freshly-requested run. The
// aggregate starts in StatusPending and records no event — the
// SummaryStarted event fires when the workflow's first step executes
// MarkStarted, not at enqueue time.
func NewRepoSummary(userID shared.UserID, repoURL RepoURL) *RepoSummary {
	return &RepoSummary{
		UserID:  userID,
		RepoURL: repoURL,
		Status:  StatusPending,
	}
}

// MarkStarted transitions pending → running and records SummaryStarted.
func (r *RepoSummary) MarkStarted(at time.Time) error {
	if r.Status != StatusPending {
		return fmt.Errorf("cannot start: status is %s, want pending", r.Status)
	}
	r.Status = StatusRunning
	r.StartedAt = at
	r.Record(SummaryStarted{
		SummaryID: r.ID,
		UserID:    r.UserID,
		RepoURL:   r.RepoURL,
		StartedAt: at,
	})
	return nil
}

// AppendFileSummary records one per-file summary and emits a
// FileSummarized event. totalFiles is the count discovered by Traverse,
// known to the workflow but not to the aggregate; passing it through
// keeps the event self-contained for SSE consumers downstream.
func (r *RepoSummary) AppendFileSummary(fs FileSummary, totalFiles int) error {
	if r.Status != StatusRunning {
		return fmt.Errorf("cannot append file summary: status is %s, want running", r.Status)
	}
	r.Files = append(r.Files, fs)
	r.Record(FileSummarized{
		SummaryID: r.ID,
		UserID:    r.UserID,
		Filename:  fs.Filename(),
		FileIndex: len(r.Files),
		FileCount: totalFiles,
	})
	return nil
}

// MarkCompleted transitions running → completed, stores the repo-level
// summary text, and records SummaryCompleted.
func (r *RepoSummary) MarkCompleted(summary string, at time.Time) error {
	if r.Status != StatusRunning {
		return fmt.Errorf("cannot complete: status is %s, want running", r.Status)
	}
	r.Status = StatusCompleted
	r.Summary = summary
	r.CompletedAt = at
	r.Record(SummaryCompleted{
		SummaryID:   r.ID,
		UserID:      r.UserID,
		CompletedAt: at,
	})
	return nil
}

// MarkFailed transitions pending/running → failed and records the reason.
func (r *RepoSummary) MarkFailed(reason string, at time.Time) error {
	if r.Status.IsTerminal() {
		return fmt.Errorf("cannot fail run already in terminal status %s", r.Status)
	}
	r.Status = StatusFailed
	r.FailReason = reason
	r.CompletedAt = at
	r.Record(SummaryFailed{
		SummaryID: r.ID,
		UserID:    r.UserID,
		Reason:    reason,
	})
	return nil
}

// MarkCancelled transitions pending/running → cancelled. Terminal runs
// (already completed/failed/cancelled) cannot be cancelled.
func (r *RepoSummary) MarkCancelled(at time.Time) error {
	if r.Status.IsTerminal() {
		return fmt.Errorf("cannot cancel run already in terminal status %s", r.Status)
	}
	r.Status = StatusCancelled
	r.CompletedAt = at
	r.Record(SummaryCancelled{
		SummaryID: r.ID,
		UserID:    r.UserID,
	})
	return nil
}
