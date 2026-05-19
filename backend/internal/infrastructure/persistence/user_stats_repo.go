package persistence

import (
	"context"
	"errors"

	"github.com/atilladeniz/next-go-pg/backend/internal/application"
	"github.com/atilladeniz/next-go-pg/backend/internal/domain"
	"gorm.io/gorm"
)

// UserStatsRepository is the GORM-backed implementation of
// application.StatsRepository. It is the only place that knows about
// the gormUserStats persistence model.
type UserStatsRepository struct {
	db *gorm.DB
}

// Compile-time guarantee that the implementation satisfies the port.
var _ application.StatsRepository = (*UserStatsRepository)(nil)

// NewUserStatsRepository constructs a repository over the given GORM
// handle. The database is expected to already have the schema in place
// (AutoMigrate runs from the composition root).
func NewUserStatsRepository(db *gorm.DB) *UserStatsRepository {
	return &UserStatsRepository{db: db}
}

// GetOrCreate returns the user's stats row, creating it with the seeded
// defaults if it does not yet exist.
func (r *UserStatsRepository) GetOrCreate(ctx context.Context, userID domain.UserID) (*domain.UserStats, error) {
	var g gormUserStats
	err := r.db.WithContext(ctx).Where("user_id = ?", string(userID)).First(&g).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		g = gormUserStats{
			UserID:        string(userID),
			ProjectCount:  3,
			ActivityToday: 10,
			Notifications: 2,
		}
		if err := r.db.WithContext(ctx).Create(&g).Error; err != nil {
			return nil, err
		}
		d := userStatsToDomain(g)
		return &d, nil
	}
	if err != nil {
		return nil, err
	}
	d := userStatsToDomain(g)
	return &d, nil
}

// Save persists changes to the stats row.
func (r *UserStatsRepository) Save(ctx context.Context, stats *domain.UserStats) error {
	g := userStatsFromDomain(*stats)
	if err := r.db.WithContext(ctx).Save(&g).Error; err != nil {
		return err
	}
	// Reflect GORM-managed timestamp updates back into the caller's
	// domain value so subsequent reads see the same data.
	*stats = userStatsToDomain(g)
	return nil
}
