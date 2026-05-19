// Package application is the stats bounded context's use-case layer.
// Ports declared here are consumed by use cases and implemented by
// adapters under internal/stats/infrastructure/.
package application

import (
	"context"

	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
	stats "github.com/atilladeniz/next-go-pg/backend/internal/stats/domain"
)

// Repository persists and retrieves user statistics aggregates.
type Repository interface {
	GetOrCreate(ctx context.Context, userID shared.UserID) (*stats.UserStats, error)
	Save(ctx context.Context, agg *stats.UserStats) error
}

// DomainEventPublisher dispatches the domain events an aggregate
// recorded during a use case. The adapter routes events to whatever
// fan-out medium is appropriate (SSE, audit log, message bus).
type DomainEventPublisher interface {
	Publish(ctx context.Context, events ...shared.DomainEvent) error
}
