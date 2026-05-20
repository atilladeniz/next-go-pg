package jobs

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"

	exportsapp "github.com/atilladeniz/next-go-pg/backend/internal/exports/application"
	exports "github.com/atilladeniz/next-go-pg/backend/internal/exports/domain"
)

// RiverClient is the subset of *river.Client this adapter needs.
type RiverClient interface {
	InsertTx(ctx context.Context, tx pgx.Tx, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
	Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
}

// Enqueuer is the exports context's River-backed JobEnqueuer.
type Enqueuer struct {
	client RiverClient
}

var _ exportsapp.JobEnqueuer = (*Enqueuer)(nil)

func NewEnqueuer(client RiverClient) *Enqueuer {
	return &Enqueuer{client: client}
}

func (e *Enqueuer) EnqueueDataExport(ctx context.Context, jobID, userID, format, dataType string) error {
	_, err := e.client.Insert(ctx, DataExportArgs{
		JobID:    jobID,
		UserID:   userID,
		Format:   exports.Format(format),
		DataType: dataType,
	}, nil)
	return err
}

// Register hooks this context's worker into a River workers registry.
func Register(
	workers *river.Workers,
	progress exportsapp.ProgressPublisher,
	store exportsapp.Store,
	stats exportsapp.StatsReader,
) {
	river.AddWorker(workers, NewDataExportWorker(progress, store, stats))
}
