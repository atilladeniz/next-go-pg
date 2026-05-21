// Package events is the aiworkflows domain-event adapter. It translates
// the bounded context's domain events into messages on the platform SSE
// broker so the frontend can render progress without polling.
package events

import (
	"context"
	"encoding/json"

	aiapp "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/application"
	ai "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/domain"
	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
	"github.com/atilladeniz/next-go-pg/backend/pkg/metrics"
)

// Broadcaster is the minimal upstream this publisher needs. The platform
// SSE broker satisfies it without knowing about aiworkflows.
type Broadcaster interface {
	Broadcast(eventName, payload string)
}

// Publisher dispatches aiworkflows domain events to the broadcaster.
// All events flow on a single SSE event name `ai-progress` so the
// frontend has one stream to listen on regardless of step.
type Publisher struct {
	broadcaster Broadcaster
}

var _ aiapp.ProgressPublisher = (*Publisher)(nil)

func NewPublisher(broadcaster Broadcaster) *Publisher {
	return &Publisher{broadcaster: broadcaster}
}

// progressPayload is the SSE event body. Two flavours flow over the same
// `ai-progress` channel — `kind` distinguishes them on the frontend:
//   - kind=lifecycle: started/completed/failed/cancelled from the
//     RepoSummary aggregate's domain events
//   - kind=step: step-level transitions emitted directly by the
//     workflow (clone/traverse/.../store with started/completed/failed/progress)
type progressPayload struct {
	Kind       string `json:"kind"`
	SummaryID  uint   `json:"summaryId"`
	UserID     string `json:"userId"`
	Step       string `json:"step"`
	State      string `json:"state,omitempty"`      // step transitions: started/completed/failed/progress
	Status     string `json:"status,omitempty"`     // lifecycle (run-level): running/completed/failed/cancelled
	DurationMs int64  `json:"durationMs,omitempty"` // step-level duration on completion
	Filename   string `json:"filename,omitempty"`
	FileIndex  int    `json:"fileIndex,omitempty"`
	FileCount  int    `json:"fileCount,omitempty"`
	Reason     string `json:"reason,omitempty"`
}

const sseEventName = "ai-progress"

// PublishStep broadcasts a step-level progress event. Step events are
// the granular update the live UI renders against (one row per step,
// with state + duration + per-file counters during fan-out).
func (p *Publisher) PublishStep(_ context.Context, step aiapp.StepProgress) {
	if p.broadcaster == nil {
		return
	}
	payload := progressPayload{
		Kind:       "step",
		SummaryID:  step.SummaryID,
		UserID:     step.UserID.String(),
		Step:       string(step.Step),
		State:      string(step.State),
		DurationMs: step.DurationMs,
		Filename:   step.Filename,
		FileIndex:  step.FileIndex,
		FileCount:  step.FileCount,
		Reason:     step.Reason,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return
	}
	p.broadcaster.Broadcast(sseEventName, string(raw))
}

// Publish routes each domain event to the broadcast topic. Unknown
// events are silently skipped (probably emitted by another context).
func (p *Publisher) Publish(_ context.Context, events ...shared.DomainEvent) error {
	if p.broadcaster == nil {
		return nil
	}
	for _, ev := range events {
		payload, ok := encode(ev)
		if !ok {
			continue
		}
		raw, err := json.Marshal(payload)
		if err != nil {
			continue
		}
		p.broadcaster.Broadcast(sseEventName, string(raw))
		incrementTerminalCounter(ev)
	}
	return nil
}

// incrementTerminalCounter bumps ai_workflows_completed_total on the
// three terminal events. Non-terminal events (Started, FileSummarized)
// do not touch the counter.
func incrementTerminalCounter(ev shared.DomainEvent) {
	switch ev.(type) {
	case ai.SummaryCompleted:
		metrics.AIWorkflowsCompleted.WithLabelValues("success").Inc()
	case ai.SummaryFailed:
		metrics.AIWorkflowsCompleted.WithLabelValues("failed").Inc()
	case ai.SummaryCancelled:
		metrics.AIWorkflowsCompleted.WithLabelValues("cancelled").Inc()
	}
}

func encode(ev shared.DomainEvent) (progressPayload, bool) {
	switch e := ev.(type) {
	case ai.SummaryStarted:
		return progressPayload{
			Kind:      "lifecycle",
			SummaryID: e.SummaryID,
			UserID:    e.UserID.String(),
			Status:    "running",
		}, true
	case ai.FileSummarized:
		// FileSummarized still flows on lifecycle events for downstream
		// consumers (audit, future analytics). The frontend's live UI
		// uses the granular PublishStep stream instead.
		return progressPayload{
			Kind:      "lifecycle",
			SummaryID: e.SummaryID,
			UserID:    e.UserID.String(),
			Filename:  e.Filename,
			FileIndex: e.FileIndex,
			FileCount: e.FileCount,
		}, true
	case ai.SummaryCompleted:
		return progressPayload{
			Kind:      "lifecycle",
			SummaryID: e.SummaryID,
			UserID:    e.UserID.String(),
			Status:    "completed",
		}, true
	case ai.SummaryFailed:
		return progressPayload{
			Kind:      "lifecycle",
			SummaryID: e.SummaryID,
			UserID:    e.UserID.String(),
			Status:    "failed",
			Reason:    e.Reason,
		}, true
	case ai.SummaryCancelled:
		return progressPayload{
			Kind:      "lifecycle",
			SummaryID: e.SummaryID,
			UserID:    e.UserID.String(),
			Status:    "cancelled",
		}, true
	default:
		return progressPayload{}, false
	}
}
