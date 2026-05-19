package application

import (
	"context"

	"github.com/atilladeniz/next-go-pg/backend/internal/domain"
)

// GetUserStats returns (or seeds) the stats row for a user.
type GetUserStats struct {
	Repo StatsRepository
}

func (uc GetUserStats) Execute(ctx context.Context, userID domain.UserID) (*domain.UserStats, error) {
	return uc.Repo.GetOrCreate(ctx, userID)
}

// IncrementStatField bumps one counter for a user, persists the change,
// and dispatches the domain events the aggregate recorded.
type IncrementStatField struct {
	Repo   StatsRepository
	Events DomainEventPublisher
}

func (uc IncrementStatField) Execute(ctx context.Context, userID domain.UserID, field domain.StatField, delta int) (*domain.UserStats, error) {
	stats, err := uc.Repo.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, err
	}
	stats.IncrementField(field, delta)
	if err := uc.Repo.Save(ctx, stats); err != nil {
		return nil, err
	}
	events := stats.PullEvents()
	if uc.Events != nil && len(events) > 0 {
		if err := uc.Events.Publish(ctx, events...); err != nil {
			// Publish failures must not roll back the persisted state;
			// log + continue is the caller's job. We surface the error
			// for observability but the row is already saved.
			return stats, err
		}
	}
	return stats, nil
}
