package jobs

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"time"

	"github.com/riverqueue/river"

	"github.com/atilladeniz/next-go-pg/backend/internal/sse"
	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
)

// ExportFormat defines the export file format
type ExportFormat string

const (
	ExportFormatCSV  ExportFormat = "csv"
	ExportFormatJSON ExportFormat = "json"
)

// ExportStatus represents the current state of an export job
type ExportStatus string

const (
	ExportStatusPending    ExportStatus = "pending"
	ExportStatusProcessing ExportStatus = "processing"
	ExportStatusCompleted  ExportStatus = "completed"
	ExportStatusFailed     ExportStatus = "failed"
)

// --- Data Export Job ---

// DataExportArgs defines the arguments for a data export job.
type DataExportArgs struct {
	JobID    string       `json:"jobId"`
	UserID   string       `json:"userId"`
	Format   ExportFormat `json:"format"`
	DataType string       `json:"dataType"` // e.g., "stats", "activity", "all"
}

func (DataExportArgs) Kind() string { return "data_export" }

// ExportProgress represents progress updates sent via SSE
type ExportProgress struct {
	JobID      string       `json:"jobId"`
	Status     ExportStatus `json:"status"`
	Progress   int          `json:"progress"` // 0-100
	Message    string       `json:"message"`
	FileName   string       `json:"fileName,omitempty"`
	DownloadID string       `json:"downloadId,omitempty"`
	Error      string       `json:"error,omitempty"`
}

// DataExportWorker processes data export jobs.
type DataExportWorker struct {
	river.WorkerDefaults[DataExportArgs]
	sseBroker   *sse.Broker
	exportStore *ExportStore
}

// ExportStore holds completed exports in memory (in production, use object storage)
type ExportStore struct {
	exports map[string]*ExportResult
}

