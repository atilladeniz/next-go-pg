// Package jobs defines background job workers for the River queue.
package jobs

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"

	"github.com/atilladeniz/next-go-pg/backend/internal/application"
	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
)

// Email workers delegate rendering + SMTP to application.EmailSender;
// the workers' only job is to extract job args and call the right
// sender method, then log the outcome.

// --- Magic Link ---

type SendMagicLinkArgs struct {
	Email string `json:"email"`
	URL   string `json:"url"`
}

func (SendMagicLinkArgs) Kind() string { return "send_magic_link" }

type SendMagicLinkWorker struct {
	river.WorkerDefaults[SendMagicLinkArgs]
	sender application.EmailSender
}

func NewSendMagicLinkWorker(sender application.EmailSender) *SendMagicLinkWorker {
	return &SendMagicLinkWorker{sender: sender}
}

func (w *SendMagicLinkWorker) Work(ctx context.Context, job *river.Job[SendMagicLinkArgs]) error {
	args := job.Args
	if err := w.sender.SendMagicLink(ctx, args.Email, application.MagicLinkPayload{URL: args.URL}); err != nil {
		logger.Error().Err(err).Str("email", args.Email).Msg("Failed to send magic link email")
		return fmt.Errorf("send email: %w", err)
	}
	logger.Info().Str("email", args.Email).Msg("Magic link email sent via background job")
	return nil
}

// --- Verification ---

type SendVerificationEmailArgs struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

func (SendVerificationEmailArgs) Kind() string { return "send_verification_email" }

type SendVerificationEmailWorker struct {
	river.WorkerDefaults[SendVerificationEmailArgs]
	sender application.EmailSender
}

func NewSendVerificationEmailWorker(sender application.EmailSender) *SendVerificationEmailWorker {
	return &SendVerificationEmailWorker{sender: sender}
}

func (w *SendVerificationEmailWorker) Work(ctx context.Context, job *river.Job[SendVerificationEmailArgs]) error {
	args := job.Args
	if err := w.sender.SendVerification(ctx, args.Email, application.VerificationPayload{URL: args.URL}); err != nil {
		logger.Error().Err(err).Str("email", args.Email).Msg("Failed to send verification email")
		return fmt.Errorf("send email: %w", err)
	}
	logger.Info().Str("email", args.Email).Msg("Verification email sent via background job")
	return nil
}

// --- 2FA OTP ---

type Send2FAOTPArgs struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	OTP   string `json:"otp"`
}

func (Send2FAOTPArgs) Kind() string { return "send_2fa_otp" }

type Send2FAOTPWorker struct {
	river.WorkerDefaults[Send2FAOTPArgs]
	sender application.EmailSender
}

func NewSend2FAOTPWorker(sender application.EmailSender) *Send2FAOTPWorker {
	return &Send2FAOTPWorker{sender: sender}
}

func (w *Send2FAOTPWorker) Work(ctx context.Context, job *river.Job[Send2FAOTPArgs]) error {
	args := job.Args
	userName := args.Name
	if userName == "" {
		userName = "Nutzer"
	}
	if err := w.sender.Send2FAOTP(ctx, args.Email, application.TwoFactorOTPPayload{
		UserName: userName,
		OTP:      args.OTP,
	}); err != nil {
		logger.Error().Err(err).Str("email", args.Email).Msg("Failed to send 2FA OTP email")
		return fmt.Errorf("send email: %w", err)
	}
	logger.Info().Str("email", args.Email).Msg("2FA OTP email sent via background job")
	return nil
}

// --- Login Notification ---

type SendLoginNotificationArgs struct {
	Email     string `json:"email"`
	UserName  string `json:"userName"`
	Device    string `json:"device"`
	IPAddress string `json:"ipAddress"`
	Time      string `json:"time"`
}

func (SendLoginNotificationArgs) Kind() string { return "send_login_notification" }

type SendLoginNotificationWorker struct {
	river.WorkerDefaults[SendLoginNotificationArgs]
	sender application.EmailSender
}

func NewSendLoginNotificationWorker(sender application.EmailSender) *SendLoginNotificationWorker {
	return &SendLoginNotificationWorker{sender: sender}
}

func (w *SendLoginNotificationWorker) Work(ctx context.Context, job *river.Job[SendLoginNotificationArgs]) error {
	args := job.Args
	userName := args.UserName
	if userName == "" {
		userName = "Nutzer"
	}
	if err := w.sender.SendLoginNotification(ctx, args.Email, application.LoginNotificationPayload{
		UserName:  userName,
		Device:    args.Device,
		IPAddress: args.IPAddress,
		Time:      args.Time,
	}); err != nil {
		logger.Error().Err(err).Str("email", args.Email).Msg("Failed to send login notification email")
		return fmt.Errorf("send email: %w", err)
	}
	logger.Info().Str("email", args.Email).Msg("Login notification email sent via background job")
	return nil
}
