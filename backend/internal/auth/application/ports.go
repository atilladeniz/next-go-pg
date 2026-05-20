// Package application is the auth bounded context's use-case layer.
// Today it exposes only a directory-read port; identity mutation is
// owned by Better Auth itself.
package application

import (
	"context"

	auth "github.com/atilladeniz/next-go-pg/backend/internal/auth/domain"
	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
)

// UserDirectory reads user records owned by the external auth
// provider. Adapters live under internal/auth/infrastructure/.
type UserDirectory interface {
	UserByID(ctx context.Context, userID shared.UserID) (*auth.User, error)
	HasKnownDevice(ctx context.Context, userID shared.UserID, userAgent, ipAddress, excludeSessionID string) (bool, error)
}
