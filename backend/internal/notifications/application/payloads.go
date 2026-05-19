package application

// Call-specific payloads for each email type. The sender enriches them
// with shared template data (AppURL, SettingsURL) at render time, so
// callers stay focused on what's actually unique per call.

type MagicLinkPayload struct {
	URL string
}

type VerificationPayload struct {
	URL string
}

type TwoFactorOTPPayload struct {
	UserName string
	OTP      string
}

type TwoFactorEnabledPayload struct {
	UserName   string
	MethodName string
}

type PasskeyAddedPayload struct {
	UserName    string
	PasskeyName string
	Device      string
	Time        string
}

type LoginNotificationPayload struct {
	UserName  string
	Device    string
	IPAddress string
	Time      string
}
