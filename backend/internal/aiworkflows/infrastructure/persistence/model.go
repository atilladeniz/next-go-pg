// Package persistence holds the GORM-backed adapter for the aiworkflows
// bounded context. The gormRepoSummary twin is intentionally unexported —
// callers exchange domain types via the mapper functions.
package persistence

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type gormRepoSummary struct {
	ID          uint              `gorm:"primaryKey"`
	UserID      string            `gorm:"index;not null"`
	RepoURL     string            `gorm:"not null"`
	Status      string            `gorm:"index;not null"`
	Files       fileSummariesJSON `gorm:"type:jsonb;default:'[]'"`
	Summary     string            `gorm:"type:text"`
	FailReason  string            `gorm:"type:text"`
	StartedAt   time.Time
	CompletedAt time.Time
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (gormRepoSummary) TableName() string { return "repo_summaries" }

// fileSummaryRecord is the JSONB shape persisted for each per-file
// summary. We do NOT serialize the domain FileSummary directly because
// its fields are unexported (deliberate — invariants enforced via
// constructor).
type fileSummaryRecord struct {
	Filename string `json:"filename"`
	Summary  string `json:"summary"`
}

// fileSummariesJSON is a slice of fileSummaryRecord with GORM
// Valuer/Scanner methods so it round-trips through a JSONB column.
type fileSummariesJSON []fileSummaryRecord

func (f fileSummariesJSON) Value() (driver.Value, error) {
	if f == nil {
		return "[]", nil
	}
	return json.Marshal(f)
}

func (f *fileSummariesJSON) Scan(src any) error {
	if src == nil {
		*f = nil
		return nil
	}
	var raw []byte
	switch v := src.(type) {
	case []byte:
		raw = v
	case string:
		raw = []byte(v)
	default:
		return errors.New("fileSummariesJSON: unsupported scan source")
	}
	if len(raw) == 0 {
		*f = nil
		return nil
	}
	return json.Unmarshal(raw, f)
}

// Entities returns the GORM models that AutoMigrate must process for
// the aiworkflows context. Called from composition.runAutoMigrations.
func Entities() []any {
	return []any{&gormRepoSummary{}}
}
