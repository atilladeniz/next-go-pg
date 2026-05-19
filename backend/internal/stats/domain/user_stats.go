// Package domain is the stats bounded context's aggregate model.
// Pure Go, no persistence, no HTTP, no other context's domain.
package domain

import (
	"time"

	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
)

// Default seed values for a freshly created UserStats aggregate. Kept
// here (in the domain) rather than in the repository so the invariant
// "a new user starts with these counts" is owned by the model.
const (
	defaultProjectCount  = 3
	defaultActivityToday = 10
	defaultNotifications = 2
)

// UserStats is the aggregate root for a user's statistics. It carries
// the counter invariants and records domain events on mutation.
type UserStats struct {
	shared.AggregateBase

	ID            uint
	UserID        shared.UserID
	ProjectCount  int
	ActivityToday int
	Notifications int
	LastLogin     time.Time
	MemberSince   time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Compile-time aggregate check.
var _ shared.AggregateRoot = (*UserStats)(nil)

// NewUserStats is the aggregate factory for a freshly-seeded row.
// Persistence adapters call this when GetOrCreate sees no existing row.
func NewUserStats(userID shared.UserID) *UserStats {
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
		return
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
