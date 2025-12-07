package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
	"github.com/mileusna/useragent"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WebhookHandler struct {
	db     *gorm.DB
	mailer *gomail.Dialer
}

func NewWebhookHandler(db *gorm.DB) *WebhookHandler {
	// Setup mailer
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "127.0.0.1"
	}
	smtpPort := 1025 // Default Mailpit port

	dialer := gomail.NewDialer(smtpHost, smtpPort, "", "")
	dialer.SSL = false

	return &WebhookHandler{
		db:     db,
		mailer: dialer,
	}
}

// SessionCreatedRequest is the payload from Better Auth session hook
type SessionCreatedRequest struct {
	SessionID string `json:"sessionId"`
	UserID    string `json:"userId"`
	UserAgent string `json:"userAgent"`
	IPAddress string `json:"ipAddress"`
}

// SessionCreated godoc
// @Summary Handle new session webhook
// @Description Called by Better Auth when a new session is created. Sends notification email for new devices.
// @Tags webhooks
// @Accept json
// @Produce json
// @Param X-Webhook-Secret header string true "Webhook secret for authentication"
// @Param request body SessionCreatedRequest true "Session data"
// @Success 200 {object} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /webhooks/session-created [post]
func (h *WebhookHandler) SessionCreated(w http.ResponseWriter, r *http.Request) {
	// Verify webhook secret
	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	if webhookSecret != "" {
		providedSecret := r.Header.Get("X-Webhook-Secret")
		if !verifyWebhookSecret(providedSecret, webhookSecret) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid webhook secret"})
			return
		}
	}

	var req SessionCreatedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	// Check if this device/IP combination is known
	// Use transaction with FOR UPDATE to prevent race conditions when multiple
	// sessions are created simultaneously from the same device
	if h.db != nil {
		var count int64
		err := h.db.Transaction(func(tx *gorm.DB) error {
			return tx.Table("session").
				Where(`"userId" = ? AND "userAgent" = ? AND "ipAddress" = ? AND id != ?`,
					req.UserID, req.UserAgent, req.IPAddress, req.SessionID).
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Limit(1).
				Count(&count).Error
		})

		if err != nil {
			logger.Warn().Err(err).Msg("Failed to check for existing sessions")
		}

		if count > 0 {
			// Known device, skip notification
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(MessageResponse{Message: "known device, notification skipped"})
			return
		}
	}

	// Get user info
	var user struct {
		Email string `gorm:"column:email"`
		Name  string `gorm:"column:name"`
	}

	if h.db != nil {
		if err := h.db.Table("user").Where("id = ?", req.UserID).First(&user).Error; err != nil {
			logger.Warn().Err(err).Str("user_id", req.UserID).Msg("User not found for session notification")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(MessageResponse{Message: "user not found"})
			return
		}
	}

	// Parse device info
	device := parseUserAgent(req.UserAgent)

	// Format time in German
	timeStr := time.Now().Format("02.01.2006, 15:04")

	// Build email
	appURL := os.Getenv("NEXT_PUBLIC_APP_URL")
	if appURL == "" {
		appURL = "http://localhost:3000"
	}
	settingsURL := appURL + "/settings"

	ipAddress := req.IPAddress
	if ipAddress == "" {
		ipAddress = "Unbekannt"
	}

	emailBody := `
		<h1>Neue Anmeldung in deinem Konto</h1>
		<p>Hallo ` + user.Name + `,</p>
		<p>Wir haben eine Anmeldung von einem neuen Gerät oder Standort festgestellt:</p>
		<ul>
			<li><strong>Gerät:</strong> ` + device + `</li>
			<li><strong>IP-Adresse:</strong> ` + ipAddress + `</li>
			<li><strong>Zeit:</strong> ` + timeStr + `</li>
		</ul>
		<p>Wenn du das nicht warst, überprüfe bitte sofort deine aktiven Sessions und beende verdächtige Sitzungen:</p>
		<p><a href="` + settingsURL + `" style="display: inline-block; padding: 12px 24px; background-color: #000; color: #fff; text-decoration: none; border-radius: 6px;">Sessions verwalten</a></p>
		<p style="margin-top: 16px; font-size: 14px; color: #666;">
			Oder kopiere diesen Link: <a href="` + settingsURL + `">` + settingsURL + `</a>
		</p>
	`

	// Send email
	smtpFrom := os.Getenv("SMTP_FROM")
	if smtpFrom == "" {
		smtpFrom = "noreply@localhost"
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtpFrom)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Neue Anmeldung von neuem Gerät")
	m.SetBody("text/html", emailBody)

	if err := h.mailer.DialAndSend(m); err != nil {
		logger.Error().Err(err).Str("email", user.Email).Msg("Failed to send login notification email")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to send email"})
		return
	}

	logger.Info().
		Str("user_id", req.UserID).
		Str("email", user.Email).
		Str("device", device).
		Msg("Login notification email sent")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MessageResponse{Message: "notification sent"})
}

