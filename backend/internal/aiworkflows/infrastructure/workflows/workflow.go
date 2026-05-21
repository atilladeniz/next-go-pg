package workflows

import (
	hatchet "github.com/hatchet-dev/hatchet/sdks/go"
)

const (
	WorkflowName       = "summarize-repo"
	StandaloneFileTask = "summarize-file"
)

// Definitions bundles the workflow handle and the child task handle so
// the worker bootstrap can register both with one call.
type Definitions struct {
	Workflow *hatchet.Workflow
	FileTask *hatchet.StandaloneTask
}

// Build wires the DAG: clone → traverse → summarize-files → aggregate → store.
// The fan-out child `summarize-file` is registered as a separate
// StandaloneTask so each per-file call gets its own checkpoint and its
// own retry policy.
func Build(client *hatchet.Client, deps Deps) Definitions {
	// Child task: one Hatchet run per file. 5× retry with exponential
	// backoff because the LLM gateway can be transiently slow, rate-
	// limited, or fail upstream.
	fileTask := client.NewStandaloneTask(
		StandaloneFileTask,
		deps.SummarizeFileStep,
		hatchet.WithRetries(5),
		hatchet.WithRetryBackoff(2, 60),
	)

	wf := client.NewWorkflow(WorkflowName)

	cloneT := wf.NewTask(
		"clone",
		func(ctx hatchet.Context, in WorkflowInput) (CloneOutput, error) {
			return deps.CloneStep(ctx, in)
		},
		hatchet.WithRetries(3),
	)

	traverseT := wf.NewTask(
		"traverse",
		func(ctx hatchet.Context, in WorkflowInput) (TraverseOutput, error) {
			var clone CloneOutput
			if err := ctx.ParentOutput(cloneT, &clone); err != nil {
				return TraverseOutput{}, err
			}
			return deps.TraverseStep(ctx, in, clone.Path)
		},
		hatchet.WithParents(cloneT),
		// No WithRetries — Traverse is pure/deterministic. Any error here
		// is a real bug or filesystem fault, not transient.
	)

	summarizeT := wf.NewTask(
		"summarize-files",
		func(ctx hatchet.Context, in WorkflowInput) (SummarizeFilesOutput, error) {
			var traverse TraverseOutput
			if err := ctx.ParentOutput(traverseT, &traverse); err != nil {
				return SummarizeFilesOutput{}, err
			}
			return deps.SummarizeFilesStep(ctx, ctx, in, traverse, fileTask)
		},
		hatchet.WithParents(traverseT),
		hatchet.WithRetries(3),
	)

	aggregateT := wf.NewTask(
		"aggregate",
		func(ctx hatchet.Context, in WorkflowInput) (AggregateOutput, error) {
			var summaries SummarizeFilesOutput
			if err := ctx.ParentOutput(summarizeT, &summaries); err != nil {
				return AggregateOutput{}, err
			}
			return deps.AggregateStep(ctx, in, summaries)
		},
		hatchet.WithParents(summarizeT),
		hatchet.WithRetries(3),
	)

	_ = wf.NewTask(
		"store",
		func(ctx hatchet.Context, in WorkflowInput) (StoreOutput, error) {
			var traverse TraverseOutput
			if err := ctx.ParentOutput(traverseT, &traverse); err != nil {
				return StoreOutput{}, err
			}
			var aggregateOut AggregateOutput
			if err := ctx.ParentOutput(aggregateT, &aggregateOut); err != nil {
				return StoreOutput{}, err
			}
			return deps.StoreStep(ctx, in, traverse, aggregateOut)
		},
		hatchet.WithParents(aggregateT),
		hatchet.WithRetries(3),
	)

	// OnFailure hook: any unrecoverable step error funnels through here
	// so the aggregate doesn't stay stuck in `running`.
	wf.OnFailure(func(ctx hatchet.Context, in WorkflowInput) (struct{}, error) {
		deps.HandleFailure(ctx, in, "workflow failure")
		return struct{}{}, nil
	})

	return Definitions{Workflow: wf, FileTask: fileTask}
}
