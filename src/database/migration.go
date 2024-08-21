package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/rs/zerolog/log"
)

func Migrate(config *DatabaseConfig) {
	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode,
	)

	// Connect to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Err(err).Msg("Failed to connect to the database")
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(db, "file://src/database/migrations"); err != nil {
		log.Err(err).Msg("Failed to run migrations")
	}

	log.Info().Msg("Migrations applied successfully")
}

func runMigrations(db *sql.DB, migrationsPath string) error {
	// Create an instance of the Postgres driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// Create a new migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres", driver)
	if err != nil {
		return err
	}

	// Run the migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
