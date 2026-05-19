package domain

import "time"

// UserStats is the pure-domain representation of a user's statistics row.
// It has no persistence concerns — GORM-tagged models live in
// internal/infrastructure/persistence and are translated via a mapper.
type UserStats struct {
	ID            uint
	UserID        UserID
	ProjectCount  int
	ActivityToday int
	Notifications int
	LastLogin     time.Time
	MemberSince   time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// IncrementField adjusts the counter selected by field. The result is
// clamped at zero — the domain refuses to represent negative counts.
func (s *UserStats) IncrementField(field StatField, delta int) {
	bump := func(v *int) {
		*v += delta
		if *v < 0 {
			*v = 0
		}
	}
	switch field {
	case StatFieldProjects:
		bump(&s.ProjectCount)
	case StatFieldActivity:
		bump(&s.ActivityToday)
	case StatFieldNotifications:
		bump(&s.Notifications)
	}
}
