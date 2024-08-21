package database

import (
	"database/sql"

	"github.com/rs/zerolog/log"
	"github.com/yourusername/bike-rental/src/database/models"
)

func SeedDatabase(db *sql.DB) {
	// Seed Users
	users := []models.User{
		{
			ID:   "d0ab33d7-8fcc-463d-bade-fefd53b77a96",
			Name: "Alice",
			Role: "Customer",
		},
		{
			ID:   "0b28a7ed-39ef-418f-a0e3-8ad3f794dfc7",
			Name: "Bob",
			Role: "Customer",
		},
		{
			ID:   "da690323-5a78-4d46-a214-943b2ec9d49e",
			Name: "Charlie",
			Role: "Admin",
		},
	}

	for _, user := range users {
		// Check if the user already exists
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", user.ID).Scan(&exists)
		if err != nil {

			log.Err(err).Msg("Failed to check if user exists")
		}

		// Insert the user if they don't already exist
		if !exists {
			_, err := db.Exec("INSERT INTO users (id, name, role) VALUES ($1, $2, $3)", user.ID, user.Name, user.Role)
			if err != nil {
				log.Err(err).Msg("Failed to seed user")
			}
		}
	}

	// Seed Bikes
	bikes := []models.Bike{
		{
			ID:         "331e7ffb-e583-4535-ba41-4c28dc34016d",
			UsageCount: 0,
			IsAssigned: false,
		},
		{
			ID:         "e4ef2d9b-5d5a-4f85-bb3a-b2df8bf42ac1",
			UsageCount: 0,
			IsAssigned: false,
		},
	}

	for _, bike := range bikes {
		// Check if the bike already exists
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM bikes WHERE id = $1)", bike.ID).Scan(&exists)
		if err != nil {
			log.Err(err).Msg("Failed to check if bike exists")

		}

		// Insert the bike if it doesn't already exist
		if !exists {
			_, err := db.Exec("INSERT INTO bikes (id, usage_count, is_assigned) VALUES ($1, $2, $3)", bike.ID, bike.UsageCount, bike.IsAssigned)
			if err != nil {
				log.Err(err).Msg("Failed to seed bike")
			}
		}
	}

	log.Info().Msg("Database seeded successfully")
}
