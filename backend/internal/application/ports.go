package application

import (
	"context"

	"github.com/atilladeniz/next-go-pg/backend/internal/domain"
)

// Ports (hexagonal interfaces) consumed by use cases in this package.
// Concrete implementations live under internal/infrastructure/.

// StatsRepository persists and retrieves user statistics.
type StatsRepository interface {
	GetOrCreate(ctx context.Context, userID domain.UserID) (*domain.UserStats, error)
	Save(ctx context.Context, stats *domain.UserStats) error
}

// EventBroadcaster publishes server-sent events to connected clients.
type EventBroadcaster interface {
	Broadcast(eventName, payload string)
}

// UserDirectory reads user records owned by an external auth provider
// (Better Auth at the time of writing). Webhook handlers consume this
// instead of touching `gorm.DB` directly.
type UserDirectory interface {
	UserByID(ctx context.Context, userID domain.UserID) (*domain.User, error)
	HasKnownDevice(ctx context.Context, userID domain.UserID, userAgent, ipAddress, excludeSessionID string) (bool, error)
}
