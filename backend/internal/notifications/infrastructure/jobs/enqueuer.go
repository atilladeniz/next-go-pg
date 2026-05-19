package jobs

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"

	notifapp "github.com/atilladeniz/next-go-pg/backend/internal/notifications/application"
)

// RiverClient is the subset of *river.Client this adapter needs.
// River's real client satisfies it implicitly.
type RiverClient interface {
	InsertTx(ctx context.Context, tx pgx.Tx, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
	Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
}

// Enqueuer is the notifications context's River-backed JobEnqueuer.
type Enqueuer struct {
	client RiverClient
}

var _ notifapp.JobEnqueuer = (*Enqueuer)(nil)

func NewEnqueuer(client RiverClient) *Enqueuer {
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
