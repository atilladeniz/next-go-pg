// Package jobs is the notifications context's queue-side adapter:
// River job args, workers that consume application.EmailSender, and
// the enqueuer that implements application.JobEnqueuer.
package jobs

type SendMagicLinkArgs struct {
	Email string `json:"email"`
	URL   string `json:"url"`
}

func (SendMagicLinkArgs) Kind() string { return "send_magic_link" }

type SendVerificationEmailArgs struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

func (SendVerificationEmailArgs) Kind() string { return "send_verification_email" }

type Send2FAOTPArgs struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	OTP   string `json:"otp"`
}

func (Send2FAOTPArgs) Kind() string { return "send_2fa_otp" }

type SendLoginNotificationArgs struct {
	Email     string `json:"email"`
	UserName  string `json:"userName"`
	Device    string `json:"device"`
	IPAddress string `json:"ipAddress"`
	Time      string `json:"time"`
}

func (SendLoginNotificationArgs) Kind() string { return "send_login_notification" }
