package persistence

import (
	ai "github.com/atilladeniz/next-go-pg/backend/internal/aiworkflows/domain"
	shared "github.com/atilladeniz/next-go-pg/backend/internal/shared/domain"
)

func toDomain(m gormRepoSummary) (*ai.RepoSummary, error) {
	status, err := ai.NewStatus(m.Status)
	if err != nil {
		return nil, err
	}
	files := make([]ai.FileSummary, 0, len(m.Files))
	for _, r := range m.Files {
		fs, err := ai.NewFileSummary(r.Filename, r.Summary)
		if err != nil {
			return nil, err
		}
		files = append(files, fs)
	}
	url, err := ai.NewRepoURL(m.RepoURL)
	if err != nil {
		return nil, err
	}
	durations := make(map[string]int64, len(m.StepDurations))
	for k, v := range m.StepDurations {
		durations[k] = v
	}
	return &ai.RepoSummary{
		ID:            m.ID,
		UserID:        shared.UserID(m.UserID),
		RepoURL:       url,
		Status:        status,
		Files:         files,
		Summary:       m.Summary,
		FailReason:    m.FailReason,
		StartedAt:     m.StartedAt,
		CompletedAt:   m.CompletedAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		StepDurations: durations,
	}, nil
}

func fromDomain(d *ai.RepoSummary) gormRepoSummary {
	files := make(fileSummariesJSON, 0, len(d.Files))
	for _, fs := range d.Files {
		files = append(files, fileSummaryRecord{
			Filename: fs.Filename(),
			Summary:  fs.Summary(),
		})
	}
	durations := make(stepDurationsJSON, len(d.StepDurations))
	for k, v := range d.StepDurations {
		durations[k] = v
	}
	return gormRepoSummary{
		ID:            d.ID,
		UserID:        d.UserID.String(),
		RepoURL:       d.RepoURL.String(),
		Status:        d.Status.String(),
		Files:         files,
		Summary:       d.Summary,
		FailReason:    d.FailReason,
		StepDurations: durations,
		StartedAt:     d.StartedAt,
		CompletedAt:   d.CompletedAt,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}
