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

// Default seed values for a freshly created UserStats aggregate. Kept
// here (in the domain) rather than in the repository so the invariant
// "a new user starts with these counts" is owned by the model.
const (
	defaultProjectCount  = 3
	defaultActivityToday = 10
	defaultNotifications = 2
)

// UserStats is the aggregate root for a user's statistics. It carries
// the counter invariants and records domain events on mutation. The
// GORM-tagged twin lives in internal/infrastructure/persistence.
type UserStats struct {
	AggregateBase

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

// Compile-time aggregate check.
var _ AggregateRoot = (*UserStats)(nil)

// NewUserStats is the aggregate factory for a freshly-seeded row.
// Persistence adapters call this when GetOrCreate sees no existing row.
func NewUserStats(userID UserID) *UserStats {
	return &UserStats{
		UserID:        userID,
		ProjectCount:  defaultProjectCount,
		ActivityToday: defaultActivityToday,
		Notifications: defaultNotifications,
	}
}

// IncrementField adjusts the counter selected by field. The result is
// clamped at zero — the domain refuses to represent negative counts.
// A StatIncremented event is recorded reflecting the post-clamp delta.
func (s *UserStats) IncrementField(field StatField, delta int) {
	var counter *int
	switch field {
	case StatFieldProjects:
		counter = &s.ProjectCount
	case StatFieldActivity:
		counter = &s.ActivityToday
	case StatFieldNotifications:
		counter = &s.Notifications
	default:
		return // unknown field — no-op
	}
	before := *counter
	*counter += delta
	if *counter < 0 {
		*counter = 0
	}
	after := *counter
	s.Record(StatIncremented{
		UserID:   s.UserID,
		Field:    field,
		Delta:    after - before,
		NewValue: after,
	})
}
