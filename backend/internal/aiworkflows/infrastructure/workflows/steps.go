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

// CloneStep performs a shallow clone of the requested repo and marks the
// aggregate as `running`. The output's Path is the on-disk working copy.
func (d Deps) CloneStep(ctx context.Context, in WorkflowInput) (CloneOutput, error) {
	agg, err := d.loadAndStart(ctx, in)
	if err != nil {
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
	// Cleanup is the worker's responsibility on workflow termination —
	// we intentionally do NOT defer cleanup here because the path needs
	// to survive into the next step. The store step performs cleanup as
	// its final side effect.
	_ = agg
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
func (d Deps) TraverseStep(_ context.Context, path string) (TraverseOutput, error) {
	files, err := selectFiles(path, d.MaxFiles, d.MaxBytes)
	if err != nil {
		return TraverseOutput{}, fmt.Errorf("traverse: %w", err)
	}
	return TraverseOutput{Path: path, Files: files}, nil
}

// SummarizeFileStep is the fan-out child task. Called per file by the
// SummarizeFiles orchestrator. Idempotent: same (Path, Filename, model)
// yields the same output if Ollama is configured deterministically. The
// signature uses hatchet.Context (not plain context.Context) because the
// Hatchet SDK validates task function signatures at registration time.
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
// crash resumes from the last in-flight file. We append per-file
// summaries to the aggregate as they arrive and publish progress events.
func (d Deps) SummarizeFilesStep(
	ctx context.Context,
	hctx hatchet.Context,
	in WorkflowInput,
	traverse TraverseOutput,
	childTask *hatchet.StandaloneTask,
) (SummarizeFilesOutput, error) {
	total := len(traverse.Files)
	results := make([]SummarizeFileOutput, total)
	errs := make([]error, total)

	var wg sync.WaitGroup
	wg.Add(total)
	for i, file := range traverse.Files {
		go func(idx int, name string) {
			defer wg.Done()
			out, err := childTask.Run(hctx, SummarizeFileInput{
				SummaryID: in.SummaryID,
				UserID:    in.UserID,
				Path:      traverse.Path,
				Filename:  name,
				Total:     total,
			})
			if err != nil {
				errs[idx] = err
				return
			}
			var typed SummarizeFileOutput
			if err := out.Into(&typed); err != nil {
				errs[idx] = fmt.Errorf("decode child %q: %w", name, err)
				return
			}
			results[idx] = typed
		}(i, file)
	}
	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return SummarizeFilesOutput{}, err
		}
	}

	// Append in deterministic order (input order is preserved by index)
	// and emit per-file progress events.
	agg, err := d.Store.GetByID(ctx, in.SummaryID)
	if err != nil {
		return SummarizeFilesOutput{}, fmt.Errorf("load aggregate: %w", err)
	}
	for _, r := range results {
		fs, err := ai.NewFileSummary(r.Filename, r.Summary)
		if err != nil {
			return SummarizeFilesOutput{}, fmt.Errorf("file summary value object: %w", err)
		}
		if err := agg.AppendFileSummary(fs, total); err != nil {
			return SummarizeFilesOutput{}, fmt.Errorf("append file: %w", err)
		}
	}
	events := agg.PullEvents()
	if err := d.Store.Save(ctx, agg); err != nil {
		return SummarizeFilesOutput{}, fmt.Errorf("save after fan-out: %w", err)
	}
	_ = d.Progress.Publish(ctx, events...)

	return SummarizeFilesOutput{Summaries: results}, nil
}

// AggregateStep asks the LLM to produce a repo-level summary by stitching
// the per-file summaries into one prompt.
func (d Deps) AggregateStep(ctx context.Context, summaries SummarizeFilesOutput) (AggregateOutput, error) {
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
) (StoreOutput, error) {
	agg, err := d.Store.GetByID(ctx, in.SummaryID)
	if err != nil {
		return StoreOutput{}, fmt.Errorf("load aggregate: %w", err)
	}
	now := time.Now().UTC()
	if err := agg.MarkCompleted(aggregateOut.Summary, now); err != nil {
		return StoreOutput{}, fmt.Errorf("mark completed: %w", err)
	}
	events := agg.PullEvents()
	if err := d.Store.Save(ctx, agg); err != nil {
		return StoreOutput{}, fmt.Errorf("save after complete: %w", err)
	}
	_ = d.Progress.Publish(ctx, events...)

	if traverse.Path != "" {
		_ = os.RemoveAll(traverse.Path)
	}
	return StoreOutput{OK: true}, nil
}

// HandleFailure marks the aggregate as failed and publishes the event.
// Wired to Hatchet's workflow OnFailure hook.
func (d Deps) HandleFailure(ctx context.Context, in WorkflowInput, reason string) {
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
			return nil // skip unreadable entries
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

// usedShared is here to keep go-vet happy when shared is otherwise unused.
var _ = shared.NewUserID
