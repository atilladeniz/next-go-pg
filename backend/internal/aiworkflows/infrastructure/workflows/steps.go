package workflows

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	hatchet "github.com/hatchet-dev/hatchet/sdks/go"

	aiapp "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/application"
	ai "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/domain"
	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
)

// Deps holds the workflow's runtime dependencies. The dependency graph
// is closed over in step functions when the workflow is built — keeping
// the workflow itself trivially testable in isolation (the Deps struct
// is a single hand-off point for fakes).
type Deps struct {
	Cloner   aiapp.RepoCloner
	LLM      aiapp.LLMClient
	Store    aiapp.Store
	Progress aiapp.ProgressPublisher
	MaxFiles int
	MaxBytes int64
}

// publishStep is a small helper to keep the per-step start/end emissions
// readable. Wrapping in a helper avoids repeating the same five-line
// boilerplate at every step boundary. When `state == completed` and we
// have a duration, we also persist it on the aggregate so refreshes
// after the run keep the timing.
func (d Deps) publishStep(ctx context.Context, in WorkflowInput, name aiapp.StepName, state aiapp.StepState, durationMs int64, reason string) {
	d.Progress.PublishStep(ctx, aiapp.StepProgress{
		SummaryID:  in.SummaryID,
		UserID:     shared.UserID(in.UserID),
		Step:       name,
		State:      state,
		DurationMs: durationMs,
		Reason:     reason,
	})
	if state == aiapp.StepStateCompleted && durationMs > 0 {
		d.recordDuration(ctx, in.SummaryID, string(name), durationMs)
	}
}

// recordDuration persists a completed step's duration onto the
// aggregate. Best-effort: failures here only mean refresh shows ?ms
// instead of the precise time, not a workflow break.
func (d Deps) recordDuration(ctx context.Context, summaryID uint, step string, ms int64) {
	agg, err := d.Store.GetByID(ctx, summaryID)
	if err != nil {
		return
	}
	agg.RecordStepDuration(step, ms)
	_ = d.Store.Save(ctx, agg)
}

// CloneStep performs a shallow clone of the requested repo and marks the
// aggregate as `running`. The output's Path is the on-disk working copy.
func (d Deps) CloneStep(ctx context.Context, in WorkflowInput) (out CloneOutput, err error) {
	start := time.Now()
	d.publishStep(ctx, in, aiapp.StepClone, aiapp.StepStateStarted, 0, "")
	defer func() {
		state := aiapp.StepStateCompleted
		reason := ""
		if err != nil {
			state = aiapp.StepStateFailed
			reason = err.Error()
		}
		d.publishStep(ctx, in, aiapp.StepClone, state, time.Since(start).Milliseconds(), reason)
	}()

	if _, err = d.loadAndStart(ctx, in); err != nil {
		return CloneOutput{}, err
	}
	url, err := ai.NewRepoURL(in.RepoURL)
	if err != nil {
		return CloneOutput{}, fmt.Errorf("clone: invalid repo url: %w", err)
	}
	cloned, err := d.Cloner.Clone(ctx, url)
	if err != nil {
		return CloneOutput{}, fmt.Errorf("clone: %w", err)
	}
	// Cleanup runs in StoreStep at the natural end of the workflow.
	return CloneOutput{Path: cloned.Path}, nil
}

// loadAndStart loads the aggregate, transitions pending → running, and
// persists. Re-runs (retry of clone) tolerate already-running rows.
func (d Deps) loadAndStart(ctx context.Context, in WorkflowInput) (*ai.RepoSummary, error) {
	agg, err := d.Store.GetByID(ctx, in.SummaryID)
	if err != nil {
		return nil, fmt.Errorf("load aggregate: %w", err)
	}
	if agg.Status == ai.StatusPending {
		if err := agg.MarkStarted(time.Now().UTC()); err != nil {
			return nil, fmt.Errorf("mark started: %w", err)
		}
		events := agg.PullEvents()
		if err := d.Store.Save(ctx, agg); err != nil {
			return nil, fmt.Errorf("save after start: %w", err)
		}
		_ = d.Progress.Publish(ctx, events...)
	}
	return agg, nil
}

