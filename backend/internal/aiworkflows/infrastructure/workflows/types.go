// Package workflows wires the aiworkflows bounded context to the
// Hatchet Go SDK. This package is the ONLY place in the codebase that
// imports `github.com/hatchet-dev/hatchet/sdks/go` — the application
// layer talks to it through the HatchetEnqueuer port.
package workflows

// WorkflowInput is the JSON payload enqueued for one summarize-repo run.
// Carries the minimum the steps need; the aggregate is loaded by
// SummaryID inside the workflow so we don't ship mutable state across
// the network boundary.
type WorkflowInput struct {
	SummaryID uint   `json:"summaryId"`
	UserID    string `json:"userId"`
	RepoURL   string `json:"repoUrl"`
}

// CloneOutput is the result of the clone step. Path is the on-disk
// location of the shallow checkout; the parent workflow run owns the
// cleanup (registered via deferred function in the worker bootstrap).
type CloneOutput struct {
	Path string `json:"path"`
}

// TraverseOutput is the list of repo-relative file paths the
// summarize-file fan-out will iterate over.
type TraverseOutput struct {
	Path  string   `json:"path"`
	Files []string `json:"files"`
}

// SummarizeFileInput is the typed payload for each child `summarize-file`
// task spawned during fan-out.
type SummarizeFileInput struct {
	SummaryID uint   `json:"summaryId"`
	UserID    string `json:"userId"`
	Path      string `json:"path"`
	Filename  string `json:"filename"`
	Total     int    `json:"total"`
}

// SummarizeFileOutput is the produced summary for one file.
type SummarizeFileOutput struct {
	Filename string `json:"filename"`
	Summary  string `json:"summary"`
}

// SummarizeFilesOutput collects all per-file results once the fan-out
// completes.
type SummarizeFilesOutput struct {
	Summaries []SummarizeFileOutput `json:"summaries"`
}

// AggregateOutput is the LLM-produced repo-level summary text.
type AggregateOutput struct {
	Summary string `json:"summary"`
}

// StoreOutput is empty; the persistence step's side effect (RepoSummary
// row updated to `completed`) is the meaningful result.
type StoreOutput struct {
	OK bool `json:"ok"`
}
