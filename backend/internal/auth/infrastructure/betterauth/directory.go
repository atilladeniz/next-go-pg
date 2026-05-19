// Package betterauth is the auth bounded context's adapter against
// Better Auth's tables (user, session). The schema is owned by Better
// Auth — this code only queries it.
package betterauth

import (
	"context"

	"gorm.io/gorm"

	authapp "github.com/atilladeniz/next-go-pg/backend/internal/auth/application"
	auth "github.com/atilladeniz/next-go-pg/backend/internal/auth/domain"
	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
)

// Directory reads users and session metadata from Better Auth tables.
type Directory struct {
	db *gorm.DB
}

var _ authapp.UserDirectory = (*Directory)(nil)

func NewDirectory(db *gorm.DB) *Directory {
	return &Directory{db: db}
}

type userRow struct {
	Email string `gorm:"column:email"`
	Name  string `gorm:"column:name"`
}

// UserByID returns the user record from Better Auth's user table.
func (d *Directory) UserByID(ctx context.Context, userID shared.UserID) (*auth.User, error) {
	var row userRow
	if err := d.db.WithContext(ctx).Table("user").Where("id = ?", string(userID)).First(&row).Error; err != nil {
		return nil, err
	}
	return &auth.User{ID: userID, Email: row.Email, Name: row.Name}, nil
}

// HasKnownDevice reports whether the user already has a session from
// the same user agent / IP, excluding the given session ID.
func (d *Directory) HasKnownDevice(ctx context.Context, userID shared.UserID, userAgent, ipAddress, excludeSessionID string) (bool, error) {
	var count int64
	err := d.db.WithContext(ctx).
		Table("session").
		Where(`"userId" = ? AND "userAgent" = ? AND "ipAddress" = ? AND id != ?`,
			string(userID), userAgent, ipAddress, excludeSessionID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
