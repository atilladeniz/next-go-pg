package templates

import (
	"bytes"
	"html/template"
)

// EmailData contains common data for email templates
type EmailData struct {
	AppURL      string
	SettingsURL string
}

// MagicLinkData contains data for magic link emails
type MagicLinkData struct {
	EmailData
	MagicLinkURL string
}

// VerificationData contains data for verification emails
type VerificationData struct {
	EmailData
	VerifyURL string
}

// LoginNotificationData contains data for login notification emails
type LoginNotificationData struct {
	EmailData
	UserName  string
	Device    string
	IPAddress string
	Time      string
}

// TwoFactorOTPData contains data for 2FA OTP emails
type TwoFactorOTPData struct {
	EmailData
	UserName string
	OTP      string
}

// TwoFactorEnabledData contains data for 2FA enabled notification emails
type TwoFactorEnabledData struct {
	EmailData
	UserName   string
	MethodName string
}

// PasskeyAddedData contains data for passkey added notification emails
type PasskeyAddedData struct {
	EmailData
	UserName    string
	PasskeyName string
	Device      string
	Time        string
}

// Template definitions
var (
	MagicLinkTemplate = template.Must(template.New("magic_link").Parse(`
<h1>Anmeldung</h1>
<p>Klicke auf den Button, um dich anzumelden:</p>
<p style="margin: 24px 0;">
	<a href="{{.MagicLinkURL}}" style="display: inline-block; padding: 12px 24px; background-color: #000; color: #fff; text-decoration: none; border-radius: 6px;">Jetzt anmelden</a>
</p>
<p style="font-size: 14px; color: #666;">
	Oder kopiere diesen Link: <a href="{{.MagicLinkURL}}">{{.MagicLinkURL}}</a>
</p>
<p style="font-size: 14px; color: #666;">Der Link ist 10 Minuten gültig.</p>
`))

	VerificationTemplate = template.Must(template.New("verification").Parse(`
<h1>Willkommen!</h1>
<p>Klicke auf den folgenden Link, um deine E-Mail-Adresse zu bestätigen:</p>
<p style="margin: 24px 0;">
	<a href="{{.VerifyURL}}" style="display: inline-block; padding: 12px 24px; background-color: #000; color: #fff; text-decoration: none; border-radius: 6px;">E-Mail bestätigen</a>
</p>
<p style="font-size: 14px; color: #666;">
	Oder kopiere diesen Link: <a href="{{.VerifyURL}}">{{.VerifyURL}}</a>
</p>
<p style="font-size: 14px; color: #666;">Der Link ist 24 Stunden gültig.</p>
`))

	LoginNotificationTemplate = template.Must(template.New("login_notification").Parse(`
<h1>Neue Anmeldung in deinem Konto</h1>
<p>Hallo {{.UserName}},</p>
<p>Wir haben eine Anmeldung von einem neuen Gerät oder Standort festgestellt:</p>
<ul>
	<li><strong>Gerät:</strong> {{.Device}}</li>
	<li><strong>IP-Adresse:</strong> {{.IPAddress}}</li>
	<li><strong>Zeit:</strong> {{.Time}}</li>
</ul>
<p>Wenn du das nicht warst, überprüfe bitte sofort deine aktiven Sessions und beende verdächtige Sitzungen:</p>
<p><a href="{{.SettingsURL}}" style="display: inline-block; padding: 12px 24px; background-color: #000; color: #fff; text-decoration: none; border-radius: 6px;">Sessions verwalten</a></p>
<p style="margin-top: 16px; font-size: 14px; color: #666;">
	Oder kopiere diesen Link: <a href="{{.SettingsURL}}">{{.SettingsURL}}</a>
</p>
`))

	TwoFactorOTPTemplate = template.Must(template.New("2fa_otp").Parse(`
<h1>Dein Sicherheitscode</h1>
<p>Hallo {{.UserName}},</p>
<p>Dein Einmal-Code für die Zwei-Faktor-Authentifizierung lautet:</p>
<p style="margin: 24px 0; text-align: center;">
	<span style="display: inline-block; padding: 16px 32px; background-color: #f4f4f4; font-size: 32px; font-weight: bold; letter-spacing: 8px; font-family: monospace; border-radius: 8px;">{{.OTP}}</span>
</p>
<p style="font-size: 14px; color: #666;">Dieser Code ist 3 Minuten gültig.</p>
<p style="font-size: 14px; color: #666;">Falls du diesen Code nicht angefordert hast, ignoriere diese E-Mail.</p>
`))

	TwoFactorEnabledTemplate = template.Must(template.New("2fa_enabled").Parse(`
<h1>Sicherheit aktiviert</h1>
<p>Hallo {{.UserName}},</p>
<p>Die folgende Sicherheitsmethode wurde für dein Konto aktiviert:</p>
<p style="margin: 16px 0; padding: 12px 16px; background-color: #f0f9ff; border-left: 4px solid #0284c7; font-weight: bold;">{{.MethodName}}</p>
<p style="font-size: 14px; color: #666;">Falls du das nicht warst, überprüfe sofort deine Sicherheitseinstellungen:</p>
<p style="margin: 16px 0;">
	<a href="{{.SettingsURL}}" style="display: inline-block; padding: 12px 24px; background-color: #000; color: #fff; text-decoration: none; border-radius: 6px;">Einstellungen öffnen</a>
</p>
`))

	PasskeyAddedTemplate = template.Must(template.New("passkey_added").Parse(`
<h1>Neuer Passkey hinzugefügt</h1>
<p>Hallo {{.UserName}},</p>
<p>Ein neuer Passkey wurde zu deinem Konto hinzugefügt:</p>
<ul style="margin: 16px 0; padding: 16px; background-color: #f4f4f4; border-radius: 8px; list-style: none;">
	<li><strong>Name:</strong> {{.PasskeyName}}</li>
	<li><strong>Gerät:</strong> {{.Device}}</li>
	<li><strong>Hinzugefügt:</strong> {{.Time}}</li>
</ul>
<p style="font-size: 14px; color: #666;">Falls du das nicht warst, entferne den Passkey sofort in deinen Einstellungen:</p>
<p style="margin: 16px 0;">
	<a href="{{.SettingsURL}}" style="display: inline-block; padding: 12px 24px; background-color: #dc2626; color: #fff; text-decoration: none; border-radius: 6px;">Passkey entfernen</a>
</p>
`))
)

// RenderMagicLink renders the magic link email template
func RenderMagicLink(data MagicLinkData) (string, error) {
	return render(MagicLinkTemplate, data)
}

// RenderVerification renders the verification email template
func RenderVerification(data VerificationData) (string, error) {
	return render(VerificationTemplate, data)
}

// RenderLoginNotification renders the login notification email template
func RenderLoginNotification(data LoginNotificationData) (string, error) {
	return render(LoginNotificationTemplate, data)
}

// RenderTwoFactorOTP renders the 2FA OTP email template
func RenderTwoFactorOTP(data TwoFactorOTPData) (string, error) {
	return render(TwoFactorOTPTemplate, data)
}

// RenderTwoFactorEnabled renders the 2FA enabled notification email template
func RenderTwoFactorEnabled(data TwoFactorEnabledData) (string, error) {
	return render(TwoFactorEnabledTemplate, data)
}

// RenderPasskeyAdded renders the passkey added notification email template
func RenderPasskeyAdded(data PasskeyAddedData) (string, error) {
	return render(PasskeyAddedTemplate, data)
}

func render(tmpl *template.Template, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
