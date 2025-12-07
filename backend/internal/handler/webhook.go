package handler

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/atilladeniz/next-go-pg/backend/internal/jobs"
	"github.com/atilladeniz/next-go-pg/backend/internal/templates"
	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
	"github.com/mileusna/useragent"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// WebhookHandler handles all webhook endpoints for email notifications
type WebhookHandler struct {
	db          *gorm.DB
	mailer      *gomail.Dialer
	config      *webhookConfig
	jobEnqueuer jobs.JobEnqueuer // Optional: if set, emails are sent via background jobs
}

type webhookConfig struct {
	smtpFrom    string
	appURL      string
	settingsURL string
}

// NewWebhookHandler creates a new webhook handler with mailer configuration
func NewWebhookHandler(db *gorm.DB) *WebhookHandler {
	smtpHost := getEnvOrDefault("SMTP_HOST", "127.0.0.1")
	smtpPort := getEnvAsIntOrDefault("SMTP_PORT", 1025)
	appURL := getEnvOrDefault("NEXT_PUBLIC_APP_URL", "http://localhost:3000")

	dialer := gomail.NewDialer(smtpHost, smtpPort, "", "")
	dialer.SSL = false

	return &WebhookHandler{
		db:     db,
		mailer: dialer,
		config: &webhookConfig{
			smtpFrom:    getEnvOrDefault("SMTP_FROM", "noreply@localhost"),
			appURL:      appURL,
			settingsURL: appURL + "/settings",
		},
	}
}

// WithJobEnqueuer sets the job enqueuer for background email processing
func (h *WebhookHandler) WithJobEnqueuer(enqueuer jobs.JobEnqueuer) *WebhookHandler {
	h.jobEnqueuer = enqueuer
	return h
}

// --- Request Types ---

type SessionCreatedRequest struct {
	SessionID string `json:"sessionId"`
	UserID    string `json:"userId"`
	UserAgent string `json:"userAgent"`
	IPAddress string `json:"ipAddress"`
}

type SendMagicLinkRequest struct {
	Email string `json:"email"`
	URL   string `json:"url"`
}

type SendVerificationEmailRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type Send2FAOTPRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	OTP   string `json:"otp"`
}

type Send2FAEnabledNotificationRequest struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	Method string `json:"method"`
}

type SendPasskeyAddedNotificationRequest struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	PasskeyName string `json:"passkeyName"`
	Device      string `json:"device"`
}

