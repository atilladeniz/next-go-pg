package persistence

import "time"

// gormUserStats is the GORM-tagged persistence representation of
// domain.UserStats. Stays unexported — callers exchange domain types via
// the mapper functions in this package.
type gormUserStats struct {
	ID            uint      `gorm:"primaryKey"`
	UserID        string    `gorm:"uniqueIndex;not null"`
	ProjectCount  int       `gorm:"default:0"`
	ActivityToday int       `gorm:"default:0"`
	Notifications int       `gorm:"default:0"`
	LastLogin     time.Time `gorm:"autoUpdateTime"`
	MemberSince   time.Time `gorm:"autoCreateTime"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (gormUserStats) TableName() string { return "user_stats" }