// TraverseStep walks the cloned repo and selects files to summarize.
// Deterministic, no retries.
func (d Deps) TraverseStep(ctx context.Context, in WorkflowInput, path string) (out TraverseOutput, err error) {
	start := time.Now()
	d.publishStep(ctx, in, aiapp.StepTraverse, aiapp.StepStateStarted, 0, "")
	defer func() {
		state := aiapp.StepStateCompleted
		reason := ""
		if err != nil {
			state = aiapp.StepStateFailed
			reason = err.Error()
		}
		d.publishStep(ctx, in, aiapp.StepTraverse, state, time.Since(start).Milliseconds(), reason)
	}()

	files, err := selectFiles(path, d.MaxFiles, d.MaxBytes)
	if err != nil {
		return TraverseOutput{}, fmt.Errorf("traverse: %w", err)
	}
	return TraverseOutput{Path: path, Files: files}, nil
}

// SummarizeFileStep is the fan-out child task. Called per file by the
// SummarizeFiles orchestrator. Idempotent on the input side: same
// (Path, Filename) always produces the same prompt — actual LLM
// determinism depends on the upstream provider's settings.
//
// Hatchet's SDK validates task function signatures via reflection, so the
// first parameter MUST be hatchet.Context (which embeds context.Context
// anyway — the LLM client treats it as a regular context).
func (d Deps) SummarizeFileStep(ctx hatchet.Context, in SummarizeFileInput) (SummarizeFileOutput, error) {
	full := filepath.Join(in.Path, in.Filename)
	body, err := os.ReadFile(full)
	if err != nil {
		return SummarizeFileOutput{}, fmt.Errorf("read %s: %w", in.Filename, err)
	}
	if int64(len(body)) > d.MaxBytes {
		// Trim huge files so the LLM context window doesn't blow up.
		body = body[:d.MaxBytes]
	}
	prompt := fmt.Sprintf(
		"Summarize the following source file in 2-3 sentences. Focus on what it does, not the syntax.\n\nFILENAME: %s\n\n---\n%s\n---\n\nSUMMARY:",
		in.Filename, string(body),
	)
	summary, err := d.LLM.Generate(ctx, prompt)
	if err != nil {
		return SummarizeFileOutput{}, fmt.Errorf("llm generate: %w", err)
	}
	return SummarizeFileOutput{
		Filename: in.Filename,
		Summary:  strings.TrimSpace(summary),
	}, nil
}

