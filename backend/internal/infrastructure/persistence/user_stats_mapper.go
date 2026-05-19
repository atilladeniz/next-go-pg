package persistence

import "github.com/atilladeniz/next-go-pg/backend/internal/domain"

func userStatsToDomain(m gormUserStats) domain.UserStats {
	return domain.UserStats{
		ID:            m.ID,
		UserID:        domain.UserID(m.UserID),
		ProjectCount:  m.ProjectCount,
		ActivityToday: m.ActivityToday,
		Notifications: m.Notifications,
		LastLogin:     m.LastLogin,
		MemberSince:   m.MemberSince,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func userStatsFromDomain(d domain.UserStats) gormUserStats {
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
