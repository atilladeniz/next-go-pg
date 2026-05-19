package jobs

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"

	"github.com/atilladeniz/next-go-pg/backend/internal/application"
)

// riverClient is the subset of *river.Client the enqueuer needs. It's
// satisfied by River's real client; declaring it lets us swap the
// dependency for tests without pulling in River's heavy types.
type riverClient interface {
	InsertTx(ctx context.Context, tx pgx.Tx, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
	Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
}

// Enqueuer is the River-backed adapter for application.JobEnqueuer.
type Enqueuer struct {
	client riverClient
}

var _ application.JobEnqueuer = (*Enqueuer)(nil)

// NewEnqueuer wraps a River client into the application port.
func NewEnqueuer(client riverClient) *Enqueuer {
	return &Enqueuer{client: client}
}

func (e *Enqueuer) EnqueueMagicLink(ctx context.Context, email, url string) error {
	_, err := e.client.Insert(ctx, SendMagicLinkArgs{Email: email, URL: url}, nil)
	return err
}

func (e *Enqueuer) EnqueueVerificationEmail(ctx context.Context, email, name, url string) error {
	_, err := e.client.Insert(ctx, SendVerificationEmailArgs{Email: email, Name: name, URL: url}, nil)
	return err
}

func (e *Enqueuer) Enqueue2FAOTP(ctx context.Context, email, name, otp string) error {
	_, err := e.client.Insert(ctx, Send2FAOTPArgs{Email: email, Name: name, OTP: otp}, nil)
	return err
}

func (e *Enqueuer) EnqueueLoginNotification(ctx context.Context, email, userName, device, ipAddress string) error {
	_, err := e.client.Insert(ctx, SendLoginNotificationArgs{
		Email:     email,
		UserName:  userName,
		Device:    device,
		IPAddress: ipAddress,
	}, nil)
	return err
}

func (e *Enqueuer) EnqueueDataExport(ctx context.Context, jobID, userID, format, dataType string) error {
	_, err := e.client.Insert(ctx, DataExportArgs{
		JobID:    jobID,
		UserID:   userID,
		Format:   ExportFormat(format),
		DataType: dataType,
	}, nil)
	return err
}