type ExportResult struct {
	Data        []byte
	ContentType string
	FileName    string
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

func NewExportStore() *ExportStore {
	return &ExportStore{
		exports: make(map[string]*ExportResult),
	}
}

func (s *ExportStore) Store(id string, result *ExportResult) {
	s.exports[id] = result
}

func (s *ExportStore) Get(id string) (*ExportResult, bool) {
	result, ok := s.exports[id]
	if !ok {
		return nil, false
	}
	// Check expiration
	if time.Now().After(result.ExpiresAt) {
		delete(s.exports, id)
		return nil, false
	}
	return result, true
}

func (s *ExportStore) Delete(id string) {
	delete(s.exports, id)
}

func NewDataExportWorker(sseBroker *sse.Broker, exportStore *ExportStore) *DataExportWorker {
	return &DataExportWorker{
		sseBroker:   sseBroker,
		exportStore: exportStore,
	}
}

func (w *DataExportWorker) Work(ctx context.Context, job *river.Job[DataExportArgs]) error {
	args := job.Args

	logger.Info().
		Str("job_id", args.JobID).
		Str("user_id", args.UserID).
		Str("format", string(args.Format)).
		Str("data_type", args.DataType).
		Msg("Starting data export job")

	// Send initial progress
	w.sendProgress(ExportProgress{
		JobID:    args.JobID,
		Status:   ExportStatusProcessing,
		Progress: 0,
		Message:  "Export wird vorbereitet...",
	})

	// Simulate data gathering (in production, fetch from database)
	time.Sleep(500 * time.Millisecond)
	w.sendProgress(ExportProgress{
		JobID:    args.JobID,
		Status:   ExportStatusProcessing,
		Progress: 20,
		Message:  "Daten werden gesammelt...",
	})

	// Generate sample data
	data := w.generateSampleData(args.UserID, args.DataType)

	time.Sleep(500 * time.Millisecond)
	w.sendProgress(ExportProgress{
		JobID:    args.JobID,
		Status:   ExportStatusProcessing,
		Progress: 50,
		Message:  "Daten werden verarbeitet...",
	})

	// Convert to requested format
	var exportData []byte
	var contentType string
	var fileName string
	var err error

	switch args.Format {
	case ExportFormatCSV:
		exportData, err = w.convertToCSV(data)
		contentType = "text/csv"
		fileName = fmt.Sprintf("export_%s_%s.csv", args.DataType, time.Now().Format("2006-01-02"))
	case ExportFormatJSON:
		exportData, err = json.MarshalIndent(data, "", "  ")
		contentType = "application/json"
		fileName = fmt.Sprintf("export_%s_%s.json", args.DataType, time.Now().Format("2006-01-02"))
	default:
		err = fmt.Errorf("unsupported format: %s", args.Format)
	}

	if err != nil {
		logger.Error().Err(err).Str("job_id", args.JobID).Msg("Export conversion failed")
		w.sendProgress(ExportProgress{
			JobID:    args.JobID,
			Status:   ExportStatusFailed,
			Progress: 0,
			Message:  "Export fehlgeschlagen",
			Error:    err.Error(),
		})
		return err
	}

	time.Sleep(500 * time.Millisecond)
	w.sendProgress(ExportProgress{
		JobID:    args.JobID,
		Status:   ExportStatusProcessing,
		Progress: 80,
		Message:  "Export wird finalisiert...",
	})

	// Store the export result
	downloadID := fmt.Sprintf("%s_%d", args.JobID, time.Now().UnixNano())
	w.exportStore.Store(downloadID, &ExportResult{
		Data:        exportData,
		ContentType: contentType,
		FileName:    fileName,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(1 * time.Hour), // Export expires after 1 hour
	})

	// Send completion
	w.sendProgress(ExportProgress{
		JobID:      args.JobID,
		Status:     ExportStatusCompleted,
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

func (w *DataExportWorker) sendProgress(progress ExportProgress) {
	if w.sseBroker == nil {
		return
	}
	data, _ := json.Marshal(progress)
	w.sseBroker.Broadcast("export-progress", string(data))
}

func (w *DataExportWorker) generateSampleData(userID, dataType string) []map[string]interface{} {
	// Generate sample data based on data type
	now := time.Now()

	switch dataType {
	case "stats":
		return []map[string]interface{}{
			{
				"datum":         now.AddDate(0, 0, -6).Format("02.01.2006"),
				"projekte":      12,
				"aufgaben":      45,
				"abgeschlossen": 38,
			},
			{
				"datum":         now.AddDate(0, 0, -5).Format("02.01.2006"),
				"projekte":      13,
				"aufgaben":      52,
				"abgeschlossen": 44,
			},
			{
				"datum":         now.AddDate(0, 0, -4).Format("02.01.2006"),
				"projekte":      13,
				"aufgaben":      58,
				"abgeschlossen": 51,
			},
			{
				"datum":         now.AddDate(0, 0, -3).Format("02.01.2006"),
				"projekte":      14,
				"aufgaben":      61,
				"abgeschlossen": 55,
			},
			{
				"datum":         now.AddDate(0, 0, -2).Format("02.01.2006"),
				"projekte":      14,
				"aufgaben":      67,
				"abgeschlossen": 60,
			},
			{
				"datum":         now.AddDate(0, 0, -1).Format("02.01.2006"),
				"projekte":      15,
				"aufgaben":      72,
				"abgeschlossen": 65,
			},
			{
				"datum":         now.Format("02.01.2006"),
				"projekte":      15,
				"aufgaben":      78,
				"abgeschlossen": 70,
			},
		}
	case "activity":
		return []map[string]interface{}{
			{
				"zeitpunkt": now.Add(-2 * time.Hour).Format("02.01.2006 15:04"),
				"aktion":    "Projekt erstellt",
				"details":   "Neues Projekt 'Website Redesign'",
			},
			{
				"zeitpunkt": now.Add(-1 * time.Hour).Format("02.01.2006 15:04"),
				"aktion":    "Aufgabe abgeschlossen",
				"details":   "UI Mockups fertiggestellt",
			},
			{
				"zeitpunkt": now.Add(-30 * time.Minute).Format("02.01.2006 15:04"),
				"aktion":    "Kommentar hinzugefügt",
				"details":   "Feedback zu Design-Entwürfen",
			},
			{
				"zeitpunkt": now.Add(-15 * time.Minute).Format("02.01.2006 15:04"),
				"aktion":    "Datei hochgeladen",
				"details":   "final_design_v2.fig",
			},
		}
	default: // "all"
		return []map[string]interface{}{
			{
				"kategorie": "Übersicht",
				"metrik":    "Aktive Projekte",
				"wert":      15,
				"trend":     "+2 diese Woche",
			},
			{
				"kategorie": "Übersicht",
				"metrik":    "Offene Aufgaben",
				"wert":      78,
				"trend":     "-5 seit gestern",
			},
			{
				"kategorie": "Übersicht",
				"metrik":    "Abschlussrate",
				"wert":      "89%",
				"trend":     "+3% diesen Monat",
			},
			{
				"kategorie": "Team",
				"metrik":    "Aktive Mitglieder",
				"wert":      8,
				"trend":     "Unverändert",
			},
		}
	}
}

func (w *DataExportWorker) convertToCSV(data []map[string]interface{}) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("no data to export")
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Get headers from first row
	var headers []string
	for key := range data[0] {
		headers = append(headers, key)
	}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// Write data rows
	for _, row := range data {
		var values []string
		for _, header := range headers {
			val := row[header]
			values = append(values, fmt.Sprintf("%v", val))
		}
		if err := writer.Write(values); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), writer.Error()
}