// SummarizeFilesStep fans out across all files via child task calls.
// Each child is independently checkpointed in Hatchet, so a mid-run
// crash resumes from the last in-flight file. As each child completes,
// we immediately publish a `summarize_files:progress` SSE event so the
// frontend's counter advances in real time, rather than only firing the
// final batch after wg.Wait().
func (d Deps) SummarizeFilesStep(
	ctx context.Context,
	hctx hatchet.Context,
	in WorkflowInput,
	traverse TraverseOutput,
	childTask *hatchet.StandaloneTask,
) (out SummarizeFilesOutput, err error) {
	start := time.Now()
	total := len(traverse.Files)
	d.Progress.PublishStep(ctx, aiapp.StepProgress{
		SummaryID: in.SummaryID,
		UserID:    shared.UserID(in.UserID),
		Step:      aiapp.StepSummarizeFiles,
		State:     aiapp.StepStateStarted,
		FileCount: total,
	})
	defer func() {
		state := aiapp.StepStateCompleted
		reason := ""
		if err != nil {
			state = aiapp.StepStateFailed
			reason = err.Error()
		}
		durMs := time.Since(start).Milliseconds()
		d.Progress.PublishStep(ctx, aiapp.StepProgress{
			SummaryID:  in.SummaryID,
			UserID:     shared.UserID(in.UserID),
			Step:       aiapp.StepSummarizeFiles,
			State:      state,
			DurationMs: durMs,
			FileCount:  total,
			Reason:     reason,
		})
		if state == aiapp.StepStateCompleted {
			d.recordDuration(ctx, in.SummaryID, string(aiapp.StepSummarizeFiles), durMs)
		}
	}()

	results := make([]SummarizeFileOutput, total)
	errs := make([]error, total)
	var completed atomic.Int32

	var wg sync.WaitGroup
	wg.Add(total)
	for i, file := range traverse.Files {
		go func(idx int, name string) {
			defer wg.Done()
			res, runErr := childTask.Run(hctx, SummarizeFileInput{
				SummaryID: in.SummaryID,
				UserID:    in.UserID,
				Path:      traverse.Path,
				Filename:  name,
				Total:     total,
			})
			if runErr != nil {
				errs[idx] = runErr
				return
			}
			var typed SummarizeFileOutput
			if decErr := res.Into(&typed); decErr != nil {
				errs[idx] = fmt.Errorf("decode child %q: %w", name, decErr)
				return
			}
			results[idx] = typed

			// Per-file progress event — fires the moment THIS file is
			// summarized, not after wg.Wait(). Counter is the number
			// completed so far (1-based, monotonic).
			n := int(completed.Add(1))
			d.Progress.PublishStep(ctx, aiapp.StepProgress{
				SummaryID: in.SummaryID,
				UserID:    shared.UserID(in.UserID),
				Step:      aiapp.StepSummarizeFiles,
				State:     aiapp.StepStateProgress,
				FileIndex: n,
				FileCount: total,
				Filename:  name,
			})
		}(i, file)
	}
	wg.Wait()

	for _, e := range errs {
		if e != nil {
			return SummarizeFilesOutput{}, e
		}
	}

	// Persist per-file summaries on the aggregate in deterministic order.
	agg, err := d.Store.GetByID(ctx, in.SummaryID)
	if err != nil {
		return SummarizeFilesOutput{}, fmt.Errorf("load aggregate: %w", err)
	}
	for _, r := range results {
		fs, fsErr := ai.NewFileSummary(r.Filename, r.Summary)
		if fsErr != nil {
			return SummarizeFilesOutput{}, fmt.Errorf("file summary value object: %w", fsErr)
		}
		if appendErr := agg.AppendFileSummary(fs, total); appendErr != nil {
			return SummarizeFilesOutput{}, fmt.Errorf("append file: %w", appendErr)
		}
	}
	events := agg.PullEvents()
	if saveErr := d.Store.Save(ctx, agg); saveErr != nil {
		return SummarizeFilesOutput{}, fmt.Errorf("save after fan-out: %w", saveErr)
	}
	_ = d.Progress.Publish(ctx, events...)

	return SummarizeFilesOutput{Summaries: results}, nil
}

// AggregateStep asks the LLM to produce a repo-level summary by stitching
// the per-file summaries into one prompt.
func (d Deps) AggregateStep(ctx context.Context, in WorkflowInput, summaries SummarizeFilesOutput) (out AggregateOutput, err error) {
	start := time.Now()
	d.publishStep(ctx, in, aiapp.StepAggregate, aiapp.StepStateStarted, 0, "")
	defer func() {
		state := aiapp.StepStateCompleted
		reason := ""
		if err != nil {
			state = aiapp.StepStateFailed
			reason = err.Error()
		}
		d.publishStep(ctx, in, aiapp.StepAggregate, state, time.Since(start).Milliseconds(), reason)
	}()

	if len(summaries.Summaries) == 0 {
		return AggregateOutput{}, errors.New("aggregate: empty per-file summaries")
	}
	var b strings.Builder
	b.WriteString("You are summarizing a Git repository. Below are short summaries of individual files. Produce a single 4-6 sentence overview describing what the repository does as a whole.\n\nFILE SUMMARIES:\n")
	for _, s := range summaries.Summaries {
		b.WriteString("- ")
		b.WriteString(s.Filename)
		b.WriteString(": ")
		b.WriteString(s.Summary)
		b.WriteString("\n")
	}
	b.WriteString("\nOVERVIEW:")

	overview, err := d.LLM.Generate(ctx, b.String())
	if err != nil {
		return AggregateOutput{}, fmt.Errorf("llm aggregate: %w", err)
	}
	return AggregateOutput{Summary: strings.TrimSpace(overview)}, nil
}

