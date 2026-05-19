// Package application is the notifications bounded context's use-case
// layer. It defines two ports: EmailSender (synchronous send) and
// JobEnqueuer (asynchronous send via the queue).
package application

import "context"

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
