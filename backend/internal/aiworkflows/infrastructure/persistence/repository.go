package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	aiapp "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/application"
	ai "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/domain"
	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
)

// Repository is the GORM-backed implementation of the aiworkflows
// context's application.Store port.
type Repository struct {
	db *gorm.DB
}

var _ aiapp.Store = (*Repository)(nil)

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a fresh RepoSummary and writes the assigned ID back
// onto the aggregate. The domain's pending events are NOT pulled here;
// the caller (use case) owns event lifecycle.
func (r *Repository) Create(ctx context.Context, agg *ai.RepoSummary) error {
	m := fromDomain(agg)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	agg.ID = m.ID
	agg.CreatedAt = m.CreatedAt
	agg.UpdatedAt = m.UpdatedAt
	return nil
}

// Save persists changes. Mutates fields back onto the aggregate without
// replacing *agg, so AggregateBase.pendingEvents survives.
func (r *Repository) Save(ctx context.Context, agg *ai.RepoSummary) error {
	m := fromDomain(agg)
	if err := r.db.WithContext(ctx).Save(&m).Error; err != nil {
		return err
	}
	agg.UpdatedAt = m.UpdatedAt
	return nil
}

// GetByID returns ErrNotFound when the row is missing.
func (r *Repository) GetByID(ctx context.Context, id uint) (*ai.RepoSummary, error) {
	var m gormRepoSummary
	err := r.db.WithContext(ctx).First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, aiapp.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(m)
}

// Delete removes the row in a single owner-scoped statement. The WHERE
// clause does the auth check inline, so a cross-user request and a
// missing row are indistinguishable on the wire — both return
// ErrNotFound (see Store contract).
func (r *Repository) Delete(ctx context.Context, userID shared.UserID, id uint) error {
	res := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, string(userID)).
		Delete(&gormRepoSummary{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return aiapp.ErrNotFound
	}
	return nil
}

// ListByUserID returns the user's recent summaries, newest first.
func (r *Repository) ListByUserID(ctx context.Context, userID shared.UserID, limit int) ([]*ai.RepoSummary, error) {
	if limit <= 0 {
		limit = 20
	}
	var rows []gormRepoSummary
	err := r.db.WithContext(ctx).
		Where("user_id = ?", string(userID)).
		Order("id DESC").
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]*ai.RepoSummary, 0, len(rows))
	for _, row := range rows {
		agg, err := toDomain(row)
		if err != nil {
			return nil, err
		}
		out = append(out, agg)
	}
	return out, nil
}