// --- Handlers ---

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
	if !h.verifySecret(w, r) {
		return
	}

	var req SessionCreatedRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Check if device is known (with locking to prevent race conditions)
	if h.isKnownDevice(req) {
		respondJSON(w, MessageResponse{Message: "known device, notification skipped"})
		return
	}

	// Get user info
	user, err := h.getUserByID(req.UserID)
	if err != nil {
		logger.Warn().Err(err).Str("user_id", req.UserID).Msg("User not found for session notification")
		respondJSON(w, MessageResponse{Message: "user not found"})
		return
	}

	// Render and send email
	body, err := templates.RenderLoginNotification(templates.LoginNotificationData{
		EmailData: templates.EmailData{
			AppURL:      h.config.appURL,
			SettingsURL: h.config.settingsURL,
		},
		UserName:  coalesce(user.Name, "Nutzer"),
		Device:    parseUserAgent(req.UserAgent),
		IPAddress: coalesce(req.IPAddress, "Unbekannt"),
		Time:      formatTimeGerman(time.Now()),
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to render login notification template")
		respondError(w, http.StatusInternalServerError, "template error")
		return
	}

	if err := h.sendEmail(user.Email, "Neue Anmeldung von neuem Gerät", body); err != nil {
		logger.Error().Err(err).Str("email", user.Email).Msg("Failed to send login notification email")
		respondError(w, http.StatusInternalServerError, "failed to send email")
		return
	}

	logger.Info().Str("user_id", req.UserID).Str("email", user.Email).Msg("Login notification email sent")
	respondJSON(w, MessageResponse{Message: "notification sent"})
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
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Use background job if available
	if h.jobEnqueuer != nil {
		if err := jobs.EnqueueMagicLink(context.Background(), h.jobEnqueuer, req.Email, req.URL); err != nil {
			logger.Error().Err(err).Str("email", req.Email).Msg("Failed to enqueue magic link job")
			respondError(w, http.StatusInternalServerError, "failed to enqueue email job")
			return
		}
		logger.Info().Str("email", req.Email).Msg("Magic link email enqueued")
		respondJSON(w, MessageResponse{Message: "magic link enqueued"})
		return
	}

	// Fallback to synchronous sending
	body, err := templates.RenderMagicLink(templates.MagicLinkData{
		EmailData:    templates.EmailData{AppURL: h.config.appURL},
		MagicLinkURL: req.URL,
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to render magic link template")
		respondError(w, http.StatusInternalServerError, "template error")
		return
	}

	if err := h.sendEmail(req.Email, "Dein Anmelde-Link", body); err != nil {
		logger.Error().Err(err).Str("email", req.Email).Msg("Failed to send magic link email")
		respondError(w, http.StatusInternalServerError, "failed to send email")
		return
	}

	logger.Info().Str("email", req.Email).Msg("Magic link email sent")
	respondJSON(w, MessageResponse{Message: "magic link sent"})
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
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Use background job if available
	if h.jobEnqueuer != nil {
		if err := jobs.EnqueueVerificationEmail(context.Background(), h.jobEnqueuer, req.Email, req.Name, req.URL); err != nil {
			logger.Error().Err(err).Str("email", req.Email).Msg("Failed to enqueue verification email job")
			respondError(w, http.StatusInternalServerError, "failed to enqueue email job")
			return
		}
		logger.Info().Str("email", req.Email).Msg("Verification email enqueued")
		respondJSON(w, MessageResponse{Message: "verification email enqueued"})
		return
	}

	// Fallback to synchronous sending
	body, err := templates.RenderVerification(templates.VerificationData{
		EmailData: templates.EmailData{AppURL: h.config.appURL},
		VerifyURL: req.URL,
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to render verification template")
		respondError(w, http.StatusInternalServerError, "template error")
		return
	}

	if err := h.sendEmail(req.Email, "E-Mail bestätigen", body); err != nil {
		logger.Error().Err(err).Str("email", req.Email).Msg("Failed to send verification email")
		respondError(w, http.StatusInternalServerError, "failed to send email")
		return
	}

	logger.Info().Str("email", req.Email).Msg("Verification email sent")
	respondJSON(w, MessageResponse{Message: "verification email sent"})
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
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Use background job if available
	if h.jobEnqueuer != nil {
		if err := jobs.Enqueue2FAOTP(context.Background(), h.jobEnqueuer, req.Email, req.Name, req.OTP); err != nil {
			logger.Error().Err(err).Str("email", req.Email).Msg("Failed to enqueue 2FA OTP job")
			respondError(w, http.StatusInternalServerError, "failed to enqueue email job")
			return
		}
		logger.Info().Str("email", req.Email).Msg("2FA OTP email enqueued")
		respondJSON(w, MessageResponse{Message: "2fa otp enqueued"})
		return
	}

	// Fallback to synchronous sending
	body, err := templates.RenderTwoFactorOTP(templates.TwoFactorOTPData{
		EmailData: templates.EmailData{AppURL: h.config.appURL},
		UserName:  coalesce(req.Name, "Nutzer"),
		OTP:       req.OTP,
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to render 2FA OTP template")
		respondError(w, http.StatusInternalServerError, "template error")
		return
	}

	if err := h.sendEmail(req.Email, "Dein Sicherheitscode", body); err != nil {
		logger.Error().Err(err).Str("email", req.Email).Msg("Failed to send 2FA OTP email")
		respondError(w, http.StatusInternalServerError, "failed to send email")
		return
	}

	logger.Info().Str("email", req.Email).Msg("2FA OTP email sent")
	respondJSON(w, MessageResponse{Message: "2fa otp sent"})
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
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	body, err := templates.RenderTwoFactorEnabled(templates.TwoFactorEnabledData{
		EmailData:  templates.EmailData{AppURL: h.config.appURL, SettingsURL: h.config.settingsURL},
		UserName:   coalesce(req.Name, "Nutzer"),
		MethodName: mapMethodToGerman(req.Method),
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to render 2FA enabled template")
		respondError(w, http.StatusInternalServerError, "template error")
		return
	}

	if err := h.sendEmail(req.Email, "Sicherheitsmethode aktiviert", body); err != nil {
		logger.Error().Err(err).Str("email", req.Email).Msg("Failed to send 2FA enabled notification")
		respondError(w, http.StatusInternalServerError, "failed to send email")
		return
	}

	logger.Info().Str("email", req.Email).Str("method", req.Method).Msg("2FA enabled notification sent")
	respondJSON(w, MessageResponse{Message: "notification sent"})
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
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	body, err := templates.RenderPasskeyAdded(templates.PasskeyAddedData{
		EmailData:   templates.EmailData{AppURL: h.config.appURL, SettingsURL: h.config.settingsURL},
		UserName:    coalesce(req.Name, "Nutzer"),
		PasskeyName: coalesce(req.PasskeyName, "Unbenannt"),
		Device:      coalesce(req.Device, "Unbekanntes Gerät"),
		Time:        formatTimeGerman(time.Now()),
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to render passkey added template")
		respondError(w, http.StatusInternalServerError, "template error")
		return
	}

	if err := h.sendEmail(req.Email, "Neuer Passkey hinzugefügt", body); err != nil {
		logger.Error().Err(err).Str("email", req.Email).Msg("Failed to send passkey added notification")
		respondError(w, http.StatusInternalServerError, "failed to send email")
		return
	}

	logger.Info().Str("email", req.Email).Str("passkey_name", req.PasskeyName).Msg("Passkey added notification sent")
	respondJSON(w, MessageResponse{Message: "notification sent"})
}

// --- Helper Methods ---

func (h *WebhookHandler) verifySecret(w http.ResponseWriter, r *http.Request) bool {
	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	if webhookSecret == "" {
		return true // No secret configured, allow all
	}

	providedSecret := r.Header.Get("X-Webhook-Secret")
	if !verifyWebhookSecret(providedSecret, webhookSecret) {
		respondError(w, http.StatusUnauthorized, "invalid webhook secret")
		return false
	}
	return true
}

func (h *WebhookHandler) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", h.config.smtpFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return h.mailer.DialAndSend(m)
}

func (h *WebhookHandler) isKnownDevice(req SessionCreatedRequest) bool {
	if h.db == nil {
		return false
	}

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
		return false
	}
	return count > 0
}

type userInfo struct {
	Email string `gorm:"column:email"`
	Name  string `gorm:"column:name"`
}

func (h *WebhookHandler) getUserByID(userID string) (*userInfo, error) {
	if h.db == nil {
		return &userInfo{}, nil
	}

	var user userInfo
	if err := h.db.Table("user").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// --- Utility Functions ---

func verifyWebhookSecret(provided, expected string) bool {
	if provided == "" || expected == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(provided), []byte(expected)) == 1
}

// ComputeHMAC computes HMAC-SHA256 for webhook signature verification
func ComputeHMAC(message, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func parseUserAgent(uaString string) string {
	if uaString == "" {
		return "Unbekanntes Gerät"
	}

	ua := useragent.Parse(uaString)

	// Handle bots
	if ua.Bot {
		return "Automatisierter Zugriff (Bot)"
	}

	// Handle CLI tools
	lowerUA := strings.ToLower(uaString)
	switch {
	case strings.Contains(lowerUA, "curl"):
		return "Kommandozeile (curl)"
	case strings.Contains(lowerUA, "wget"):
		return "Kommandozeile (wget)"
	case strings.Contains(lowerUA, "httpie"):
		return "Kommandozeile (HTTPie)"
	case strings.Contains(lowerUA, "postman"):
		return "API-Client (Postman)"
	case strings.Contains(lowerUA, "insomnia"):
		return "API-Client (Insomnia)"
	}

	browser := coalesce(ua.Name, "Browser")
	osName := coalesce(ua.OS, "System")

	// Normalize OS names
	switch {
	case strings.Contains(osName, "Mac"):
		osName = "macOS"
	case strings.Contains(osName, "iOS"):
		osName = "iOS"
	}

	return browser + " auf " + osName
}

func mapMethodToGerman(method string) string {
	switch method {
	case "passkey":
		return "Passkey"
	case "totp":
		return "Authenticator-App (TOTP)"
	default:
		return "Zwei-Faktor-Authentifizierung"
	}
}

func formatTimeGerman(t time.Time) string {
	return t.Format("02.01.2006, 15:04")
}

func coalesce(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func getEnvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsIntOrDefault(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func decodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func respondJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
