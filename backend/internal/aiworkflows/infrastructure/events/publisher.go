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

// progressPayload is the SSE event body. Frontend consumers decode it
// directly. `step` identifies which domain event fired so the UI can
// pick a label and render fileIndex/fileCount when present.
type progressPayload struct {
	SummaryID uint   `json:"summaryId"`
	UserID    string `json:"userId"`
	Step      string `json:"step"`
	Status    string `json:"status,omitempty"`
	Filename  string `json:"filename,omitempty"`
	FileIndex int    `json:"fileIndex,omitempty"`
	FileCount int    `json:"fileCount,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

const sseEventName = "ai-progress"

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
			SummaryID: e.SummaryID,
			UserID:    e.UserID.String(),
			Step:      "started",
			Status:    "running",
		}, true
	case ai.FileSummarized:
		return progressPayload{
			SummaryID: e.SummaryID,
			UserID:    e.UserID.String(),
			Step:      "summarize_file",
			Filename:  e.Filename,
			FileIndex: e.FileIndex,
			FileCount: e.FileCount,
		}, true
	case ai.SummaryCompleted:
		return progressPayload{
			SummaryID: e.SummaryID,
			UserID:    e.UserID.String(),
			Step:      "store",
			Status:    "completed",
		}, true
	case ai.SummaryFailed:
		return progressPayload{
			SummaryID: e.SummaryID,
			UserID:    e.UserID.String(),
			Step:      "store",
			Status:    "failed",
			Reason:    e.Reason,
		}, true
	case ai.SummaryCancelled:
		return progressPayload{
			SummaryID: e.SummaryID,
			UserID:    e.UserID.String(),
			Step:      "store",
			Status:    "cancelled",
		}, true
	default:
		return progressPayload{}, false
	}
}
