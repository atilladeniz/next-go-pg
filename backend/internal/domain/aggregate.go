package domain

// DomainEvent is the marker interface every domain event satisfies.
// Aggregates record events on themselves; the application layer pulls
// them after a successful Save and publishes them via the
// application.DomainEventPublisher port.
type DomainEvent interface {
	eventName() string
}

// AggregateRoot is the marker interface for aggregate roots. Embedding
// AggregateBase gives a type the event-recording machinery for free.
type AggregateRoot interface {
	PullEvents() []DomainEvent
}

// AggregateBase is the embeddable base for aggregate roots. It owns
// the pending-event slice and the methods to record + drain it.
type AggregateBase struct {
	pendingEvents []DomainEvent
}

// Record appends an event to the aggregate's pending queue.
func (b *AggregateBase) Record(events ...DomainEvent) {
	b.pendingEvents = append(b.pendingEvents, events...)
}

// PullEvents returns the recorded events and clears the queue. The
// application layer calls this after the aggregate has been persisted.
func (b *AggregateBase) PullEvents() []DomainEvent {
	events := b.pendingEvents
	b.pendingEvents = nil
	return events
}
