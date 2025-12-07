package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWebhookSecretVerification tests the webhook secret verification
func TestWebhookSecretVerification(t *testing.T) {
	tests := []struct {
		name           string
		webhookSecret  string
		providedSecret string
		expectedStatus int
	}{
		{
			name:           "valid secret",
			webhookSecret:  "test-secret-123",
			providedSecret: "test-secret-123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid secret",
			webhookSecret:  "test-secret-123",
			providedSecret: "wrong-secret",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "missing secret header",
			webhookSecret:  "test-secret-123",
			providedSecret: "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "no webhook secret configured - allows all",
			webhookSecret:  "",
			providedSecret: "",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment
			if tt.webhookSecret != "" {
				os.Setenv("WEBHOOK_SECRET", tt.webhookSecret)
				defer os.Unsetenv("WEBHOOK_SECRET")
			} else {
				os.Unsetenv("WEBHOOK_SECRET")
			}

			handler := NewWebhookHandler(nil)

			// Create request with magic link payload
			payload := SendMagicLinkRequest{
				Email: "test@example.com",
				URL:   "https://example.com/magic-link?token=abc123",
			}
			body, _ := json.Marshal(payload)
			req := httptest.NewRequest(http.MethodPost, "/webhooks/send-magic-link", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.providedSecret != "" {
				req.Header.Set("X-Webhook-Secret", tt.providedSecret)
			}

			rr := httptest.NewRecorder()
			handler.SendMagicLink(rr, req)

			// For unauthorized, we check status directly
			// For authorized, we expect either success or email error (SMTP not configured)
			if tt.expectedStatus == http.StatusUnauthorized {
				assert.Equal(t, http.StatusUnauthorized, rr.Code)

				var response ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "invalid webhook secret", response.Error)
			} else {
				// If secret is valid, we either get 200 or 500 (SMTP error)
				// The important thing is we're not getting 401
				assert.NotEqual(t, http.StatusUnauthorized, rr.Code)
			}
		})
	}
}

// TestVerifyWebhookSecret tests the constant-time secret comparison
func TestVerifyWebhookSecret(t *testing.T) {
	tests := []struct {
		name     string
		provided string
		expected string
		want     bool
	}{
		{
			name:     "matching secrets",
			provided: "my-secret-key",
			expected: "my-secret-key",
			want:     true,
		},
		{
			name:     "different secrets",
			provided: "wrong-secret",
			expected: "my-secret-key",
			want:     false,
		},
		{
			name:     "empty provided",
			provided: "",
			expected: "my-secret-key",
			want:     false,
		},
		{
			name:     "empty expected",
			provided: "my-secret-key",
			expected: "",
			want:     false,
		},
		{
			name:     "both empty",
			provided: "",
			expected: "",
			want:     false,
		},
		{
			name:     "similar secrets different length",
			provided: "secret",
			expected: "secrets",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := verifyWebhookSecret(tt.provided, tt.expected)
			assert.Equal(t, tt.want, result)
		})
	}
}

// TestSendMagicLinkRequestValidation tests request body parsing
func TestSendMagicLinkRequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "valid request",
			body:           `{"email":"test@example.com","url":"https://example.com/link"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid json",
			body:           `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty body",
			body:           ``,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("WEBHOOK_SECRET") // No secret required
			handler := NewWebhookHandler(nil)

			req := httptest.NewRequest(http.MethodPost, "/webhooks/send-magic-link", bytes.NewReader([]byte(tt.body)))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.SendMagicLink(rr, req)

			if tt.expectedStatus == http.StatusBadRequest {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			} else {
				// Valid request should not return 400
				assert.NotEqual(t, http.StatusBadRequest, rr.Code)
			}
		})
	}
}

