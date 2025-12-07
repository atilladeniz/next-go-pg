// Package jobs defines background job workers for the River queue.
package jobs

import (
	"github.com/riverqueue/river"

	"github.com/atilladeniz/next-go-pg/backend/internal/repository"
	"github.com/atilladeniz/next-go-pg/backend/internal/sse"
)

// WorkerDeps holds dependencies for job workers.
type WorkerDeps struct {
	EmailConfig *EmailConfig
	SSEBroker   *sse.Broker
	ExportStore *ExportStore
	StatsRepo   *repository.UserStatsRepository
}

// RegisterWorkers registers all job workers with the given workers registry.
func RegisterWorkers(workers *river.Workers, deps *WorkerDeps) {
	// Email workers
	if deps.EmailConfig != nil {
		river.AddWorker(workers, NewSendMagicLinkWorker(deps.EmailConfig))
		river.AddWorker(workers, NewSendVerificationEmailWorker(deps.EmailConfig))
		river.AddWorker(workers, NewSend2FAOTPWorker(deps.EmailConfig))
		river.AddWorker(workers, NewSendLoginNotificationWorker(deps.EmailConfig))
	}

	// Export workers
	if deps.SSEBroker != nil && deps.ExportStore != nil {
		river.AddWorker(workers, NewDataExportWorker(deps.SSEBroker, deps.ExportStore, deps.StatsRepo))
	}
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
