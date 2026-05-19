// Package repository is a transitional facade kept alive while handlers
// still depend on it directly. Phase 3 of the backend-clean-architecture
// refactor introduces application use cases and removes this package.
package repository

import (
	"context"

	"github.com/atilladeniz/next-go-pg/backend/internal/domain"
	"github.com/atilladeniz/next-go-pg/backend/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

// UserStatsRepository delegates persistence to the infrastructure layer
// and adds the soon-to-be-relocated business orchestration (IncrementField).
type UserStatsRepository struct {
	inner *persistence.UserStatsRepository
}

// NewUserStatsRepository constructs the facade over the GORM handle.
func NewUserStatsRepository(db *gorm.DB) *UserStatsRepository {
	return &UserStatsRepository{inner: persistence.NewUserStatsRepository(db)}
}

// GetOrCreate retrieves stats for a user, creating defaults if absent.
func (r *UserStatsRepository) GetOrCreate(userID string) (*domain.UserStats, error) {
	return r.inner.GetOrCreate(context.Background(), userID)
}

// Update saves changes to user stats.
func (r *UserStatsRepository) Update(stats *domain.UserStats) error {
	return r.inner.Save(context.Background(), stats)
}

// IncrementField increments a stat counter by delta and persists the
// change. Business orchestration belongs in a use case — moves there
// in Phase 3.
func (r *UserStatsRepository) IncrementField(userID, field string, delta int) (*domain.UserStats, error) {
	stats, err := r.GetOrCreate(userID)
	if err != nil {
		return nil, err
	}
	parsedField, err := domain.NewStatField(field)
	if err != nil {
		return nil, err
	}
	stats.IncrementField(parsedField, delta)
	if err := r.Update(stats); err != nil {
		return nil, err
	}
	return stats, nil
}
