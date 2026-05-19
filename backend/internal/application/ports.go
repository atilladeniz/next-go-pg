package application

import (
	"context"

	"github.com/atilladeniz/next-go-pg/backend/internal/domain"
)

// Ports (hexagonal interfaces) consumed by use cases and handlers in
// this package. Concrete implementations live under internal/infrastructure/
// or in jobs/ for the queue-side adapter.

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
// instead of touching gorm.DB directly.
type UserDirectory interface {
	UserByID(ctx context.Context, userID domain.UserID) (*domain.User, error)
	HasKnownDevice(ctx context.Context, userID domain.UserID, userAgent, ipAddress, excludeSessionID string) (bool, error)
}

// JobEnqueuer schedules background jobs by semantic name. Handlers
// depend on this port; the concrete adapter (River-backed) lives in
// the jobs package and is wired by the composition root.
type JobEnqueuer interface {
	EnqueueMagicLink(ctx context.Context, email, url string) error
	EnqueueVerificationEmail(ctx context.Context, email, name, url string) error
	Enqueue2FAOTP(ctx context.Context, email, name, otp string) error
	EnqueueLoginNotification(ctx context.Context, email, userName, device, ipAddress string) error
	EnqueueDataExport(ctx context.Context, jobID, userID, format, dataType string) error
}
