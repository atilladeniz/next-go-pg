// Package application is the exports bounded context's use-case layer.
// Adapters live under internal/exports/infrastructure/.
package application

import (
	"context"
	"time"
)

// Store holds completed export results (in-memory cache today;
// could become object storage tomorrow). Adapters live in
// internal/exports/infrastructure/.
type Store interface {
	Save(id string, result *Result)
	Get(id string) (*Result, bool)
	Delete(id string)
}

// Result is the rendered export, ready for download.
type Result struct {
	Data        []byte
	ContentType string
	FileName    string
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

// ProgressPublisher emits per-job progress events to clients (SSE).
// The exports context broadcasts its own "export-progress" topic; it
// does not share the stats context's domain event publisher because
// progress is an infrastructure-level fan-out, not a domain event.
type ProgressPublisher interface {
	Broadcast(eventName, payload string)
}

// JobEnqueuer schedules export work. Only one method today.
type JobEnqueuer interface {
	EnqueueDataExport(ctx context.Context, jobID, userID, format, dataType string) error
}

// StatsSnapshot is exports' view of another context's data. Strict-DDD
// inter-context communication: the exports context declares what it
// needs; the composition root wires an adapter against the stats
// context's repository. No direct import of stats internals.
type StatsSnapshot struct {
	Projects      int
	Activity      int
	Notifications int
}

// StatsReader is exports' anti-corruption layer over the stats
// context. The implementation lives in the composition root.
type StatsReader interface {
	Read(ctx context.Context, userID string) (StatsSnapshot, error)
}
