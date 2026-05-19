// Package domain is the auth bounded context's model. It carries the
// projection of the user record managed by an external auth provider
// (Better Auth at the time of writing). The auth context does NOT own
// the user table; it reads from it.
package domain

import shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"

// User is the pure-domain projection of a user, regardless of where
// the row physically lives. Authentication and persistence are
// external concerns — this type carries the fields the application
// reads.
type User struct {
	ID    shared.UserID
	Email string
	Name  string
}
