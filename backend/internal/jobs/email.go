// Package jobs defines background job workers for the River queue.
package jobs

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"
	"gopkg.in/gomail.v2"

	"github.com/atilladeniz/next-go-pg/backend/internal/templates"
	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
)

// EmailConfig holds email configuration for job workers.
type EmailConfig struct {
	SMTPHost    string
	SMTPPort    int
	SMTPFrom    string
	AppURL      string
	SettingsURL string
}

// --- Magic Link Email Job ---

// SendMagicLinkArgs defines the arguments for sending a magic link email.
type SendMagicLinkArgs struct {
	Email string `json:"email"`
	URL   string `json:"url"`
}

func (SendMagicLinkArgs) Kind() string { return "send_magic_link" }

// SendMagicLinkWorker processes magic link email jobs.
type SendMagicLinkWorker struct {
	river.WorkerDefaults[SendMagicLinkArgs]
	mailer *gomail.Dialer
	config *EmailConfig
}

func NewSendMagicLinkWorker(config *EmailConfig) *SendMagicLinkWorker {
	dialer := gomail.NewDialer(config.SMTPHost, config.SMTPPort, "", "")
	dialer.SSL = false
	return &SendMagicLinkWorker{
		mailer: dialer,
		config: config,
	}
}

func (w *SendMagicLinkWorker) Work(ctx context.Context, job *river.Job[SendMagicLinkArgs]) error {
	args := job.Args

	body, err := templates.RenderMagicLink(templates.MagicLinkData{
		EmailData:    templates.EmailData{AppURL: w.config.AppURL},
		MagicLinkURL: args.URL,
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to render magic link template")
		return fmt.Errorf("render template: %w", err)
	}

	if err := w.sendEmail(args.Email, "Dein Anmelde-Link", body); err != nil {
		logger.Error().Err(err).Str("email", args.Email).Msg("Failed to send magic link email")
		return fmt.Errorf("send email: %w", err)
	}

	logger.Info().Str("email", args.Email).Msg("Magic link email sent via background job")
	return nil
}

func (w *SendMagicLinkWorker) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", w.config.SMTPFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return w.mailer.DialAndSend(m)
}

// --- Verification Email Job ---

// SendVerificationEmailArgs defines the arguments for sending a verification email.
type SendVerificationEmailArgs struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

func (SendVerificationEmailArgs) Kind() string { return "send_verification_email" }

// SendVerificationEmailWorker processes verification email jobs.
type SendVerificationEmailWorker struct {
	river.WorkerDefaults[SendVerificationEmailArgs]
	mailer *gomail.Dialer
	config *EmailConfig
}

func NewSendVerificationEmailWorker(config *EmailConfig) *SendVerificationEmailWorker {
	dialer := gomail.NewDialer(config.SMTPHost, config.SMTPPort, "", "")
	dialer.SSL = false
	return &SendVerificationEmailWorker{
		mailer: dialer,
		config: config,
	}
}

func (w *SendVerificationEmailWorker) Work(ctx context.Context, job *river.Job[SendVerificationEmailArgs]) error {
	args := job.Args

	body, err := templates.RenderVerification(templates.VerificationData{
		EmailData: templates.EmailData{AppURL: w.config.AppURL},
		VerifyURL: args.URL,
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to render verification template")
		return fmt.Errorf("render template: %w", err)
	}

	if err := w.sendEmail(args.Email, "E-Mail bestätigen", body); err != nil {
		logger.Error().Err(err).Str("email", args.Email).Msg("Failed to send verification email")
		return fmt.Errorf("send email: %w", err)
	}

	logger.Info().Str("email", args.Email).Msg("Verification email sent via background job")
	return nil
}

func (w *SendVerificationEmailWorker) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", w.config.SMTPFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return w.mailer.DialAndSend(m)
}

// --- 2FA OTP Email Job ---

// Send2FAOTPArgs defines the arguments for sending a 2FA OTP email.
type Send2FAOTPArgs struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	OTP   string `json:"otp"`
}

