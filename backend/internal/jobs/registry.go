// Package jobs defines background job workers for the River queue.
package jobs

import (
	"github.com/riverqueue/river"

	"github.com/atilladeniz/next-go-pg/backend/internal/application"
)

// WorkerDeps holds dependencies for job workers.
type WorkerDeps struct {
	EmailSender application.EmailSender
	Events      application.EventBroadcaster
	ExportStore *ExportStore
	StatsRepo   application.StatsRepository
}

// RegisterWorkers registers all job workers with the given workers registry.
func RegisterWorkers(workers *river.Workers, deps *WorkerDeps) {
	if deps.EmailSender != nil {
		river.AddWorker(workers, NewSendMagicLinkWorker(deps.EmailSender))
		river.AddWorker(workers, NewSendVerificationEmailWorker(deps.EmailSender))
		river.AddWorker(workers, NewSend2FAOTPWorker(deps.EmailSender))
		river.AddWorker(workers, NewSendLoginNotificationWorker(deps.EmailSender))
	}

	if deps.Events != nil && deps.ExportStore != nil {
		river.AddWorker(workers, NewDataExportWorker(deps.Events, deps.ExportStore, deps.StatsRepo))
	}
}
