package email

import "html/template"

// Inline templates. Issue #59 tracks the move to embedded .html files
// for syntax highlighting + shared `_base.html` / `_components.html`.

var magicLinkTmpl = template.Must(template.New("magic_link").Parse(`
<h1>Anmeldung</h1>
<p>Klicke auf den Button, um dich anzumelden:</p>
<p style="margin: 24px 0;">
	<a href="{{.URL}}" style="display: inline-block; padding: 12px 24px; background-color: #000; color: #fff; text-decoration: none; border-radius: 6px;">Jetzt anmelden</a>
</p>
<p style="font-size: 14px; color: #666;">
	Oder kopiere diesen Link: <a href="{{.URL}}">{{.URL}}</a>
</p>
<p style="font-size: 14px; color: #666;">Der Link ist 10 Minuten gültig.</p>
`))

var verificationTmpl = template.Must(template.New("verification").Parse(`
<h1>Willkommen!</h1>
<p>Klicke auf den folgenden Link, um deine E-Mail-Adresse zu bestätigen:</p>
<p style="margin: 24px 0;">
	<a href="{{.URL}}" style="display: inline-block; padding: 12px 24px; background-color: #000; color: #fff; text-decoration: none; border-radius: 6px;">E-Mail bestätigen</a>
</p>
<p style="font-size: 14px; color: #666;">
	Oder kopiere diesen Link: <a href="{{.URL}}">{{.URL}}</a>
</p>
<p style="font-size: 14px; color: #666;">Der Link ist 24 Stunden gültig.</p>
`))

var twoFactorOTPTmpl = template.Must(template.New("2fa_otp").Parse(`
<h1>Dein Sicherheitscode</h1>
<p>Hallo {{.UserName}},</p>
<p>Dein Einmal-Code für die Zwei-Faktor-Authentifizierung lautet:</p>
<p style="margin: 24px 0; text-align: center;">
	<span style="display: inline-block; padding: 16px 32px; background-color: #f4f4f4; font-size: 32px; font-weight: bold; letter-spacing: 8px; font-family: monospace; border-radius: 8px;">{{.OTP}}</span>
</p>
<p style="font-size: 14px; color: #666;">Dieser Code ist 3 Minuten gültig.</p>
<p style="font-size: 14px; color: #666;">Falls du diesen Code nicht angefordert hast, ignoriere diese E-Mail.</p>
`))

var twoFactorEnabledTmpl = template.Must(template.New("2fa_enabled").Parse(`
<h1>Sicherheit aktiviert</h1>
<p>Hallo {{.UserName}},</p>
<p>Die folgende Sicherheitsmethode wurde für dein Konto aktiviert:</p>
<p style="margin: 16px 0; padding: 12px 16px; background-color: #f0f9ff; border-left: 4px solid #0284c7; font-weight: bold;">{{.MethodName}}</p>
<p style="font-size: 14px; color: #666;">Falls du das nicht warst, überprüfe sofort deine Sicherheitseinstellungen:</p>
<p style="margin: 16px 0;">
	<a href="{{.SettingsURL}}" style="display: inline-block; padding: 12px 24px; background-color: #000; color: #fff; text-decoration: none; border-radius: 6px;">Einstellungen öffnen</a>
</p>
`))

var passkeyAddedTmpl = template.Must(template.New("passkey_added").Parse(`
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

var loginNotificationTmpl = template.Must(template.New("login_notification").Parse(`
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
