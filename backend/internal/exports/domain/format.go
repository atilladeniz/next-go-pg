// Package domain is the exports bounded context's model. The model is
// small (today): a Format enum and a Status enum. Real export data
// lives elsewhere — this context only owns the lifecycle of an export
// job and its rendered output.
package domain

// Format enumerates the supported export file formats.
type Format string

const (
	FormatCSV  Format = "csv"
	FormatJSON Format = "json"
)

// Status reports the lifecycle of an export job.
type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)
