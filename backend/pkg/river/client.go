// Package river provides a River job queue client for background job processing.
package river

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"

	"github.com/atilladeniz/next-go-pg/backend/pkg/logger"
)

// Client wraps the River client with convenience methods.
type Client struct {
	*river.Client[pgx.Tx]
	pool *pgxpool.Pool
}

// Config holds the configuration for the River client.
type Config struct {
	// MaxWorkers is the maximum number of concurrent workers per queue.
	// Default: 100
	MaxWorkers int

	// Queues defines custom queue configurations.
	// If empty, only the default queue will be used.
	Queues map[string]int
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		MaxWorkers: 100,
		Queues:     nil, // Use default queue only
	}
}

// NewClient creates a new River client.
func NewClient(ctx context.Context, pool *pgxpool.Pool, workers *river.Workers, cfg Config) (*Client, error) {
	if cfg.MaxWorkers <= 0 {
		cfg.MaxWorkers = 100
	}

	// Build queue configuration
	queues := make(map[string]river.QueueConfig)
	if len(cfg.Queues) == 0 {
		queues[river.QueueDefault] = river.QueueConfig{MaxWorkers: cfg.MaxWorkers}
	} else {
		for name, maxWorkers := range cfg.Queues {
			queues[name] = river.QueueConfig{MaxWorkers: maxWorkers}
		}
	}

	riverClient, err := river.NewClient(riverpgxv5.New(pool), &river.Config{
		Queues:  queues,
		Workers: workers,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create river client: %w", err)
	}

	return &Client{
		Client: riverClient,
		pool:   pool,
	}, nil
}

// NewInsertOnlyClient creates a River client that can only insert jobs (no workers).
// Use this for frontend processes that enqueue jobs but don't process them.
func NewInsertOnlyClient(ctx context.Context, pool *pgxpool.Pool) (*Client, error) {
	riverClient, err := river.NewClient(riverpgxv5.New(pool), &river.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create insert-only river client: %w", err)
	}

	return &Client{
		Client: riverClient,
		pool:   pool,
	}, nil
}

// Start starts the River client and begins processing jobs.
func (c *Client) Start(ctx context.Context) error {
	logger.Info().Msg("Starting River job queue client")
	return c.Client.Start(ctx)
}

// Stop stops the River client gracefully.
func (c *Client) Stop(ctx context.Context) error {
	logger.Info().Msg("Stopping River job queue client")
	return c.Client.Stop(ctx)
}

// RunMigrations runs River's database migrations.
func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migrator, err := rivermigrate.New(riverpgxv5.New(pool), nil)
	if err != nil {
		return fmt.Errorf("failed to create river migrator: %w", err)
	}

	res, err := migrator.Migrate(ctx, rivermigrate.DirectionUp, nil)
	if err != nil {
		return fmt.Errorf("failed to run river migrations: %w", err)
	}

	for _, version := range res.Versions {
		logger.Info().
			Int("version", version.Version).
			Str("direction", "up").
			Msg("River migration applied")
	}

	return nil
}

// IsMigrated checks if River tables exist in the database.
func IsMigrated(ctx context.Context, pool *pgxpool.Pool) (bool, error) {
	// Check if river_job table exists
	var exists bool
	err := pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema = 'river'
			AND table_name = 'job'
		)
	`).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check river tables: %w", err)
	}
	return exists, nil
}