func parseUserAgent(uaString string) string {
	if uaString == "" {
		return "Unbekanntes Gerät"
	}

	// Use battle-tested library for user agent parsing
	ua := useragent.Parse(uaString)

	browser := ua.Name
	if browser == "" {
		browser = "Browser"
	}

	osName := ua.OS
	if osName == "" {
		osName = "System"
	}

	// Normalize OS names for German display
	switch {
	case strings.Contains(osName, "Mac"):
		osName = "macOS"
	case strings.Contains(osName, "iOS"):
		osName = "iOS"
	}

	return browser + " auf " + osName
}

func verifyWebhookSecret(provided, expected string) bool {
	if provided == "" || expected == "" {
		return false
	}
	// Use constant-time comparison to prevent timing attacks
	// subtle.ConstantTimeCompare is the correct function for comparing secrets
	return subtle.ConstantTimeCompare([]byte(provided), []byte(expected)) == 1
}

// ComputeHMAC computes HMAC-SHA256 for webhook signature verification
func ComputeHMAC(message, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// SendMagicLinkRequest is the payload for magic link email
type SendMagicLinkRequest struct {
	Email string `json:"email"`
	URL   string `json:"url"`
}

// SendMagicLink godoc
// @Summary Send magic link email
// @Description Called by Better Auth to send magic link login email
// @Tags webhooks
// @Accept json
// @Produce json
// @Param X-Webhook-Secret header string true "Webhook secret for authentication"
// @Param request body SendMagicLinkRequest true "Magic link data"
// @Success 200 {object} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /webhooks/send-magic-link [post]
func (h *WebhookHandler) SendMagicLink(w http.ResponseWriter, r *http.Request) {
	if !h.verifySecret(w, r) {
		return
	}

	var req SendMagicLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	appURL := os.Getenv("NEXT_PUBLIC_APP_URL")
	if appURL == "" {
		appURL = "http://localhost:3000"
	}

	emailBody := `
		<h1>Anmeldung</h1>
		<p>Klicke auf den Button, um dich anzumelden:</p>
		<p style="margin: 24px 0;">
			<a href="` + req.URL + `" style="display: inline-block; padding: 12px 24px; background-color: #000; color: #fff; text-decoration: none; border-radius: 6px;">Jetzt anmelden</a>
		</p>
		<p style="font-size: 14px; color: #666;">
			Oder kopiere diesen Link: <a href="` + req.URL + `">` + req.URL + `</a>
		</p>
		<p style="font-size: 14px; color: #666;">Der Link ist 10 Minuten gültig.</p>
	`

	if err := h.sendEmail(req.Email, "Dein Anmelde-Link", emailBody); err != nil {
		logger.Error().Err(err).Str("email", req.Email).Msg("Failed to send magic link email")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to send email"})
		return
	}

	logger.Info().Str("email", req.Email).Msg("Magic link email sent")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MessageResponse{Message: "magic link sent"})
}

// SendVerificationEmailRequest is the payload for verification email
type SendVerificationEmailRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

