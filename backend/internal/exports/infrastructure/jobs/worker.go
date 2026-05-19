package jobs

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"time"

	"github.com/riverqueue/river"

	exportsapp "github.com/atilladeniz/next-go-pg/backend/internal/exports/application"
	exports "github.com/atilladeniz/next-go-pg/backend/internal/exports/domain"
	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
)

// DataExportWorker runs an export job: gather data via the StatsReader
// ACL, format, persist to the result store, broadcast progress.
type DataExportWorker struct {
	river.WorkerDefaults[DataExportArgs]
	progress exportsapp.ProgressPublisher
	store    exportsapp.Store
	stats    exportsapp.StatsReader
}

func NewDataExportWorker(
	progress exportsapp.ProgressPublisher,
	store exportsapp.Store,
	stats exportsapp.StatsReader,
) *DataExportWorker {
	return &DataExportWorker{progress: progress, store: store, stats: stats}
}

func (w *DataExportWorker) Work(ctx context.Context, job *river.Job[DataExportArgs]) error {
	args := job.Args

	logger.Info().
		Str("job_id", args.JobID).
		Str("user_id", args.UserID).
		Str("format", string(args.Format)).
		Str("data_type", args.DataType).
		Msg("Starting data export job")

	w.sendProgress(ProgressUpdate{
		JobID:    args.JobID,
		Status:   exports.StatusProcessing,
		Progress: 0,
		Message:  "Export wird vorbereitet...",
	})
	time.Sleep(500 * time.Millisecond)

	w.sendProgress(ProgressUpdate{
		JobID:    args.JobID,
		Status:   exports.StatusProcessing,
		Progress: 20,
		Message:  "Daten werden gesammelt...",
	})

	data := w.gather(ctx, args.UserID, args.DataType)
	time.Sleep(500 * time.Millisecond)

	w.sendProgress(ProgressUpdate{
		JobID:    args.JobID,
		Status:   exports.StatusProcessing,
		Progress: 50,
		Message:  "Daten werden verarbeitet...",
	})

	var (
		exportData  []byte
		contentType string
		fileName    string
		err         error
	)
	switch args.Format {
	case exports.FormatCSV:
		exportData, err = toCSV(data)
		contentType = "text/csv"
		fileName = fmt.Sprintf("export_%s_%s.csv", args.DataType, time.Now().Format("2006-01-02"))
	case exports.FormatJSON:
		exportData, err = json.MarshalIndent(data, "", "  ")
		contentType = "application/json"
		fileName = fmt.Sprintf("export_%s_%s.json", args.DataType, time.Now().Format("2006-01-02"))
	default:
		err = fmt.Errorf("unsupported format: %s", args.Format)
	}
	if err != nil {
		logger.Error().Err(err).Str("job_id", args.JobID).Msg("Export conversion failed")
		w.sendProgress(ProgressUpdate{
			JobID:    args.JobID,
			Status:   exports.StatusFailed,
			Progress: 0,
			Message:  "Export fehlgeschlagen",
			Error:    err.Error(),
		})
		return err
	}

	time.Sleep(500 * time.Millisecond)
	w.sendProgress(ProgressUpdate{
		JobID:    args.JobID,
		Status:   exports.StatusProcessing,
		Progress: 80,
		Message:  "Export wird finalisiert...",
	})

	downloadID := fmt.Sprintf("%s_%d", args.JobID, time.Now().UnixNano())
	w.store.Save(downloadID, &exportsapp.Result{
		Data:        exportData,
		ContentType: contentType,
		FileName:    fileName,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	})

	w.sendProgress(ProgressUpdate{
		JobID:      args.JobID,
		Status:     exports.StatusCompleted,
		Progress:   100,
		Message:    "Export abgeschlossen!",
		FileName:   fileName,
		DownloadID: downloadID,
	})

	logger.Info().
		Str("job_id", args.JobID).
		Str("download_id", downloadID).
		Str("file_name", fileName).
		Int("size_bytes", len(exportData)).
		Msg("Data export completed")
	return nil
}

func (w *DataExportWorker) sendProgress(update ProgressUpdate) {
	if w.progress == nil {
		return
	}
	payload, _ := json.Marshal(update)
	w.progress.Broadcast("export-progress", string(payload))
}

func (w *DataExportWorker) gather(ctx context.Context, userID, dataType string) []map[string]any {
	now := time.Now()
	var snapshot exportsapp.StatsSnapshot
	if w.stats != nil {
		if s, err := w.stats.Read(ctx, userID); err == nil {
			snapshot = s
		}
	}
	switch dataType {
	case "stats":
		return []map[string]any{
			{
				"datum":              now.Format("02.01.2006"),
				"projekte":           snapshot.Projects,
				"aktivitaet":         snapshot.Activity,
				"benachrichtigungen": snapshot.Notifications,
			},
		}
	case "activity":
		return []map[string]any{
			{
				"zeitpunkt": now.Format("02.01.2006 15:04"),
				"aktion":    "Export erstellt",
				"details":   fmt.Sprintf("Aktivität heute: %d", snapshot.Activity),
			},
		}
	default:
		return []map[string]any{
			{"kategorie": "Dashboard", "metrik": "Projekte", "wert": snapshot.Projects},
			{"kategorie": "Dashboard", "metrik": "Aktivität Heute", "wert": snapshot.Activity},
			{"kategorie": "Dashboard", "metrik": "Benachrichtigungen", "wert": snapshot.Notifications},
		}
	}
}

func toCSV(data []map[string]any) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("no data to export")
	}
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	headers := make([]string, 0, len(data[0]))
	for key := range data[0] {
		headers = append(headers, key)
	}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}
	for _, row := range data {
		values := make([]string, 0, len(headers))
		for _, h := range headers {
			values = append(values, fmt.Sprintf("%v", row[h]))
		}
		if err := writer.Write(values); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	return buf.Bytes(), writer.Error()
}
