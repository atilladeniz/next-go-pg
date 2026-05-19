// Package email is the concrete adapter for application.EmailSender.
// It owns the template rendering and the SMTP transport. Templates are
// currently inline raw strings; issue #59 tracks moving them into
// proper .html files via embed.FS.
package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"gopkg.in/gomail.v2"

	"github.com/atilladeniz/next-go-pg/backend/internal/application"
)

// Config configures the SMTP transport and the base URLs used inside
// templates (footer links, settings URL, ...).
type Config struct {
	SMTPHost string
	SMTPPort int
	SMTPFrom string
	AppURL   string
}

// Sender is the gomail-backed implementation of application.EmailSender.
type Sender struct {
	dialer      *gomail.Dialer
	from        string
	appURL      string
	settingsURL string
}

var _ application.EmailSender = (*Sender)(nil)

// NewSender builds a Sender from SMTP config. The dialer is held but
// not connected — DialAndSend opens a new connection per message.
func NewSender(cfg Config) *Sender {
	dialer := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, "", "")
	dialer.SSL = false
	appURL := cfg.AppURL
	if appURL == "" {
		appURL = "http://localhost:3000"
	}
	return &Sender{
		dialer:      dialer,
		from:        cfg.SMTPFrom,
		appURL:      appURL,
		settingsURL: appURL + "/settings",
	}
}

// Template data structs — payload + sender-supplied URLs. Kept private
// so the template surface stays an implementation detail.

type magicLinkData struct {
	URL    string
	AppURL string
}

type verificationData struct {
	URL    string
	AppURL string
}

type twoFactorOTPData struct {
	UserName string
	OTP      string
	AppURL   string
}

type twoFactorEnabledData struct {
	UserName    string
	MethodName  string
	AppURL      string
	SettingsURL string
}

type passkeyAddedData struct {
	UserName    string
	PasskeyName string
	Device      string
	Time        string
	AppURL      string
	SettingsURL string
}

type loginNotificationData struct {
	UserName    string
	Device      string
	IPAddress   string
	Time        string
	AppURL      string
	SettingsURL string
}

func (s *Sender) SendMagicLink(_ context.Context, to string, p application.MagicLinkPayload) error {
	body, err := render(magicLinkTmpl, magicLinkData{URL: p.URL, AppURL: s.appURL})
	if err != nil {
		return err
	}
	return s.send(to, "Dein Anmelde-Link", body)
}

func (s *Sender) SendVerification(_ context.Context, to string, p application.VerificationPayload) error {
	body, err := render(verificationTmpl, verificationData{URL: p.URL, AppURL: s.appURL})
	if err != nil {
		return err
	}
	return s.send(to, "E-Mail bestätigen", body)
}

func (s *Sender) Send2FAOTP(_ context.Context, to string, p application.TwoFactorOTPPayload) error {
	body, err := render(twoFactorOTPTmpl, twoFactorOTPData{
		UserName: p.UserName,
		OTP:      p.OTP,
		AppURL:   s.appURL,
	})
	if err != nil {
		return err
	}
	return s.send(to, "Dein Sicherheitscode", body)
}

func (s *Sender) SendTwoFactorEnabled(_ context.Context, to string, p application.TwoFactorEnabledPayload) error {
	body, err := render(twoFactorEnabledTmpl, twoFactorEnabledData{
		UserName:    p.UserName,
		MethodName:  p.MethodName,
		AppURL:      s.appURL,
		SettingsURL: s.settingsURL,
	})
	if err != nil {
		return err
	}
	return s.send(to, "Sicherheitsmethode aktiviert", body)
}

func (s *Sender) SendPasskeyAdded(_ context.Context, to string, p application.PasskeyAddedPayload) error {
	body, err := render(passkeyAddedTmpl, passkeyAddedData{
		UserName:    p.UserName,
		PasskeyName: p.PasskeyName,
		Device:      p.Device,
		Time:        p.Time,
		AppURL:      s.appURL,
		SettingsURL: s.settingsURL,
	})
	if err != nil {
		return err
	}
	return s.send(to, "Neuer Passkey hinzugefügt", body)
}

func (s *Sender) SendLoginNotification(_ context.Context, to string, p application.LoginNotificationPayload) error {
	body, err := render(loginNotificationTmpl, loginNotificationData{
		UserName:    p.UserName,
		Device:      p.Device,
		IPAddress:   p.IPAddress,
		Time:        p.Time,
		AppURL:      s.appURL,
		SettingsURL: s.settingsURL,
	})
	if err != nil {
		return err
	}
	return s.send(to, "Neue Anmeldung von neuem Gerät", body)
}

func (s *Sender) send(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("smtp send: %w", err)
	}
	return nil
}

func render(tmpl *template.Template, data any) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
