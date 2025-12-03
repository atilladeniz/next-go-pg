package domain

import (
	"time"
)

// UserStats represents user statistics stored in the database
type UserStats struct {
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

// TableName specifies the table name for GORM
func (UserStats) TableName() string {
	return "user_stats"
}
