package persistence

import (
	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
	stats "github.com/atilladeniz/next-go-pg/backend/internal/stats/domain"
)

func toDomain(m gormUserStats) stats.UserStats {
	return stats.UserStats{
		ID:            m.ID,
		UserID:        shared.UserID(m.UserID),
		ProjectCount:  m.ProjectCount,
		ActivityToday: m.ActivityToday,
		Notifications: m.Notifications,
		LastLogin:     m.LastLogin,
		MemberSince:   m.MemberSince,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func fromDomain(d stats.UserStats) gormUserStats {
	return gormUserStats{
		ID:            d.ID,
		UserID:        string(d.UserID),
		ProjectCount:  d.ProjectCount,
		ActivityToday: d.ActivityToday,
		Notifications: d.Notifications,
		LastLogin:     d.LastLogin,
		MemberSince:   d.MemberSince,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}
