// Package application is the aiworkflows bounded context's use-case
// layer. Ports declared here are consumed by use cases and implemented
// by adapters under internal/aiworkflows/infrastructure/. No package in
// this layer may import the Hatchet SDK, GORM, or net/http directly.
package application

import (
	"context"
	"errors"

	ai "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/domain"
	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
)

// ErrNotFound is returned by Store and use cases when a RepoSummary is
// missing OR when the requesting user does not own it. Returning the
// same value in both cases lets the HTTP layer map cleanly to a 404 and
// avoids leaking existence of other users' runs.
var ErrNotFound = errors.New("repo summary not found")

// Store persists and retrieves RepoSummary aggregates. The contract is:
//   - Create assigns a non-zero ID on success.
//   - Save MUST mutate fields in place rather than replacing *agg, so
//     the aggregate's pending domain events survive persistence.
//   - GetByID returns ErrNotFound when the row does not exist.
type Store interface {
	Create(ctx context.Context, agg *ai.RepoSummary) error
	Save(ctx context.Context, agg *ai.RepoSummary) error
	GetByID(ctx context.Context, id uint) (*ai.RepoSummary, error)
	ListByUserID(ctx context.Context, userID shared.UserID, limit int) ([]*ai.RepoSummary, error)
}

// HatchetEnqueuer hides the Hatchet SDK from the application and HTTP
// layers. Swapping the workflow engine should only require a new
// infrastructure-layer adapter.
type HatchetEnqueuer interface {
	EnqueueSummarizeRepo(ctx context.Context, in EnqueueSummarizeRepoInput) (runID string, err error)
}

// EnqueueSummarizeRepoInput is the payload the engine adapter forwards
// into the workflow run. Carries the minimum the steps need; everything
// else is loaded by SummaryID inside the workflow.
type EnqueueSummarizeRepoInput struct {
	SummaryID uint
	UserID    shared.UserID
	RepoURL   ai.RepoURL
}

// LLMClient is the LLM-runtime abstraction. Current implementation
// talks to OpenRouter; the port stays generic so swapping in another
// provider (local model, different gateway) is a one-line wire change.
type LLMClient interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

// RepoCloner produces a local working copy of a public Git repository.
// Callers MUST invoke Cleanup when done with the path, even on error.
type RepoCloner interface {
	Clone(ctx context.Context, url ai.RepoURL) (ClonedRepo, error)
}

// ClonedRepo is the result of a successful RepoCloner.Clone.
type ClonedRepo struct {
	Path    string
	Cleanup func() error
}

// ProgressPublisher dispatches workflow progress to the frontend.
//
// Two channels:
//   - `Publish` routes the bounded context's domain events (started,
//     completed, failed, cancelled). The adapter under
//     infrastructure/events maps each onto an `ai-progress` SSE event.
//   - `PublishStep` is a thin bypass for fine-grained step-state and
//     per-file fan-out progress that the aggregate does NOT own. It
//     fires on every step boundary plus once per file completion so the
//     frontend can render a live "step N of 5 — file 3 of 5" view
//     without waiting for the orchestrator's WaitGroup to drain.
type ProgressPublisher interface {
	Publish(ctx context.Context, events ...shared.DomainEvent) error
	PublishStep(ctx context.Context, step StepProgress)
}

// StepName enumerates the workflow's main steps. Kept as a typed string
// so the frontend can match on values without a wire-level fragility check.
type StepName string

const (
	StepClone          StepName = "clone"
	StepTraverse       StepName = "traverse"
	StepSummarizeFiles StepName = "summarize_files"
	StepAggregate      StepName = "aggregate"
	StepStore          StepName = "store"
)

// StepState is the wire-level state of one step.
type StepState string

const (
	StepStateStarted   StepState = "started"
	StepStateCompleted StepState = "completed"
	StepStateFailed    StepState = "failed"
	StepStateProgress  StepState = "progress" // per-file ticks within summarize_files
)

// StepProgress is the payload published on a step transition. Use the
// zero value for unset numeric fields.
type StepProgress struct {
	SummaryID  uint
	UserID     shared.UserID
	Step       StepName
	State      StepState
	DurationMs int64
	FileIndex  int    // 1-based, for summarize_files state=progress
	FileCount  int    // total files Traverse selected
	Filename   string // last completed filename
	Reason     string // populated only when State=failed
}