// SendVerificationEmail godoc
// @Summary Send verification email
// @Description Called by Better Auth to send email verification link
// @Tags webhooks
// @Accept json
// @Produce json
// @Param X-Webhook-Secret header string true "Webhook secret for authentication"
// @Param request body SendVerificationEmailRequest true "Verification data"
// @Success 200 {object} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /webhooks/send-verification-email [post]
func (h *WebhookHandler) SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	if !h.verifySecret(w, r) {
		return
	}

	var req SendVerificationEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	emailBody := `
		<h1>Willkommen!</h1>
		<p>Klicke auf den folgenden Link, um deine E-Mail-Adresse zu bestätigen:</p>
		<p style="margin: 24px 0;">
			<a href="` + req.URL + `" style="display: inline-block; padding: 12px 24px; background-color: #000; color: #fff; text-decoration: none; border-radius: 6px;">E-Mail bestätigen</a>
		</p>
		<p style="font-size: 14px; color: #666;">
			Oder kopiere diesen Link: <a href="` + req.URL + `">` + req.URL + `</a>
		</p>
		<p style="font-size: 14px; color: #666;">Der Link ist 24 Stunden gültig.</p>
	`

	if err := h.sendEmail(req.Email, "E-Mail bestätigen", emailBody); err != nil {
		logger.Error().Err(err).Str("email", req.Email).Msg("Failed to send verification email")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to send email"})
		return
	}

	logger.Info().Str("email", req.Email).Msg("Verification email sent")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MessageResponse{Message: "verification email sent"})
}

// Helper to verify webhook secret
func (h *WebhookHandler) verifySecret(w http.ResponseWriter, r *http.Request) bool {
	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	if webhookSecret != "" {
		providedSecret := r.Header.Get("X-Webhook-Secret")
		if !verifyWebhookSecret(providedSecret, webhookSecret) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid webhook secret"})
			return false
		}
	}
	return true
}

// Helper to send email
func (h *WebhookHandler) sendEmail(to, subject, body string) error {
	smtpFrom := os.Getenv("SMTP_FROM")
	if smtpFrom == "" {
		smtpFrom = "noreply@localhost"
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtpFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return h.mailer.DialAndSend(m)
}

// Send2FAOTPRequest is the payload for 2FA OTP email
type Send2FAOTPRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	OTP   string `json:"otp"`
}

// Send2FAOTP godoc
// @Summary Send 2FA OTP email
// @Description Called by Better Auth to send 2FA one-time password via email
// @Tags webhooks
// @Accept json
// @Produce json
// @Param X-Webhook-Secret header string true "Webhook secret for authentication"
// @Param request body Send2FAOTPRequest true "OTP data"
// @Success 200 {object} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /webhooks/send-2fa-otp [post]
func (h *WebhookHandler) Send2FAOTP(w http.ResponseWriter, r *http.Request) {
	if !h.verifySecret(w, r) {
		return
	}

	var req Send2FAOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	name := req.Name
	if name == "" {
		name = "Nutzer"
	}

	emailBody := `
		<h1>Dein Sicherheitscode</h1>
		<p>Hallo ` + name + `,</p>
		<p>Dein Einmal-Code für die Zwei-Faktor-Authentifizierung lautet:</p>
		<p style="margin: 24px 0; text-align: center;">
			<span style="display: inline-block; padding: 16px 32px; background-color: #f4f4f4; font-size: 32px; font-weight: bold; letter-spacing: 8px; font-family: monospace; border-radius: 8px;">` + req.OTP + `</span>
		</p>
		<p style="font-size: 14px; color: #666;">Dieser Code ist 3 Minuten gültig.</p>
		<p style="font-size: 14px; color: #666;">Falls du diesen Code nicht angefordert hast, ignoriere diese E-Mail.</p>
	`

	if err := h.sendEmail(req.Email, "Dein Sicherheitscode", emailBody); err != nil {
		logger.Error().Err(err).Str("email", req.Email).Msg("Failed to send 2FA OTP email")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to send email"})
		return
	}

	logger.Info().Str("email", req.Email).Msg("2FA OTP email sent")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MessageResponse{Message: "2fa otp sent"})
}

// Send2FAEnabledNotificationRequest is the payload for 2FA enabled notification
type Send2FAEnabledNotificationRequest struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	Method string `json:"method"` // "totp" or "passkey"
}

