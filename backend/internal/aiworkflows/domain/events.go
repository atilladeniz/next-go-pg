package domain

import (
	"time"

	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
)

// SummaryStarted is recorded when a pending run transitions to running.
type SummaryStarted struct {
	SummaryID uint
	UserID    shared.UserID
	RepoURL   RepoURL
	StartedAt time.Time
}

func (SummaryStarted) EventName() string { return "aiworkflows.summary_started" }

// FileSummarized is recorded each time the workflow appends a per-file
// summary to the aggregate. Index is 1-based (1 of N).
type FileSummarized struct {
	SummaryID uint
	UserID    shared.UserID
	Filename  string
	FileIndex int
	FileCount int
}

func (FileSummarized) EventName() string { return "aiworkflows.file_summarized" }

// SummaryCompleted is recorded when the workflow stores the final
// repo-level summary and reaches the terminal `completed` status.
type SummaryCompleted struct {
	SummaryID   uint
	UserID      shared.UserID
	CompletedAt time.Time
}

func (SummaryCompleted) EventName() string { return "aiworkflows.summary_completed" }

// SummaryFailed is recorded when the workflow exhausts retries on any
// step and the run reaches the terminal `failed` status.
type SummaryFailed struct {
	SummaryID uint
	UserID    shared.UserID
	Reason    string
}

func (SummaryFailed) EventName() string { return "aiworkflows.summary_failed" }

// SummaryCancelled is recorded when the run is cancelled before it has
// terminated naturally (e.g. user disconnect).
type SummaryCancelled struct {
	SummaryID uint
	UserID    shared.UserID
}

func (SummaryCancelled) EventName() string { return "aiworkflows.summary_cancelled" }
