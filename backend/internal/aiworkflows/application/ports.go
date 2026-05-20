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

// LLMClient is the LLM-runtime abstraction. Implementations talk to
// Ollama in dev; in production they could fan out to a hosted provider.
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

// ProgressPublisher dispatches the domain events the workflow records on
// the aggregate. The SSE adapter under infrastructure/events maps each
// event onto the `ai-progress` SSE event type.
type ProgressPublisher interface {
	Publish(ctx context.Context, events ...shared.DomainEvent) error
}
