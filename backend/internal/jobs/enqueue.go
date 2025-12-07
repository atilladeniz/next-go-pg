package jobs

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

// JobEnqueuer defines the interface for enqueueing jobs.
type JobEnqueuer interface {
	InsertTx(ctx context.Context, tx pgx.Tx, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
	Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
}

// EnqueueMagicLink enqueues a magic link email job.
func EnqueueMagicLink(ctx context.Context, enqueuer JobEnqueuer, email, url string) error {
	_, err := enqueuer.Insert(ctx, SendMagicLinkArgs{
		Email: email,
		URL:   url,
	}, nil)
	return err
}

// EnqueueVerificationEmail enqueues a verification email job.
func EnqueueVerificationEmail(ctx context.Context, enqueuer JobEnqueuer, email, name, url string) error {
	_, err := enqueuer.Insert(ctx, SendVerificationEmailArgs{
		Email: email,
		Name:  name,
		URL:   url,
	}, nil)
	return err
}

// Enqueue2FAOTP enqueues a 2FA OTP email job.
func Enqueue2FAOTP(ctx context.Context, enqueuer JobEnqueuer, email, name, otp string) error {
	_, err := enqueuer.Insert(ctx, Send2FAOTPArgs{
		Email: email,
		Name:  name,
		OTP:   otp,
	}, nil)
	return err
}

// EnqueueLoginNotification enqueues a login notification email job.
func EnqueueLoginNotification(ctx context.Context, enqueuer JobEnqueuer, email, userName, device, ipAddress string) error {
	_, err := enqueuer.Insert(ctx, SendLoginNotificationArgs{
		Email:     email,
		UserName:  userName,
		Device:    device,
		IPAddress: ipAddress,
	}, nil)
	return err
}
