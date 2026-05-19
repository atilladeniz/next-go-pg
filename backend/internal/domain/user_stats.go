package domain

import (
	"fmt"
	"time"
)

// StatField enumerates the user-stats counters that can be incremented.
type StatField int

const (
	StatFieldUnknown StatField = iota
	StatFieldProjects
	StatFieldActivity
	StatFieldNotifications
)

// NewStatField parses the wire-level field name into a typed value.
// Unknown names are rejected.
func NewStatField(s string) (StatField, error) {
	switch s {
	case "projects":
		return StatFieldProjects, nil
	case "activity":
		return StatFieldActivity, nil
	case "notifications":
		return StatFieldNotifications, nil
	default:
		return StatFieldUnknown, fmt.Errorf("unknown stat field %q (want one of: projects, activity, notifications)", s)
	}
}

// String returns the wire-level name.
func (f StatField) String() string {
	switch f {
	case StatFieldProjects:
		return "projects"
	case StatFieldActivity:
		return "activity"
	case StatFieldNotifications:
		return "notifications"
	default:
		return "unknown"
	}
}

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
