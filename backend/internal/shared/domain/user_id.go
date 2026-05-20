// Package domain holds the SHARED KERNEL of the bounded contexts in
// this backend — the minimal vocabulary that every context needs to
// reference users without depending on another context's domain.
//
// Keep it small. Anything specific to a single bounded context belongs
// in that context's own domain/ package, not here.
package domain

import (
	"errors"
	"strings"
)

// UserID identifies a user across all bounded contexts. Constructed
// via NewUserID so the empty-string invariant is enforced at the
// boundary.
type UserID string

// NewUserID returns a validated UserID. Empty input is rejected.
func NewUserID(s string) (UserID, error) {
	if strings.TrimSpace(s) == "" {
		return "", errors.New("user id must not be empty")
	}
	return UserID(s), nil
}

func (u UserID) String() string { return string(u) }
