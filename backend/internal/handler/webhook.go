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
	"strings"
	"time"

	"github.com/atilladeniz/next-go-pg/backend/internal/application"
	"github.com/atilladeniz/next-go-pg/backend/internal/domain"
	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
	"github.com/mileusna/useragent"
)

// WebhookHandler handles all webhook endpoints for email notifications.
// It owns no transport state — rendering + SMTP live behind the
// application.EmailSender port; async dispatch lives behind
// application.JobEnqueuer.
type WebhookHandler struct {
	users       application.UserDirectory
	emails      application.EmailSender
	jobEnqueuer application.JobEnqueuer
}

// NewWebhookHandler creates a webhook handler. users and emails may be
// nil — in degraded modes the relevant endpoints fail fast with 503.
func NewWebhookHandler(users application.UserDirectory, emails application.EmailSender) *WebhookHandler {
	return &WebhookHandler{users: users, emails: emails}
}

// WithJobEnqueuer sets the job enqueuer for background email processing.
func (h *WebhookHandler) WithJobEnqueuer(enqueuer application.JobEnqueuer) *WebhookHandler {
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

// --- Endpoints ---

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

	if h.isKnownDevice(r.Context(), req) {
		respondJSON(w, MessageResponse{Message: "known device, notification skipped"})
		return
	}

	user, err := h.getUserByID(r.Context(), req.UserID)
	if err != nil {
		logger.Warn().Err(err).Str("user_id", req.UserID).Msg("User not found for session notification")
		respondJSON(w, MessageResponse{Message: "user not found"})
		return
	}

	if h.emails == nil {
		respondError(w, http.StatusServiceUnavailable, "email service unavailable")
		return
	}

	if err := h.emails.SendLoginNotification(r.Context(), user.Email, application.LoginNotificationPayload{
		UserName:  coalesce(user.Name, "Nutzer"),
		Device:    parseUserAgent(req.UserAgent),
		IPAddress: coalesce(req.IPAddress, "Unbekannt"),
		Time:      formatTimeGerman(time.Now()),
	}); err != nil {
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

	if h.jobEnqueuer != nil {
		if err := h.jobEnqueuer.EnqueueMagicLink(context.Background(), req.Email, req.URL); err != nil {
			logger.Error().Err(err).Str("email", req.Email).Msg("Failed to enqueue magic link job")
			respondError(w, http.StatusInternalServerError, "failed to enqueue email job")
			return
		}
		logger.Info().Str("email", req.Email).Msg("Magic link email enqueued")
		respondJSON(w, MessageResponse{Message: "magic link enqueued"})
		return
	}

	if h.emails == nil {
		respondError(w, http.StatusServiceUnavailable, "email service unavailable")
		return
	}

	if err := h.emails.SendMagicLink(r.Context(), req.Email, application.MagicLinkPayload{URL: req.URL}); err != nil {
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

	if h.jobEnqueuer != nil {
		if err := h.jobEnqueuer.EnqueueVerificationEmail(context.Background(), req.Email, req.Name, req.URL); err != nil {
			logger.Error().Err(err).Str("email", req.Email).Msg("Failed to enqueue verification email job")
			respondError(w, http.StatusInternalServerError, "failed to enqueue email job")
			return
		}
		logger.Info().Str("email", req.Email).Msg("Verification email enqueued")
		respondJSON(w, MessageResponse{Message: "verification email enqueued"})
		return
	}

	if h.emails == nil {
		respondError(w, http.StatusServiceUnavailable, "email service unavailable")
		return
	}

	if err := h.emails.SendVerification(r.Context(), req.Email, application.VerificationPayload{URL: req.URL}); err != nil {
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

	if h.jobEnqueuer != nil {
		if err := h.jobEnqueuer.Enqueue2FAOTP(context.Background(), req.Email, req.Name, req.OTP); err != nil {
			logger.Error().Err(err).Str("email", req.Email).Msg("Failed to enqueue 2FA OTP job")
			respondError(w, http.StatusInternalServerError, "failed to enqueue email job")
			return
		}
		logger.Info().Str("email", req.Email).Msg("2FA OTP email enqueued")
		respondJSON(w, MessageResponse{Message: "2fa otp enqueued"})
		return
	}

	if h.emails == nil {
		respondError(w, http.StatusServiceUnavailable, "email service unavailable")
		return
	}

	if err := h.emails.Send2FAOTP(r.Context(), req.Email, application.TwoFactorOTPPayload{
		UserName: coalesce(req.Name, "Nutzer"),
		OTP:      req.OTP,
	}); err != nil {
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

	if h.emails == nil {
		respondError(w, http.StatusServiceUnavailable, "email service unavailable")
		return
	}

	if err := h.emails.SendTwoFactorEnabled(r.Context(), req.Email, application.TwoFactorEnabledPayload{
		UserName:   coalesce(req.Name, "Nutzer"),
		MethodName: mapMethodToGerman(req.Method),
	}); err != nil {
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

	if h.emails == nil {
		respondError(w, http.StatusServiceUnavailable, "email service unavailable")
		return
	}

	if err := h.emails.SendPasskeyAdded(r.Context(), req.Email, application.PasskeyAddedPayload{
		UserName:    coalesce(req.Name, "Nutzer"),
		PasskeyName: coalesce(req.PasskeyName, "Unbenannt"),
		Device:      coalesce(req.Device, "Unbekanntes Gerät"),
		Time:        formatTimeGerman(time.Now()),
	}); err != nil {
		logger.Error().Err(err).Str("email", req.Email).Msg("Failed to send passkey added notification")
		respondError(w, http.StatusInternalServerError, "failed to send email")
		return
	}

	logger.Info().Str("email", req.Email).Str("passkey_name", req.PasskeyName).Msg("Passkey added notification sent")
	respondJSON(w, MessageResponse{Message: "notification sent"})
}

// --- Helpers ---

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

func (h *WebhookHandler) isKnownDevice(ctx context.Context, req SessionCreatedRequest) bool {
	if h.users == nil {
		return false
	}
	uid, err := domain.NewUserID(req.UserID)
	if err != nil {
		return false
	}
	known, err := h.users.HasKnownDevice(ctx, uid, req.UserAgent, req.IPAddress, req.SessionID)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to check for existing sessions")
		return false
	}
	return known
}

func (h *WebhookHandler) getUserByID(ctx context.Context, userID string) (*domain.User, error) {
	if h.users == nil {
		return &domain.User{}, nil
	}
	uid, err := domain.NewUserID(userID)
	if err != nil {
		return nil, err
	}
	return h.users.UserByID(ctx, uid)
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

func decodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func respondJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	respondJSON(w, ErrorResponse{Error: message})
}

func coalesce(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func parseUserAgent(uaString string) string {
	if uaString == "" {
		return "Unbekanntes Gerät"
	}

	// Special-case common CLI / API tools the useragent library treats
	// as plain strings.
	lower := strings.ToLower(uaString)
	switch {
	case strings.HasPrefix(lower, "curl/"):
		return "Kommandozeile (curl)"
	case strings.HasPrefix(lower, "wget/"):
		return "Kommandozeile (wget)"
	case strings.HasPrefix(lower, "postmanruntime/"):
		return "API-Client (Postman)"
	case strings.HasPrefix(lower, "insomnia/"):
		return "API-Client (Insomnia)"
	}

	ua := useragent.Parse(uaString)
	if ua.Name == "" {
		return "Unbekanntes Gerät"
	}
	parts := []string{ua.Name}
	if ua.OS != "" {
		parts = append(parts, "auf", ua.OS)
	}
	return strings.Join(parts, " ")
}

func formatTimeGerman(t time.Time) string {
	return t.Format("02.01.2006 um 15:04 Uhr")
}

func mapMethodToGerman(method string) string {
	switch method {
	case "totp":
		return "Authenticator-App (TOTP)"
	case "email-otp":
		return "E-Mail-OTP"
	case "passkey":
		return "Passkey"
	case "backup-codes":
		return "Backup-Codes"
	default:
		return "Zwei-Faktor-Authentifizierung"
	}
}
