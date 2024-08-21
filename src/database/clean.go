package database

import (
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
)

// CleanDatabase deletes all records from the database tables
func CleanDatabase(db *sql.DB) {
	tables := []string{
		"assignments",
		"bikes",
		"users",
	}

	for _, table := range tables {
		query := fmt.Sprintf("DELETE FROM %s", table)
		if _, err := db.Exec(query); err != nil {
			log.Err(err).Msg("Failed to clean table")
		}
	}

	log.Info().Msg("All records deleted successfully")
}
