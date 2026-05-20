package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
	statsapp "github.com/atilladeniz/next-go-pg/backend/internal/stats/application"
	stats "github.com/atilladeniz/next-go-pg/backend/internal/stats/domain"
)

// Repository is the GORM-backed implementation of the stats context's
// application.Repository port.
type Repository struct {
	db *gorm.DB
}

var _ statsapp.Repository = (*Repository)(nil)

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetOrCreate returns the user's stats row. If absent, it asks the
// domain to build a freshly-seeded aggregate.
func (r *Repository) GetOrCreate(ctx context.Context, userID shared.UserID) (*stats.UserStats, error) {
	var m gormUserStats
	err := r.db.WithContext(ctx).Where("user_id = ?", string(userID)).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fresh := stats.NewUserStats(userID)
		m = fromDomain(*fresh)
		if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
			return nil, err
		}
		d := toDomain(m)
		return &d, nil
	}
	if err != nil {
		return nil, err
	}
	d := toDomain(m)
	return &d, nil
}

// Save persists changes and reflects GORM-managed timestamps back into
// the caller's domain value WITHOUT replacing the aggregate as a whole —
// that would wipe AggregateBase.pendingEvents and the use case would
// never see the events the aggregate just recorded. Mutate the fields
// the database owns; leave everything else (including pending events)
// alone.
func (r *Repository) Save(ctx context.Context, agg *stats.UserStats) error {
	m := fromDomain(*agg)
	if err := r.db.WithContext(ctx).Save(&m).Error; err != nil {
		return err
	}
	agg.ID = m.ID
	agg.LastLogin = m.LastLogin
	agg.MemberSince = m.MemberSince
	agg.CreatedAt = m.CreatedAt
	agg.UpdatedAt = m.UpdatedAt
	return nil
}
