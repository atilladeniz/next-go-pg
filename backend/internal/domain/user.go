package domain

import (
	"errors"
	"strings"
)

// UserID identifies a user. Constructed via NewUserID so the empty-string
// invariant is enforced at the boundary.
type UserID string

// NewUserID returns a validated UserID. Empty input is rejected.
func NewUserID(s string) (UserID, error) {
	if strings.TrimSpace(s) == "" {
		return "", errors.New("user id must not be empty")
	}
	return UserID(s), nil
}

func (u UserID) String() string { return string(u) }

// User is the pure-domain projection of a user, regardless of where the
// row physically lives (Better Auth's `user` table at the time of
// writing). Authentication and persistence are external concerns — this
// type carries the fields the application reads.
type User struct {
	ID    UserID
	Email string
	Name  string
}
