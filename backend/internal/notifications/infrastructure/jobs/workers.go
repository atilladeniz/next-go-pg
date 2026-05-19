package jobs

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"

	notifapp "github.com/atilladeniz/next-go-pg/backend/internal/notifications/application"
	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
)

// Email workers are thin shells: extract args, call the right
// EmailSender method, log the outcome. No rendering, no SMTP — all of
// that lives in infrastructure/email/.

type SendMagicLinkWorker struct {
	river.WorkerDefaults[SendMagicLinkArgs]
	sender notifapp.EmailSender
}

func NewSendMagicLinkWorker(sender notifapp.EmailSender) *SendMagicLinkWorker {
	return &SendMagicLinkWorker{sender: sender}
}

func (w *SendMagicLinkWorker) Work(ctx context.Context, job *river.Job[SendMagicLinkArgs]) error {
	args := job.Args
	if err := w.sender.SendMagicLink(ctx, args.Email, notifapp.MagicLinkPayload{URL: args.URL}); err != nil {
		logger.Error().Err(err).Str("email", args.Email).Msg("Failed to send magic link email")
		return fmt.Errorf("send email: %w", err)
	}
	logger.Info().Str("email", args.Email).Msg("Magic link email sent via background job")
	return nil
}

type SendVerificationEmailWorker struct {
	river.WorkerDefaults[SendVerificationEmailArgs]
	sender notifapp.EmailSender
}

func NewSendVerificationEmailWorker(sender notifapp.EmailSender) *SendVerificationEmailWorker {
	return &SendVerificationEmailWorker{sender: sender}
}

func (w *SendVerificationEmailWorker) Work(ctx context.Context, job *river.Job[SendVerificationEmailArgs]) error {
	args := job.Args
	if err := w.sender.SendVerification(ctx, args.Email, notifapp.VerificationPayload{URL: args.URL}); err != nil {
		logger.Error().Err(err).Str("email", args.Email).Msg("Failed to send verification email")
		return fmt.Errorf("send email: %w", err)
	}
	logger.Info().Str("email", args.Email).Msg("Verification email sent via background job")
	return nil
}

type Send2FAOTPWorker struct {
	river.WorkerDefaults[Send2FAOTPArgs]
	sender notifapp.EmailSender
}

func NewSend2FAOTPWorker(sender notifapp.EmailSender) *Send2FAOTPWorker {
	return &Send2FAOTPWorker{sender: sender}
}

func (w *Send2FAOTPWorker) Work(ctx context.Context, job *river.Job[Send2FAOTPArgs]) error {
	args := job.Args
	name := args.Name
	if name == "" {
		name = "Nutzer"
	}
	if err := w.sender.Send2FAOTP(ctx, args.Email, notifapp.TwoFactorOTPPayload{UserName: name, OTP: args.OTP}); err != nil {
		logger.Error().Err(err).Str("email", args.Email).Msg("Failed to send 2FA OTP email")
		return fmt.Errorf("send email: %w", err)
	}
	logger.Info().Str("email", args.Email).Msg("2FA OTP email sent via background job")
	return nil
}

type SendLoginNotificationWorker struct {
	river.WorkerDefaults[SendLoginNotificationArgs]
	sender notifapp.EmailSender
}

func NewSendLoginNotificationWorker(sender notifapp.EmailSender) *SendLoginNotificationWorker {
	return &SendLoginNotificationWorker{sender: sender}
}

func (w *SendLoginNotificationWorker) Work(ctx context.Context, job *river.Job[SendLoginNotificationArgs]) error {
	args := job.Args
	name := args.UserName
	if name == "" {
		name = "Nutzer"
	}
	if err := w.sender.SendLoginNotification(ctx, args.Email, notifapp.LoginNotificationPayload{
		UserName:  name,
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

// Register hooks this context's workers into a River workers registry.
func Register(workers *river.Workers, sender notifapp.EmailSender) {
	river.AddWorker(workers, NewSendMagicLinkWorker(sender))
	river.AddWorker(workers, NewSendVerificationEmailWorker(sender))
	river.AddWorker(workers, NewSend2FAOTPWorker(sender))
	river.AddWorker(workers, NewSendLoginNotificationWorker(sender))
}
