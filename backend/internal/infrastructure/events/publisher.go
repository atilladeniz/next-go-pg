// Package events is the concrete adapter for application.DomainEventPublisher.
// It translates typed domain events into messages on the SSE event
// broadcaster (and could fan out to audit/log/metrics sinks as the
// domain grows).
package events

import (
	"context"
	"fmt"

	"github.com/atilladeniz/next-go-pg/backend/internal/application"
	"github.com/atilladeniz/next-go-pg/backend/internal/domain"
)

// Publisher dispatches domain events to one or more sinks. Today the
// only sink is the SSE broadcaster; new sinks (audit log, message bus)
// would be added as fields here.
type Publisher struct {
	broadcaster application.EventBroadcaster
}

var _ application.DomainEventPublisher = (*Publisher)(nil)

// NewPublisher builds a Publisher over the given broadcaster.
func NewPublisher(broadcaster application.EventBroadcaster) *Publisher {
	return &Publisher{broadcaster: broadcaster}
}

// Publish routes each domain event to the right broadcast topic. New
// event types get a case here; the compiler doesn't enforce
// exhaustiveness for interface type switches, but the absence of a
// case is a no-op — explicit failure-by-omission, not silent drop.
func (p *Publisher) Publish(_ context.Context, events ...domain.DomainEvent) error {
	for _, event := range events {
		switch e := event.(type) {
		case domain.StatIncremented:
			if p.broadcaster != nil {
				p.broadcaster.Broadcast("stats-updated", fmt.Sprintf(`{"field":"%s"}`, e.Field.String()))
			}
		}
	}
	return nil
}
