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

// IncrementStatField bumps one counter for a user and broadcasts the
// change. Domain invariants (clamping) live on UserStats.IncrementField.
type IncrementStatField struct {
	Repo   StatsRepository
	Events EventBroadcaster
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
	if uc.Events != nil {
		uc.Events.Broadcast("stats-updated", `{"field":"`+field.String()+`"}`)
	}
	return stats, nil
}
