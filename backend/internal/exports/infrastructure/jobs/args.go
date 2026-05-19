// Package jobs is the exports context's River-side adapter: args,
// worker, and enqueuer for the data-export workflow.
package jobs

import exports "github.com/atilladeniz/next-go-pg/backend/internal/exports/domain"

// DataExportArgs are the queue-shaped inputs to the export worker.
type DataExportArgs struct {
	JobID    string         `json:"jobId"`
	UserID   string         `json:"userId"`
	Format   exports.Format `json:"format"`
	DataType string         `json:"dataType"`
}

func (DataExportArgs) Kind() string { return "data_export" }

// ProgressUpdate is what the worker broadcasts to subscribed clients.
type ProgressUpdate struct {
	JobID      string         `json:"jobId"`
	Status     exports.Status `json:"status"`
	Progress   int            `json:"progress"`
	Message    string         `json:"message"`
	FileName   string         `json:"fileName,omitempty"`
	DownloadID string         `json:"downloadId,omitempty"`
	Error      string         `json:"error,omitempty"`
}
