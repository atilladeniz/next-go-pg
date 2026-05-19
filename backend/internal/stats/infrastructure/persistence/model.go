// Package persistence holds the GORM-backed adapter for the stats
// bounded context. The gormUserStats twin is intentionally unexported —
// callers exchange domain types via the mapper functions.
package persistence

import "time"

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

// Entities returns the GORM models that AutoMigrate must process for
// the stats context. The composition root collects entities from every
// context that owns persistence and feeds them to AutoMigrate.
func Entities() []any {
	return []any{&gormUserStats{}}
}
