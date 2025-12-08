package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"

	"github.com/atilladeniz/next-go-pg/backend/pkg/config"
)

func main() {
	// Parse command line flags
	upFlag := flag.Bool("up", false, "Run all up migrations")
	downFlag := flag.Bool("down", false, "Rollback the last migration")
	versionFlag := flag.Bool("version", false, "Print current migration version")
	flag.Parse()

	// Load configuration
	cfg := config.Load()
	dbURL := cfg.GetDatabaseURLForPgx()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to database
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create migrator
	migrator, err := rivermigrate.New(riverpgxv5.New(pool), nil)
	if err != nil {
		log.Fatalf("Failed to create River migrator: %v", err)
	}

	// Execute command
	switch {
	case *upFlag:
		runUp(ctx, migrator)
	case *downFlag:
		runDown(ctx, migrator)
	case *versionFlag:
		printVersion(ctx, pool)
	default:
		printUsage()
	}
}

func runUp(ctx context.Context, migrator *rivermigrate.Migrator[pgx.Tx]) {
	fmt.Println("Running River migrations up...")

	res, err := migrator.Migrate(ctx, rivermigrate.DirectionUp, nil)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	if len(res.Versions) == 0 {
		fmt.Println("No migrations to apply - database is up to date")
		return
	}

	for _, version := range res.Versions {
		fmt.Printf("Applied migration version %d\n", version.Version)
	}
	fmt.Println("River migrations completed successfully")
}

func runDown(ctx context.Context, migrator *rivermigrate.Migrator[pgx.Tx]) {
	fmt.Println("Rolling back River migration...")

	res, err := migrator.Migrate(ctx, rivermigrate.DirectionDown, &rivermigrate.MigrateOpts{
		MaxSteps: 1,
	})
	if err != nil {
		log.Fatalf("Rollback failed: %v", err)
	}

	if len(res.Versions) == 0 {
		fmt.Println("No migrations to rollback")
		return
	}

	for _, version := range res.Versions {
		fmt.Printf("Rolled back migration version %d\n", version.Version)
	}
	fmt.Println("Rollback completed successfully")
}

func printVersion(ctx context.Context, pool *pgxpool.Pool) {
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
		log.Fatalf("Failed to check migration status: %v", err)
	}

	if exists {
		fmt.Println("River migrations are installed (river.job table exists)")
	} else {
		fmt.Println("No River migrations installed yet")
	}
}

func printUsage() {
	fmt.Println("River Migration Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  river-migrate -up       Run all pending River migrations")
	fmt.Println("  river-migrate -down     Rollback the last River migration")
	fmt.Println("  river-migrate -version  Print current River migration version")
	fmt.Println()
	fmt.Println("Environment variables (from .env):")
	fmt.Println("  DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSL_MODE")
	os.Exit(0)
}
