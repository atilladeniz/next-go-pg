package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
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
	if h.db != nil {
		var count int64
		h.db.Table("session").
			Where(`"userId" = ? AND "userAgent" = ? AND "ipAddress" = ? AND id != ?`,
				req.UserID, req.UserAgent, req.IPAddress, req.SessionID).
			Limit(1).
			Count(&count)

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

func parseUserAgent(ua string) string {
	if ua == "" {
		return "Unbekanntes Gerät"
	}

	browser := "Browser"
	os := "System"

	// Detect browser
	switch {
	case contains(ua, "Firefox"):
		browser = "Firefox"
	case contains(ua, "Edg/"):
		browser = "Edge"
	case contains(ua, "Chrome"):
		browser = "Chrome"
	case contains(ua, "Safari"):
		browser = "Safari"
	}

	// Detect OS
	switch {
	case contains(ua, "Windows"):
		os = "Windows"
	case contains(ua, "Mac OS"):
		os = "macOS"
	case contains(ua, "Linux"):
		os = "Linux"
	case contains(ua, "Android"):
		os = "Android"
	case contains(ua, "iPhone"), contains(ua, "iPad"):
		os = "iOS"
	}

	return browser + " auf " + os
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsImpl(s, substr))
}

func containsImpl(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func verifyWebhookSecret(provided, expected string) bool {
	if provided == "" || expected == "" {
		return false
	}
	// Use constant-time comparison to prevent timing attacks
	return hmac.Equal([]byte(provided), []byte(expected))
}

// ComputeHMAC computes HMAC-SHA256 for webhook signature verification
func ComputeHMAC(message, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
