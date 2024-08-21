package cronjobs

import (
	"database/sql"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yourusername/bike-rental/src/database/models"
)

// timeNow is a variable that returns the current time. It can be overridden in tests.
var timeNow = time.Now

func AutoUnassignOverdueBikes(db *sql.DB) {
	// Calculate the cutoff time for 24 hours ago
	cutoff := timeNow().Add(-24 * time.Hour)

	log.Debug().Msg("Scanning for overdue bike assignments...")

	// Find all assignments older than 24 hours and still active (unassigned_at is NULL)
	query := `SELECT id, user_id, bike_id FROM assignments WHERE assigned_at < $1 AND unassigned_at IS NULL`
	rows, err := db.Query(query, cutoff)
	if err != nil {
		log.Err(err).Msg("Failed to retrieve overdue assignments")
		return
	}
	defer rows.Close()

	var overdueAssignments []models.Assignment
	for rows.Next() {
		var assignment models.Assignment
		if err := rows.Scan(&assignment.ID, &assignment.UserID, &assignment.BikeID); err != nil {
			log.Err(err).Msg("Failed to scan overdue assignment")
			return
		}
		overdueAssignments = append(overdueAssignments, assignment)
	}

	// Process each overdue assignment
	for _, assignment := range overdueAssignments {
		// Unassign the bike and update the assignment
		if err := unassignBikeByUserID(db, assignment.UserID); err != nil {
			log.Err(err).Str("user", assignment.UserID).Msg("Failed to unassign bike")
		}
	}
}

func unassignBikeByUserID(db *sql.DB, userID string) error {
	// Find the active assignment for the user
	query := `SELECT id, bike_id FROM assignments WHERE user_id = $1 AND unassigned_at IS NULL LIMIT 1`
	var assignment models.Assignment
	if err := db.QueryRow(query, userID).Scan(&assignment.ID, &assignment.BikeID); err != nil {
		if err == sql.ErrNoRows {
			return nil // No active assignment found, nothing to do
		}
		return err
	}

	log.Info().Uint("assignment", uint(assignment.ID)).Msg("Found overdue bike assignment...")

	// Mark the bike as unassigned
	query = `UPDATE bikes SET is_assigned = false, last_unassigned = $1 WHERE id = $2`
	now := time.Now()
	if _, err := db.Exec(query, now, assignment.BikeID); err != nil {
		return err
	}

	// Update the assignment to mark it as unassigned
	query = `UPDATE assignments SET unassigned_at = $1 WHERE id = $2`
	if _, err := db.Exec(query, now, assignment.ID); err != nil {
		return err
	}

	log.Info().Uint("assignment", uint(assignment.ID)).Msg("Successfully unassigned overdue bike")

	return nil
}