func (Send2FAOTPArgs) Kind() string { return "send_2fa_otp" }

// Send2FAOTPWorker processes 2FA OTP email jobs.
type Send2FAOTPWorker struct {
	river.WorkerDefaults[Send2FAOTPArgs]
	mailer *gomail.Dialer
	config *EmailConfig
}

func NewSend2FAOTPWorker(config *EmailConfig) *Send2FAOTPWorker {
	dialer := gomail.NewDialer(config.SMTPHost, config.SMTPPort, "", "")
	dialer.SSL = false
	return &Send2FAOTPWorker{
		mailer: dialer,
		config: config,
	}
}

func (w *Send2FAOTPWorker) Work(ctx context.Context, job *river.Job[Send2FAOTPArgs]) error {
	args := job.Args

	userName := args.Name
	if userName == "" {
		userName = "Nutzer"
	}

	body, err := templates.RenderTwoFactorOTP(templates.TwoFactorOTPData{
		EmailData: templates.EmailData{AppURL: w.config.AppURL},
		UserName:  userName,
		OTP:       args.OTP,
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to render 2FA OTP template")
		return fmt.Errorf("render template: %w", err)
	}

	if err := w.sendEmail(args.Email, "Dein Sicherheitscode", body); err != nil {
		logger.Error().Err(err).Str("email", args.Email).Msg("Failed to send 2FA OTP email")
		return fmt.Errorf("send email: %w", err)
	}

	logger.Info().Str("email", args.Email).Msg("2FA OTP email sent via background job")
	return nil
}

func (w *Send2FAOTPWorker) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", w.config.SMTPFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return w.mailer.DialAndSend(m)
}

// --- Login Notification Email Job ---

// SendLoginNotificationArgs defines the arguments for sending a login notification email.
type SendLoginNotificationArgs struct {
	Email     string `json:"email"`
	UserName  string `json:"userName"`
	Device    string `json:"device"`
	IPAddress string `json:"ipAddress"`
}

func (SendLoginNotificationArgs) Kind() string { return "send_login_notification" }

// SendLoginNotificationWorker processes login notification email jobs.
type SendLoginNotificationWorker struct {
	river.WorkerDefaults[SendLoginNotificationArgs]
	mailer *gomail.Dialer
	config *EmailConfig
}

func NewSendLoginNotificationWorker(config *EmailConfig) *SendLoginNotificationWorker {
	dialer := gomail.NewDialer(config.SMTPHost, config.SMTPPort, "", "")
	dialer.SSL = false
	return &SendLoginNotificationWorker{
		mailer: dialer,
		config: config,
	}
}

func (w *SendLoginNotificationWorker) Work(ctx context.Context, job *river.Job[SendLoginNotificationArgs]) error {
	args := job.Args

	userName := args.UserName
	if userName == "" {
		userName = "Nutzer"
	}

	ipAddress := args.IPAddress
	if ipAddress == "" {
		ipAddress = "Unbekannt"
	}

	body, err := templates.RenderLoginNotification(templates.LoginNotificationData{
		EmailData: templates.EmailData{
			AppURL:      w.config.AppURL,
			SettingsURL: w.config.SettingsURL,
		},
		UserName:  userName,
		Device:    args.Device,
		IPAddress: ipAddress,
		Time:      "", // Will be set by template with current time
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to render login notification template")
		return fmt.Errorf("render template: %w", err)
	}

	if err := w.sendEmail(args.Email, "Neue Anmeldung von neuem Gerät", body); err != nil {
		logger.Error().Err(err).Str("email", args.Email).Msg("Failed to send login notification email")
		return fmt.Errorf("send email: %w", err)
	}

	logger.Info().Str("email", args.Email).Msg("Login notification email sent via background job")
	return nil
}

func (w *SendLoginNotificationWorker) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", w.config.SMTPFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return w.mailer.DialAndSend(m)
}
