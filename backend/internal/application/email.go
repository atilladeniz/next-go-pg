package application

import "context"

// MagicLinkPayload carries the call-specific data for the magic-link
// email. Templates also need the app base URL — the sender supplies
// that at render time so handlers don't have to thread it.
type MagicLinkPayload struct {
	URL string
}

// VerificationPayload — call-specific data for the email-verification email.
type VerificationPayload struct {
	URL string
}

// TwoFactorOTPPayload — call-specific data for the 2FA-OTP email.
type TwoFactorOTPPayload struct {
	UserName string
	OTP      string
}

// TwoFactorEnabledPayload — call-specific data for the 2FA-enabled email.
type TwoFactorEnabledPayload struct {
	UserName   string
	MethodName string
}

// PasskeyAddedPayload — call-specific data for the passkey-added email.
type PasskeyAddedPayload struct {
	UserName    string
	PasskeyName string
	Device      string
	Time        string
}

// LoginNotificationPayload — call-specific data for the new-device login email.
type LoginNotificationPayload struct {
	UserName  string
	Device    string
	IPAddress string
	Time      string
}

// EmailSender renders and dispatches transactional emails. Concrete
// implementations live under internal/infrastructure/email/ and own the
// template rendering plus the SMTP transport. The app base URL is the
// sender's responsibility, not the caller's.
type EmailSender interface {
	SendMagicLink(ctx context.Context, to string, payload MagicLinkPayload) error
	SendVerification(ctx context.Context, to string, payload VerificationPayload) error
	Send2FAOTP(ctx context.Context, to string, payload TwoFactorOTPPayload) error
	SendTwoFactorEnabled(ctx context.Context, to string, payload TwoFactorEnabledPayload) error
	SendPasskeyAdded(ctx context.Context, to string, payload PasskeyAddedPayload) error
	SendLoginNotification(ctx context.Context, to string, payload LoginNotificationPayload) error
}
