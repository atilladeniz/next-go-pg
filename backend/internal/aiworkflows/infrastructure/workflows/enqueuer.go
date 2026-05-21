package workflows

import (
	"context"

	hatchet "github.com/hatchet-dev/hatchet/sdks/go"

	aiapp "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/application"
)

// Enqueuer is the HatchetEnqueuer adapter. It hides the Hatchet client
// from the application layer; swapping engines means only this file
// (plus workflow.go and steps.go) changes.
type Enqueuer struct {
	Client *hatchet.Client
}

// NewEnqueuer constructs an Enqueuer over a configured Hatchet client.
func NewEnqueuer(client *hatchet.Client) *Enqueuer {
	return &Enqueuer{Client: client}
}

// EnqueueSummarizeRepo kicks off a `summarize-repo` workflow run.
// Returns the Hatchet run ID so the HTTP layer can echo it to the
// client (the frontend uses it as a correlation key for SSE events).
func (e *Enqueuer) EnqueueSummarizeRepo(ctx context.Context, in aiapp.EnqueueSummarizeRepoInput) (string, error) {
	ref, err := e.Client.RunNoWait(ctx, WorkflowName, WorkflowInput{
		SummaryID: in.SummaryID,
		UserID:    in.UserID.String(),
		RepoURL:   in.RepoURL.String(),
	})
	if err != nil {
		return "", err
	}
	return ref.RunId, nil
}

// Static port-conformance check.
var _ aiapp.HatchetEnqueuer = (*Enqueuer)(nil)
