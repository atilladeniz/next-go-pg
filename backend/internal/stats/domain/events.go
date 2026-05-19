package domain

import shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"

// StatIncremented is recorded when one of a user's stat counters
// changes (via UserStats.IncrementField). NewValue is the post-clamp
// value persisted, so subscribers don't have to re-derive it.
type StatIncremented struct {
	UserID   shared.UserID
	Field    StatField
	Delta    int
	NewValue int
}

func (StatIncremented) EventName() string { return "stats.incremented" }