// TestSessionCreatedRequestValidation tests session created webhook
func TestSessionCreatedRequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "valid request",
			body:           `{"sessionId":"sess-123","userId":"user-456","userAgent":"Mozilla/5.0","ipAddress":"192.168.1.1"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid json",
			body:           `not valid json`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("WEBHOOK_SECRET")
			handler := NewWebhookHandler(nil)

			req := httptest.NewRequest(http.MethodPost, "/webhooks/session-created", bytes.NewReader([]byte(tt.body)))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.SessionCreated(rr, req)

			if tt.expectedStatus == http.StatusBadRequest {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			} else {
				assert.NotEqual(t, http.StatusBadRequest, rr.Code)
			}
		})
	}
}

// TestSendVerificationEmailRequestValidation tests verification email webhook
func TestSendVerificationEmailRequestValidation(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "valid request",
			body:           `{"email":"test@example.com","name":"Test User","url":"https://example.com/verify"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid json",
			body:           `{broken`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("WEBHOOK_SECRET")
			handler := NewWebhookHandler(nil)

			req := httptest.NewRequest(http.MethodPost, "/webhooks/send-verification-email", bytes.NewReader([]byte(tt.body)))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.SendVerificationEmail(rr, req)

			if tt.expectedStatus == http.StatusBadRequest {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			} else {
				assert.NotEqual(t, http.StatusBadRequest, rr.Code)
			}
		})
	}
}

