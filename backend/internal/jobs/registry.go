// Package jobs defines background job workers for the River queue.
package jobs

import (
	"github.com/riverqueue/river"
)

// RegisterWorkers registers all job workers with the given workers registry.
func RegisterWorkers(workers *river.Workers, config *EmailConfig) {
	// Email workers
	river.AddWorker(workers, NewSendMagicLinkWorker(config))
	river.AddWorker(workers, NewSendVerificationEmailWorker(config))
	river.AddWorker(workers, NewSend2FAOTPWorker(config))
	river.AddWorker(workers, NewSendLoginNotificationWorker(config))
}

// NewEmailConfig creates a new EmailConfig from environment settings.
func NewEmailConfig(smtpHost string, smtpPort int, smtpFrom, appURL string) *EmailConfig {
	settingsURL := appURL + "/settings"
	return &EmailConfig{
		SMTPHost:    smtpHost,
		SMTPPort:    smtpPort,
		SMTPFrom:    smtpFrom,
		AppURL:      appURL,
		SettingsURL: settingsURL,
	}
}