// StoreStep marks the aggregate as completed and persists the final
// summary. Also cleans up the working copy from disk.
func (d Deps) StoreStep(
	ctx context.Context,
	in WorkflowInput,
	traverse TraverseOutput,
	aggregateOut AggregateOutput,
) (out StoreOutput, err error) {
	start := time.Now()
	d.publishStep(ctx, in, aiapp.StepStore, aiapp.StepStateStarted, 0, "")
	defer func() {
		state := aiapp.StepStateCompleted
		reason := ""
		if err != nil {
			state = aiapp.StepStateFailed
			reason = err.Error()
		}
		d.publishStep(ctx, in, aiapp.StepStore, state, time.Since(start).Milliseconds(), reason)
	}()

	agg, err := d.Store.GetByID(ctx, in.SummaryID)
	if err != nil {
		return StoreOutput{}, fmt.Errorf("load aggregate: %w", err)
	}
	now := time.Now().UTC()
	if err = agg.MarkCompleted(aggregateOut.Summary, now); err != nil {
		return StoreOutput{}, fmt.Errorf("mark completed: %w", err)
	}
	events := agg.PullEvents()
	if err = d.Store.Save(ctx, agg); err != nil {
		return StoreOutput{}, fmt.Errorf("save after complete: %w", err)
	}
	_ = d.Progress.Publish(ctx, events...)

	if traverse.Path != "" {
		_ = os.RemoveAll(traverse.Path)
	}
	return StoreOutput{OK: true}, nil
}

// HandleFailure marks the aggregate as failed and publishes the event.
// Wired to Hatchet's workflow OnFailure hook. Uses context.Background()
// because the hatchet.Context handed to the failure hook may already be
// cancelled by the time we get here — and we still need to write the
// terminal state to the DB regardless.
func (d Deps) HandleFailure(_ context.Context, in WorkflowInput, reason string) {
	ctx := context.Background()
	agg, err := d.Store.GetByID(ctx, in.SummaryID)
	if err != nil {
		return
	}
	if agg.Status.IsTerminal() {
		return
	}
	if err := agg.MarkFailed(reason, time.Now().UTC()); err != nil {
		return
	}
	events := agg.PullEvents()
	if err := d.Store.Save(ctx, agg); err != nil {
		return
	}
	_ = d.Progress.Publish(ctx, events...)
}

// selectFiles walks the cloned repo and returns up to MaxFiles paths to
// summarize. Filters by extension and skips obvious noise (.git, vendored
// node_modules, lockfiles). Deterministic ordering — same repo state
// yields the same list across reruns.
func selectFiles(root string, maxFiles int, maxBytes int64) ([]string, error) {
	if maxFiles <= 0 {
		maxFiles = 25
	}
	include := map[string]struct{}{
		".go": {}, ".ts": {}, ".tsx": {}, ".js": {}, ".jsx": {},
		".py": {}, ".rs": {}, ".java": {}, ".rb": {}, ".sql": {},
		".md": {}, ".yaml": {}, ".yml": {}, ".toml": {},
	}
	skipDir := map[string]struct{}{
		".git": {}, "node_modules": {}, "vendor": {}, "dist": {},
		"build": {}, ".next": {}, "target": {}, "__pycache__": {},
	}
	var picked []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if _, skip := skipDir[info.Name()]; skip {
				return filepath.SkipDir
			}
			return nil
		}
		if maxBytes > 0 && info.Size() > maxBytes {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if _, ok := include[ext]; !ok {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}
		picked = append(picked, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(picked)
	if len(picked) > maxFiles {
		picked = picked[:maxFiles]
	}
	return picked, nil
}
