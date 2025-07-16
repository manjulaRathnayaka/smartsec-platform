package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// Initialize creates and configures the database connection
func Initialize(databaseURL string) (*sql.DB, error) {
	log.Debug().Msg("Opening database connection")
	
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool for production use
	log.Debug().Msg("Configuring database connection pool")
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	// Test connection with retry logic
	maxRetries := 5
	log.Info().Int("max_retries", maxRetries).Msg("Testing database connection")
	
	for i := 0; i < maxRetries; i++ {
		if err := db.Ping(); err == nil {
			log.Info().Int("retry_attempt", i+1).Msg("Database connection successful")
			break
		} else {
			log.Warn().
				Err(err).
				Int("retry_attempt", i+1).
				Int("max_retries", maxRetries).
				Msg("Database connection failed, retrying...")
		}
		if i == maxRetries-1 {
			return nil, fmt.Errorf("failed to ping database after %d retries: %w", maxRetries, err)
		}
		sleepDuration := time.Duration(i+1) * time.Second
		log.Debug().
			Dur("sleep_duration", sleepDuration).
			Msg("Sleeping before next retry")
		time.Sleep(sleepDuration)
	}

	// Test basic database operations
	log.Debug().Msg("Testing basic database operations")
	if err := testDatabaseOperations(db); err != nil {
		log.Warn().Err(err).Msg("Database operations test failed")
		return nil, fmt.Errorf("database operations test failed: %w", err)
	}

	return db, nil
}

// RunMigrations runs database migrations
func RunMigrations(databaseURL string) error {
	log.Info().Msg("Starting database migrations")
	
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database for migrations: %w", err)
	}
	defer db.Close()

	log.Debug().Msg("Creating migration driver")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	log.Debug().Msg("Creating migration instance")
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	log.Info().Msg("Running migrations")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Info().Msg("No new migrations to apply")
	} else {
		log.Info().Msg("Database migrations completed successfully")
	}

	return nil
}

// testDatabaseOperations performs basic database operations to verify connectivity
func testDatabaseOperations(db *sql.DB) error {
	log.Debug().Msg("Testing SELECT 1 query")
	
	// Test basic SELECT query
	var result int
	if err := db.QueryRow("SELECT 1").Scan(&result); err != nil {
		return fmt.Errorf("failed to execute SELECT 1: %w", err)
	}
	
	if result != 1 {
		return fmt.Errorf("unexpected result from SELECT 1: got %d, expected 1", result)
	}
	
	log.Debug().Msg("Basic database operations test passed")
	return nil
}
