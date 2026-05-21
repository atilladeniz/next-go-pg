package workflows

import (
	"context"

	hatchet "github.com/hatchet-dev/hatchet/sdks/go"
)

// Worker is the AI workflow worker — a long-running goroutine that
// pulls tasks for `summarize-repo` and the `summarize-file` child.
// Started from composition.Build via Worker.Start; stopped by cancelling
// the context.
type Worker struct {
	hatchet *hatchet.Worker
	defs    Definitions
}

// NewWorker registers the workflow + child task with a freshly-created
// Hatchet worker. Slot count controls parallel task execution; 10
// matches the value used in the Hatchet docs for fan-out examples and
// is sized for the dev laptop (LLM-provider latency dominates anyway).
func NewWorker(client *hatchet.Client, deps Deps, name string) (*Worker, error) {
	defs := Build(client, deps)
	w, err := client.NewWorker(
		name,
		hatchet.WithWorkflows(defs.Workflow, defs.FileTask),
		hatchet.WithSlots(10),
	)
	if err != nil {
		return nil, err
	}
	return &Worker{hatchet: w, defs: defs}, nil
}

// Start runs the worker until ctx is cancelled. Blocks the caller —
// composition.Build runs this in a dedicated goroutine.
func (w *Worker) Start(ctx context.Context) error {
	return w.hatchet.StartBlocking(ctx)
}

// Definitions returns the registered workflow handles. Useful for tests
// that want to invoke the workflow directly without going through the
// enqueuer.
func (w *Worker) Definitions() Definitions { return w.defs }
