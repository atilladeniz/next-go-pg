package http

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

func TestVerifyWebhookSecret(t *testing.T) {
	tests := []struct {
		name, provided, expected string
		want                     bool
	}{
		{"matching", "k", "k", true},
		{"different", "a", "b", false},
		{"empty provided", "", "k", false},
		{"empty expected", "k", "", false},
		{"both empty", "", "", false},
		{"different length", "secret", "secrets", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, VerifyWebhookSecret(tt.provided, tt.expected))
		})
	}
}

func TestWebhookSecretVerification(t *testing.T) {
	tests := []struct {
		name           string
		webhookSecret  string
		providedSecret string
		expectedStatus int
	}{
		{"valid secret", "test-secret-123", "test-secret-123", http.StatusOK},
		{"invalid secret", "test-secret-123", "wrong-secret", http.StatusUnauthorized},
		{"missing header", "test-secret-123", "", http.StatusUnauthorized},
		{"no secret configured", "", "", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.webhookSecret != "" {
				os.Setenv("WEBHOOK_SECRET", tt.webhookSecret)
				defer os.Unsetenv("WEBHOOK_SECRET")
			} else {
				os.Unsetenv("WEBHOOK_SECRET")
			}

			handler := NewHandler(nil, nil)
			payload := SendMagicLinkRequest{Email: "t@e.com", URL: "https://e.com/m"}
			body, _ := json.Marshal(payload)
			req := httptest.NewRequest(http.MethodPost, "/webhooks/send-magic-link", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.providedSecret != "" {
				req.Header.Set("X-Webhook-Secret", tt.providedSecret)
			}

			rr := httptest.NewRecorder()
			handler.SendMagicLink(rr, req)

			if tt.expectedStatus == http.StatusUnauthorized {
				assert.Equal(t, http.StatusUnauthorized, rr.Code)
				var resp ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Equal(t, "invalid webhook secret", resp.Error)
			} else {
				assert.NotEqual(t, http.StatusUnauthorized, rr.Code)
			}
		})
	}
}

func TestParseUserAgent(t *testing.T) {
	tests := []struct {
		name, ua, expected string
	}{
		{"empty", "", "Unbekanntes Gerät"},
		{"curl", "curl/7.64.1", "Kommandozeile (curl)"},
		{"wget", "Wget/1.21", "Kommandozeile (wget)"},
		{"postman", "PostmanRuntime/7.28.0", "API-Client (Postman)"},
		{"insomnia", "insomnia/2021.5.0", "API-Client (Insomnia)"},
		{"chrome mac", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36", "Chrome auf macOS"},
		{"safari ios", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Mobile/15E148 Safari/604.1", "Safari auf iOS"},
		{"firefox win", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0", "Firefox auf Windows"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, parseUserAgent(tt.ua))
		})
	}
}

func TestMapMethodToGerman(t *testing.T) {
	tests := []struct{ method, expected string }{
		{"passkey", "Passkey"},
		{"totp", "Authenticator-App (TOTP)"},
		{"unknown", "Zwei-Faktor-Authentifizierung"},
		{"", "Zwei-Faktor-Authentifizierung"},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			assert.Equal(t, tt.expected, mapMethodToGerman(tt.method))
		})
	}
}

func TestCoalesce(t *testing.T) {
	assert.Equal(t, "a", coalesce("a", "b"))
	assert.Equal(t, "b", coalesce("", "b"))
	assert.Equal(t, "", coalesce("", ""))
	assert.Equal(t, "b", coalesce("   ", "b"))
}

func TestComputeHMAC(t *testing.T) {
	got1 := ComputeHMAC("hello", "key")
	got2 := ComputeHMAC("hello", "key")
	got3 := ComputeHMAC("hello", "different")
	assert.Equal(t, got1, got2)
	assert.NotEqual(t, got1, got3)
}
