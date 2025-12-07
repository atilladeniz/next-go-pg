# Background Jobs with River

This project uses [River](https://riverqueue.com) for background job processing. River is a PostgreSQL-native job queue for Go, offering:

- **Transactional guarantees**: Jobs are committed atomically with your data
- **High performance**: 46,000-66,000 jobs/second
- **Native PostgreSQL**: No Redis required, uses existing database
- **Automatic retries**: Built-in exponential backoff

## Architecture

```
┌─────────────┐      ┌──────────────┐      ┌─────────────┐
│  HTTP API   │ ──── │ Job Enqueue  │ ──── │  PostgreSQL │
│  (webhook)  │      │  (Insert)    │      │  (river.*)  │
└─────────────┘      └──────────────┘      └──────────────┘
                                                   │
                                           ┌───────────────┐
                                           │  River Worker │
                                           │  (Background) │
                                           └───────────────┘
```

## Email Jobs

The following email types are processed as background jobs:

| Job Kind | Description | Webhook Endpoint |
|----------|-------------|------------------|
| `send_magic_link` | Magic link login emails | `/api/v1/webhooks/send-magic-link` |
| `send_verification_email` | Email verification | `/api/v1/webhooks/send-verification-email` |
| `send_2fa_otp` | 2FA one-time passwords | `/api/v1/webhooks/send-2fa-otp` |
| `send_login_notification` | New device login alerts | `/api/v1/webhooks/session-created` |

## Usage

### Running Migrations

River creates its own schema in PostgreSQL. Run migrations before starting the server:

```bash
cd backend

# Run River migrations
make river-migrate-up

# Check migration status
make river-migrate-version

# Rollback (if needed)
make river-migrate-down
```

### Server Integration

River is automatically initialized when the server starts:

1. Creates pgx connection pool
2. Runs River migrations
3. Registers email workers
4. Starts background job processing

On shutdown, River gracefully stops processing and waits for running jobs.

### Fallback Behavior

If River is unavailable (e.g., database connection issues), webhooks fall back to synchronous email sending. This ensures email delivery even during infrastructure problems.

## Creating New Job Types

### 1. Define Job Arguments

```go
// internal/jobs/email.go

type SendWelcomeEmailArgs struct {
    Email    string `json:"email"`
    UserName string `json:"userName"`
}

func (SendWelcomeEmailArgs) Kind() string { return "send_welcome_email" }
```

### 2. Implement Worker

```go
type SendWelcomeEmailWorker struct {
    river.WorkerDefaults[SendWelcomeEmailArgs]
    mailer *gomail.Dialer
    config *EmailConfig
}

func (w *SendWelcomeEmailWorker) Work(ctx context.Context, job *river.Job[SendWelcomeEmailArgs]) error {
    args := job.Args

    // Render template and send email
    body, err := templates.RenderWelcome(...)
    if err != nil {
        return fmt.Errorf("render template: %w", err)
    }

    if err := w.sendEmail(args.Email, "Welcome!", body); err != nil {
        return fmt.Errorf("send email: %w", err)
    }

    logger.Info().Str("email", args.Email).Msg("Welcome email sent")
    return nil
}
```

### 3. Register Worker

```go
// internal/jobs/registry.go

func RegisterWorkers(workers *river.Workers, config *EmailConfig) {
    river.AddWorker(workers, NewSendMagicLinkWorker(config))
    river.AddWorker(workers, NewSendVerificationEmailWorker(config))
    // Add your new worker:
    river.AddWorker(workers, NewSendWelcomeEmailWorker(config))
}
```

### 4. Create Enqueue Helper

```go
// internal/jobs/enqueue.go

func EnqueueWelcomeEmail(ctx context.Context, enqueuer JobEnqueuer, email, userName string) error {
    _, err := enqueuer.Insert(ctx, SendWelcomeEmailArgs{
        Email:    email,
        UserName: userName,
    }, nil)
    return err
}
```

## Configuration

River uses default settings optimized for most workloads:

```go
// pkg/river/client.go

type Config struct {
    MaxWorkers int            // Default: 100
    Queues     map[string]int // Custom queue configurations
}
```

For high-volume scenarios, consider:

- Running dedicated worker processes
- Using multiple queues with different worker counts
- Enabling River Pro for global concurrency limits

## Monitoring

### Database Queries

Check pending jobs:
```sql
SELECT kind, state, COUNT(*)
FROM river.job
GROUP BY kind, state;
```

Check failed jobs:
```sql
SELECT id, kind, errors, created_at
FROM river.job
WHERE state = 'discarded'
ORDER BY created_at DESC;
```

### Logs

River operations are logged with structured logging:

```json
{"level":"info","msg":"Magic link email enqueued","email":"user@example.com"}
{"level":"info","msg":"Magic link email sent via background job","email":"user@example.com"}
```

## Troubleshooting

### Jobs Not Processing

1. Check if River is started:
   ```
   grep "River job queue initialized" /var/log/backend.log
   ```

2. Verify migrations:
   ```bash
   make river-migrate-version
   ```

3. Check for errors:
   ```sql
   SELECT * FROM river.job WHERE state IN ('retryable', 'discarded');
   ```

### High Latency

River uses LISTEN/NOTIFY for sub-millisecond job pickup. If experiencing latency:

1. Ensure direct PostgreSQL connection (not through PgBouncer transaction pooling)
2. Check worker concurrency settings
3. Monitor PostgreSQL connection pool

## References

- [River Documentation](https://riverqueue.com/docs)
- [River GitHub](https://github.com/riverqueue/river)
- [River Pro](https://riverqueue.com/pro) (commercial features)
