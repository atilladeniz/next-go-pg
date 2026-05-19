package application

import (
	"context"

	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
	stats "github.com/atilladeniz/next-go-pg/backend/internal/stats/domain"
)

// GetUserStats returns (or seeds) the stats row for a user.
type GetUserStats struct {
	Repo Repository
}

func (uc GetUserStats) Execute(ctx context.Context, userID shared.UserID) (*stats.UserStats, error) {
	return uc.Repo.GetOrCreate(ctx, userID)
}

// IncrementStatField bumps one counter for a user, persists the change,
// and dispatches the domain events the aggregate recorded.
type IncrementStatField struct {
	Repo   Repository
	Events DomainEventPublisher
}

func (uc IncrementStatField) Execute(ctx context.Context, userID shared.UserID, field stats.StatField, delta int) (*stats.UserStats, error) {
	agg, err := uc.Repo.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, err
	}
	agg.IncrementField(field, delta)
	if err := uc.Repo.Save(ctx, agg); err != nil {
		return nil, err
	}
	events := agg.PullEvents()
	if uc.Events != nil && len(events) > 0 {
		if err := uc.Events.Publish(ctx, events...); err != nil {
			return agg, err
		}
	}
	return agg, nil
}
