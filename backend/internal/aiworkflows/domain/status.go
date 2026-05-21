package domain

import "fmt"

// Status is the lifecycle state of a RepoSummary run. The set is closed:
// any value outside the constants below is invalid and rejected by
// NewStatus at the persistence boundary.
type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

// NewStatus parses a wire-level status string and rejects unknown values.
func NewStatus(s string) (Status, error) {
	switch Status(s) {
	case StatusPending, StatusRunning, StatusCompleted, StatusFailed, StatusCancelled:
		return Status(s), nil
	default:
		return "", fmt.Errorf("unknown status %q", s)
	}
}

func (s Status) String() string { return string(s) }

// IsTerminal returns true when the run can no longer transition.
func (s Status) IsTerminal() bool {
	return s == StatusCompleted || s == StatusFailed || s == StatusCancelled
}
