// Package application is the notifications bounded context's use-case
// layer. It defines three ports: EmailSender (synchronous send),
// JobEnqueuer (asynchronous send via the queue), and UserDirectory
// (look up the minimum user info needed to personalise an email).
package application

import (
	"context"

	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
)

// EmailSender renders and dispatches transactional emails. The
// concrete implementation lives under internal/notifications/infrastructure/email/.
type EmailSender interface {
	SendMagicLink(ctx context.Context, to string, payload MagicLinkPayload) error
	SendVerification(ctx context.Context, to string, payload VerificationPayload) error
	Send2FAOTP(ctx context.Context, to string, payload TwoFactorOTPPayload) error
	SendTwoFactorEnabled(ctx context.Context, to string, payload TwoFactorEnabledPayload) error
	SendPasskeyAdded(ctx context.Context, to string, payload PasskeyAddedPayload) error
	SendLoginNotification(ctx context.Context, to string, payload LoginNotificationPayload) error
}

// JobEnqueuer schedules notification emails for asynchronous delivery.
// Each method maps to a typed River job in the queue adapter.
type JobEnqueuer interface {
	EnqueueMagicLink(ctx context.Context, email, url string) error
	EnqueueVerificationEmail(ctx context.Context, email, name, url string) error
	Enqueue2FAOTP(ctx context.Context, email, name, otp string) error
	EnqueueLoginNotification(ctx context.Context, email, userName, device, ipAddress string) error
}

// UserDirectory is the notifications context's port for the user-info
// it needs in order to send personalised emails. The auth context owns
// the actual user records; an anti-corruption adapter in the
// composition root translates between auth's UserDirectory and this
// notifications-local port so the two contexts stay decoupled.
type UserDirectory interface {
	UserByID(ctx context.Context, userID shared.UserID) (UserSnapshot, error)
	HasKnownDevice(ctx context.Context, userID shared.UserID, userAgent, ipAddress, excludeSessionID string) (bool, error)
}

// UserSnapshot is the minimal projection notifications consumes. It
// intentionally carries only what an email template might render —
// nothing about credentials, sessions, or identity provenance.
type UserSnapshot struct {
	Email string
	Name  string
}
