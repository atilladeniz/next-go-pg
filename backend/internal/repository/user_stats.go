package repository

import (
	"github.com/atilladeniz/gocatest/backend/internal/domain"
	"gorm.io/gorm"
)

// UserStatsRepository handles database operations for user stats
type UserStatsRepository struct {
	db *gorm.DB
}

// NewUserStatsRepository creates a new repository instance
func NewUserStatsRepository(db *gorm.DB) *UserStatsRepository {
	return &UserStatsRepository{db: db}
}

// GetOrCreate retrieves stats for a user, creating default stats if not found
func (r *UserStatsRepository) GetOrCreate(userID string) (*domain.UserStats, error) {
	var stats domain.UserStats

	err := r.db.Where("user_id = ?", userID).First(&stats).Error
	if err == gorm.ErrRecordNotFound {
		// Create default stats for new user
		stats = domain.UserStats{
			UserID:        userID,
			ProjectCount:  3,
			ActivityToday: 10,
			Notifications: 2,
		}
		if err := r.db.Create(&stats).Error; err != nil {
			return nil, err
		}
		return &stats, nil
	}
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// Update saves changes to user stats
func (r *UserStatsRepository) Update(stats *domain.UserStats) error {
	return r.db.Save(stats).Error
}

// IncrementField increments a specific field by delta
func (r *UserStatsRepository) IncrementField(userID string, field string, delta int) (*domain.UserStats, error) {
	stats, err := r.GetOrCreate(userID)
	if err != nil {
		return nil, err
	}

	switch field {
	case "projects":
		stats.ProjectCount += delta
		if stats.ProjectCount < 0 {
			stats.ProjectCount = 0
		}
	case "activity":
		stats.ActivityToday += delta
		if stats.ActivityToday < 0 {
			stats.ActivityToday = 0
		}
	case "notifications":
		stats.Notifications += delta
		if stats.Notifications < 0 {
			stats.Notifications = 0
		}
	}

	if err := r.Update(stats); err != nil {
		return nil, err
	}

	return stats, nil
}