// TestSend2FAOTPRequestValidation tests 2FA OTP webhook
func TestSend2FAOTPRequestValidation(t *testing.T) {
	os.Unsetenv("WEBHOOK_SECRET")
	handler := NewWebhookHandler(nil)

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "valid request",
			body:           `{"email":"test@example.com","name":"Test User","otp":"123456"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid json",
			body:           `otp: 123456`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/webhooks/send-2fa-otp", bytes.NewReader([]byte(tt.body)))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.Send2FAOTP(rr, req)

			if tt.expectedStatus == http.StatusBadRequest {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			} else {
				assert.NotEqual(t, http.StatusBadRequest, rr.Code)
			}
		})
	}
}

// TestParseUserAgent tests user agent parsing
func TestParseUserAgent(t *testing.T) {
	tests := []struct {
		name     string
		ua       string
		expected string
	}{
		{
			name:     "empty user agent",
			ua:       "",
			expected: "Unbekanntes Gerät",
		},
		{
			name:     "curl",
			ua:       "curl/7.64.1",
			expected: "Kommandozeile (curl)",
		},
		{
			name:     "wget",
			ua:       "Wget/1.21",
			expected: "Kommandozeile (wget)",
		},
		{
			name:     "postman",
			ua:       "PostmanRuntime/7.28.0",
			expected: "API-Client (Postman)",
		},
		{
			name:     "insomnia",
			ua:       "insomnia/2021.5.0",
			expected: "API-Client (Insomnia)",
		},
		{
			name:     "chrome on mac",
			ua:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			expected: "Chrome auf macOS",
		},
		{
			name:     "safari on ios",
			ua:       "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Mobile/15E148 Safari/604.1",
			expected: "Safari auf iOS",
		},
		{
			name:     "firefox on windows",
			ua:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
			expected: "Firefox auf Windows",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseUserAgent(tt.ua)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestMapMethodToGerman tests 2FA method translation
func TestMapMethodToGerman(t *testing.T) {
	tests := []struct {
		method   string
		expected string
	}{
		{"passkey", "Passkey"},
		{"totp", "Authenticator-App (TOTP)"},
		{"unknown", "Zwei-Faktor-Authentifizierung"},
		{"", "Zwei-Faktor-Authentifizierung"},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			result := mapMethodToGerman(tt.method)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCoalesce tests the coalesce helper function
func TestCoalesce(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		fallback string
		expected string
	}{
		{
			name:     "value present",
			value:    "hello",
			fallback: "default",
			expected: "hello",
		},
		{
			name:     "value empty",
			value:    "",
			fallback: "default",
			expected: "default",
		},
		{
			name:     "both empty",
			value:    "",
			fallback: "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := coalesce(tt.value, tt.fallback)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestComputeHMAC tests HMAC computation
func TestComputeHMAC(t *testing.T) {
	// Test vector from RFC 4231
	message := "test"
	key := "secret"

	result := ComputeHMAC(message, key)

	// Verify it's a valid hex string
	assert.Regexp(t, "^[a-f0-9]+$", result)
	// HMAC-SHA256 produces 64 character hex string
	assert.Len(t, result, 64)

	// Same input should produce same output
	result2 := ComputeHMAC(message, key)
	assert.Equal(t, result, result2)

	// Different key should produce different output
	result3 := ComputeHMAC(message, "different-key")
	assert.NotEqual(t, result, result3)
}

// TestGetEnvOrDefault tests environment variable helper
func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		setValue string
		setEnv   bool
		fallback string
		expected string
	}{
		{
			name:     "env set",
			key:      "TEST_ENV_VAR",
			setValue: "custom-value",
			setEnv:   true,
			fallback: "default",
			expected: "custom-value",
		},
		{
			name:     "env not set",
			key:      "TEST_ENV_VAR_NOT_SET",
			setEnv:   false,
			fallback: "default",
			expected: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.setValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnvOrDefault(tt.key, tt.fallback)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGetEnvAsIntOrDefault tests integer environment variable helper
func TestGetEnvAsIntOrDefault(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		setValue string
		setEnv   bool
		fallback int
		expected int
	}{
		{
			name:     "valid int",
			key:      "TEST_INT_VAR",
			setValue: "8080",
			setEnv:   true,
			fallback: 3000,
			expected: 8080,
		},
		{
			name:     "invalid int",
			key:      "TEST_INT_VAR_INVALID",
			setValue: "not-a-number",
			setEnv:   true,
			fallback: 3000,
			expected: 3000,
		},
		{
			name:     "not set",
			key:      "TEST_INT_VAR_NOT_SET",
			setEnv:   false,
			fallback: 3000,
			expected: 3000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				os.Setenv(tt.key, tt.setValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnvAsIntOrDefault(tt.key, tt.fallback)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestWebhookHandlerCreation tests handler initialization
func TestWebhookHandlerCreation(t *testing.T) {
	// Test with default values
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("SMTP_FROM")
	os.Unsetenv("NEXT_PUBLIC_APP_URL")

	handler := NewWebhookHandler(nil)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.mailer)
	assert.NotNil(t, handler.config)
	assert.Equal(t, "noreply@localhost", handler.config.smtpFrom)
	assert.Equal(t, "http://localhost:3000", handler.config.appURL)
	assert.Equal(t, "http://localhost:3000/settings", handler.config.settingsURL)
}

// TestWebhookHandlerWithCustomEnv tests handler with custom environment
func TestWebhookHandlerWithCustomEnv(t *testing.T) {
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_FROM", "no-reply@example.com")
	os.Setenv("NEXT_PUBLIC_APP_URL", "https://app.example.com")
	defer func() {
		os.Unsetenv("SMTP_HOST")
		os.Unsetenv("SMTP_PORT")
		os.Unsetenv("SMTP_FROM")
		os.Unsetenv("NEXT_PUBLIC_APP_URL")
	}()

	handler := NewWebhookHandler(nil)

	assert.NotNil(t, handler)
	assert.Equal(t, "no-reply@example.com", handler.config.smtpFrom)
	assert.Equal(t, "https://app.example.com", handler.config.appURL)
	assert.Equal(t, "https://app.example.com/settings", handler.config.settingsURL)
}

// TestRespondJSON tests JSON response helper
func TestRespondJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	data := MessageResponse{Message: "test message"}
	respondJSON(rr, data)

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response MessageResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "test message", response.Message)
}

// TestRespondError tests error response helper
func TestRespondError(t *testing.T) {
	rr := httptest.NewRecorder()

	respondError(rr, http.StatusBadRequest, "something went wrong")

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "something went wrong", response.Error)
}
