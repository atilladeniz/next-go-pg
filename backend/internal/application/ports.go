package application

import (
	"context"

	"github.com/atilladeniz/next-go-pg/backend/internal/domain"
)

// Ports (hexagonal interfaces) consumed by use cases in this package.
// Concrete implementations live under internal/infrastructure/.
//
// The user-ID parameter is typed as string for now; it becomes the
// domain.UserID value object in Phase 2 of the refactor.

// StatsRepository persists and retrieves user statistics.
type StatsRepository interface {
	GetOrCreate(ctx context.Context, userID string) (*domain.UserStats, error)
	Save(ctx context.Context, stats *domain.UserStats) error
}

// EventBroadcaster publishes server-sent events to connected clients.
type EventBroadcaster interface {
	Broadcast(eventName, payload string)
}
