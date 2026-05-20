// Package events is the stats bounded context's domain-event adapter.
// It translates typed stats events into messages on a generic event
// broadcaster (SSE today; could fan out to other sinks if the context
// grows new subscribers).
package events

import (
	"context"
	"fmt"

	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
	statsapp "github.com/atilladeniz/next-go-pg/backend/internal/stats/application"
	stats "github.com/atilladeniz/next-go-pg/backend/internal/stats/domain"
)

// Broadcaster is the minimal upstream the publisher needs. The
// platform SSE broker satisfies it without knowing about stats.
type Broadcaster interface {
	Broadcast(eventName, payload string)
}

// Publisher dispatches stats domain events to the broadcaster.
type Publisher struct {
	broadcaster Broadcaster
}

var _ statsapp.DomainEventPublisher = (*Publisher)(nil)

func NewPublisher(broadcaster Broadcaster) *Publisher {
	return &Publisher{broadcaster: broadcaster}
}

// Publish routes each domain event to the right broadcast topic. The
// type switch is exhaustive for THIS context's events; an unknown
// event type is a no-op (the event was probably emitted by another
// context and should be handled there).
func (p *Publisher) Publish(_ context.Context, events ...shared.DomainEvent) error {
	if p.broadcaster == nil {
		return nil
	}
	for _, event := range events {
		switch e := event.(type) {
		case stats.StatIncremented:
			p.broadcaster.Broadcast("stats-updated", fmt.Sprintf(`{"field":"%s"}`, e.Field.String()))
		}
	}
	return nil
}