// Send2FAEnabledNotification godoc
// @Summary Send 2FA enabled notification
// @Description Notifies user that 2FA has been enabled on their account
// @Tags webhooks
// @Accept json
// @Produce json
// @Param X-Webhook-Secret header string true "Webhook secret for authentication"
// @Param request body Send2FAEnabledNotificationRequest true "Notification data"
// @Success 200 {object} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /webhooks/send-2fa-enabled [post]
func (h *WebhookHandler) Send2FAEnabledNotification(w http.ResponseWriter, r *http.Request) {
	if !h.verifySecret(w, r) {
		return
	}

	var req Send2FAEnabledNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	name := req.Name
	if name == "" {
		name = "Nutzer"
	}

	methodName := "Zwei-Faktor-Authentifizierung"
	if req.Method == "passkey" {
		methodName = "Passkey"
	} else if req.Method == "totp" {
		methodName = "Authenticator-App (TOTP)"
	}

	appURL := os.Getenv("NEXT_PUBLIC_APP_URL")
	if appURL == "" {
		appURL = "http://localhost:3000"
	}
	settingsURL := appURL + "/settings"

	emailBody := `
		<h1>Sicherheit aktiviert</h1>
		<p>Hallo ` + name + `,</p>
		<p>Die folgende Sicherheitsmethode wurde für dein Konto aktiviert:</p>
		<p style="margin: 16px 0; padding: 12px 16px; background-color: #f0f9ff; border-left: 4px solid #0284c7; font-weight: bold;">` + methodName + `</p>
		<p style="font-size: 14px; color: #666;">Falls du das nicht warst, überprüfe sofort deine Sicherheitseinstellungen:</p>
		<p style="margin: 16px 0;">
			<a href="` + settingsURL + `" style="display: inline-block; padding: 12px 24px; background-color: #000; color: #fff; text-decoration: none; border-radius: 6px;">Einstellungen öffnen</a>
		</p>
	`

	if err := h.sendEmail(req.Email, "Sicherheitsmethode aktiviert", emailBody); err != nil {
		logger.Error().Err(err).Str("email", req.Email).Msg("Failed to send 2FA enabled notification")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to send email"})
		return
	}

	logger.Info().Str("email", req.Email).Str("method", req.Method).Msg("2FA enabled notification sent")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MessageResponse{Message: "notification sent"})
}

// SendPasskeyAddedNotificationRequest is the payload for passkey added notification
type SendPasskeyAddedNotificationRequest struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	PasskeyName string `json:"passkeyName"`
	Device      string `json:"device"`
}

// SendPasskeyAddedNotification godoc
// @Summary Send passkey added notification
// @Description Notifies user that a new passkey has been added to their account
// @Tags webhooks
// @Accept json
// @Produce json
// @Param X-Webhook-Secret header string true "Webhook secret for authentication"
// @Param request body SendPasskeyAddedNotificationRequest true "Notification data"
// @Success 200 {object} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /webhooks/send-passkey-added [post]
func (h *WebhookHandler) SendPasskeyAddedNotification(w http.ResponseWriter, r *http.Request) {
	if !h.verifySecret(w, r) {
		return
	}

	var req SendPasskeyAddedNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	name := req.Name
	if name == "" {
		name = "Nutzer"
	}

	passkeyName := req.PasskeyName
	if passkeyName == "" {
		passkeyName = "Unbenannt"
	}

	device := req.Device
	if device == "" {
		device = "Unbekanntes Gerät"
	}

	timeStr := time.Now().Format("02.01.2006, 15:04")

	appURL := os.Getenv("NEXT_PUBLIC_APP_URL")
	if appURL == "" {
		appURL = "http://localhost:3000"
	}
	settingsURL := appURL + "/settings"

	emailBody := `
		<h1>Neuer Passkey hinzugefügt</h1>
		<p>Hallo ` + name + `,</p>
		<p>Ein neuer Passkey wurde zu deinem Konto hinzugefügt:</p>
		<ul style="margin: 16px 0; padding: 16px; background-color: #f4f4f4; border-radius: 8px; list-style: none;">
			<li><strong>Name:</strong> ` + passkeyName + `</li>
			<li><strong>Gerät:</strong> ` + device + `</li>
			<li><strong>Hinzugefügt:</strong> ` + timeStr + `</li>
		</ul>
		<p style="font-size: 14px; color: #666;">Falls du das nicht warst, entferne den Passkey sofort in deinen Einstellungen:</p>
		<p style="margin: 16px 0;">
			<a href="` + settingsURL + `" style="display: inline-block; padding: 12px 24px; background-color: #dc2626; color: #fff; text-decoration: none; border-radius: 6px;">Passkey entfernen</a>
		</p>
	`

	if err := h.sendEmail(req.Email, "Neuer Passkey hinzugefügt", emailBody); err != nil {
		logger.Error().Err(err).Str("email", req.Email).Msg("Failed to send passkey added notification")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to send email"})
		return
	}

	logger.Info().Str("email", req.Email).Str("passkey_name", passkeyName).Msg("Passkey added notification sent")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MessageResponse{Message: "notification sent"})
}
