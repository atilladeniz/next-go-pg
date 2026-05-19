package persistence

import (
	"context"

	"github.com/atilladeniz/next-go-pg/backend/internal/application"
	"github.com/atilladeniz/next-go-pg/backend/internal/domain"
	"gorm.io/gorm"
)

// UserDirectoryRepository reads users and session metadata from the
// Better Auth tables. Schema is owned by Better Auth — this adapter
// only queries it.
type UserDirectoryRepository struct {
	db *gorm.DB
}

var _ application.UserDirectory = (*UserDirectoryRepository)(nil)

func NewUserDirectoryRepository(db *gorm.DB) *UserDirectoryRepository {
	return &UserDirectoryRepository{db: db}
}

type betterAuthUserRow struct {
	Email string `gorm:"column:email"`
	Name  string `gorm:"column:name"`
}

// UserByID returns the user record from the Better Auth `user` table.
func (r *UserDirectoryRepository) UserByID(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	var row betterAuthUserRow
	if err := r.db.WithContext(ctx).Table("user").Where("id = ?", string(userID)).First(&row).Error; err != nil {
		return nil, err
	}
	return &domain.User{ID: userID, Email: row.Email, Name: row.Name}, nil
}

// HasKnownDevice reports whether the user already has a session from
// the same user agent / IP, excluding the given session ID.
func (r *UserDirectoryRepository) HasKnownDevice(ctx context.Context, userID domain.UserID, userAgent, ipAddress, excludeSessionID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("session").
		Where(`"userId" = ? AND "userAgent" = ? AND "ipAddress" = ? AND id != ?`,
			string(userID), userAgent, ipAddress, excludeSessionID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
