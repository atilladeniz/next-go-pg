package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/atilladeniz/next-go-pg/backend/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Parse command line flags
	upFlag := flag.Bool("up", false, "Run all up migrations")
	downFlag := flag.Bool("down", false, "Rollback the last migration")
	downAllFlag := flag.Bool("down-all", false, "Rollback all migrations")
	versionFlag := flag.Bool("version", false, "Print current migration version")
	forceFlag := flag.Int("force", -1, "Force set migration version (use with caution)")
	stepsFlag := flag.Int("steps", 0, "Number of migrations to run (positive=up, negative=down)")
	flag.Parse()

	// Load configuration
	cfg := config.Load()

	// Build database URL for golang-migrate
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	// Create migrate instance
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	// Execute command
	switch {
	case *upFlag:
		runUp(m)
	case *downFlag:
		runDown(m, 1)
	case *downAllFlag:
		runDownAll(m)
	case *versionFlag:
		printVersion(m)
	case *forceFlag >= 0:
		forceVersion(m, *forceFlag)
	case *stepsFlag != 0:
		runSteps(m, *stepsFlag)
	default:
		printUsage()
	}
}

func runUp(m *migrate.Migrate) {
	fmt.Println("Running all up migrations...")
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to apply - database is up to date")
			return
		}
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("Migrations applied successfully")
}

func runDown(m *migrate.Migrate, steps int) {
	fmt.Printf("Rolling back %d migration(s)...\n", steps)
	if err := m.Steps(-steps); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to rollback")
			return
		}
		log.Fatalf("Rollback failed: %v", err)
	}
	fmt.Println("Rollback completed successfully")
}

func runDownAll(m *migrate.Migrate) {
	fmt.Println("Rolling back all migrations...")
	if err := m.Down(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to rollback")
			return
		}
		log.Fatalf("Rollback failed: %v", err)
	}
	fmt.Println("All migrations rolled back successfully")
}

func runSteps(m *migrate.Migrate, steps int) {
	direction := "up"
	if steps < 0 {
		direction = "down"
	}
	fmt.Printf("Running %d %s migration(s)...\n", abs(steps), direction)
	if err := m.Steps(steps); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to apply")
			return
		}
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("Migrations applied successfully")
}

func printVersion(m *migrate.Migrate) {
	version, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			fmt.Println("No migrations have been applied yet")
			return
		}
		log.Fatalf("Failed to get version: %v", err)
	}
	fmt.Printf("Current version: %d\n", version)
	if dirty {
		fmt.Println("Warning: Database is in dirty state!")
	}
}

func forceVersion(m *migrate.Migrate, version int) {
	fmt.Printf("Forcing migration version to %d...\n", version)
	if err := m.Force(version); err != nil {
		log.Fatalf("Failed to force version: %v", err)
	}
	fmt.Println("Version forced successfully")
}

func printUsage() {
	fmt.Println("Database Migration Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  migrate -up           Run all pending migrations")
	fmt.Println("  migrate -down         Rollback the last migration")
	fmt.Println("  migrate -down-all     Rollback all migrations")
	fmt.Println("  migrate -version      Print current migration version")
	fmt.Println("  migrate -steps N      Run N migrations (positive=up, negative=down)")
	fmt.Println("  migrate -force N      Force set version to N (use with caution)")
	fmt.Println()
	fmt.Println("Environment variables (from .env):")
	fmt.Println("  DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSL_MODE")
	os.Exit(0)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
