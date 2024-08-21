package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/yourusername/bike-rental/src/database/models"
)

// Assume models package is properly defined
type AssignBikeRequest struct {
	UserUUID string `json:"user_uuid"`
}

func AssignBike(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Parse the JSON request body
	var req AssignBikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch the user based on UUID
	var user models.User
	query := "SELECT id, role FROM users WHERE id = $1"
	if err := db.QueryRow(query, req.UserUUID).Scan(&user.ID, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		}
		return
	}

	// Check if the user is an Admin
	if user.Role == "Admin" {
		http.Error(w, "Admins cannot be assigned bikes", http.StatusBadRequest)
		return
	}

	// Check if the user already has an active bike assignment
	var existingAssignment models.Assignment
	query = "SELECT id FROM assignments WHERE user_id = $1 AND unassigned_at IS NULL"
	if err := db.QueryRow(query, user.ID).Scan(&existingAssignment.ID); err == nil {
		http.Error(w, "User already has an active bike assignment", http.StatusBadRequest)
		return
	} else if err != sql.ErrNoRows {
		http.Error(w, "Failed to check user assignments", http.StatusInternalServerError)
		return
	}

	// Fetch the least used bike that is not assigned and was unassigned more than 5 minutes ago
	var bike models.Bike
	query = `SELECT id, is_assigned, usage_count, last_unassigned
	         FROM bikes 
	         WHERE is_assigned = false 
	         AND (last_unassigned IS NULL OR last_unassigned < $1)
	         ORDER BY usage_count ASC
	         LIMIT 1`
	if err := db.QueryRow(query, time.Now().Add(-5*time.Minute)).Scan(&bike.ID, &bike.IsAssigned, &bike.UsageCount, &bike.LastUnassigned); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No available bikes", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch bike", http.StatusInternalServerError)
		}
		return
	}

	// Update bike status and usage count
	bike.UsageCount++
	query = "UPDATE bikes SET is_assigned = true, usage_count = $1 WHERE id = $2"
	if _, err := db.Exec(query, bike.UsageCount, bike.ID); err != nil {
		http.Error(w, "Failed to assign bike", http.StatusInternalServerError)
		return
	}

	// Create a new assignment record
	query = `INSERT INTO assignments (user_id, bike_id, assigned_at)
	         VALUES ($1, $2, $3)`
	if _, err := db.Exec(query, user.ID, bike.ID, time.Now()); err != nil {
		http.Error(w, "Failed to create assignment", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Bike assigned successfully"))
}

type UnassignBikeRequest struct {
	BikeUUID string `json:"bike_uuid"`
	UserUUID string `json:"user_uuid"`
}

func UnassignBike(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Parse the JSON request body
	var req UnassignBikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate that both BikeUUID and UserUUID are provided
	if req.BikeUUID == "" || req.UserUUID == "" {
		http.Error(w, "Both bike_uuid and user_uuid are required", http.StatusBadRequest)
		return
	}

	// Fetch the bike and its assignment based on UUID and user ID
	var bikeID string
	query := `
		SELECT b.id
		FROM bikes b
		INNER JOIN assignments a ON b.id = a.bike_id
		WHERE b.id = $1 AND a.user_id = $2 AND a.unassigned_at IS NULL
	`
	if err := db.QueryRow(query, req.BikeUUID, req.UserUUID).Scan(&bikeID); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Bike not found or not assigned to the user", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch bike assignment", http.StatusInternalServerError)
		}
		return
	}

	// Update bike to be unassigned and set the last_unassigned timestamp
	now := time.Now()
	query = "UPDATE bikes SET is_assigned = false, last_unassigned = $1 WHERE id = $2"
	if _, err := db.Exec(query, now, bikeID); err != nil {
		http.Error(w, "Failed to unassign bike", http.StatusInternalServerError)
		return
	}

	// Update the corresponding assignment record to set the unassigned_at timestamp
	query = "UPDATE assignments SET unassigned_at = $1 WHERE bike_id = $2 AND unassigned_at IS NULL"
	if _, err := db.Exec(query, now, bikeID); err != nil {
		http.Error(w, "Failed to update assignment record", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Bike unassigned successfully"))
}

// GetAllAssignments retrieves all assignments from the database using database/sql
func GetAllAssignments(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Prepare the query
	query := "SELECT id, user_id, bike_id, assigned_at, unassigned_at FROM assignments"

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed to retrieve assignments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate over the rows and build the assignments slice
	var assignments []models.Assignment
	for rows.Next() {
		var assignment models.Assignment
		if err := rows.Scan(&assignment.ID, &assignment.UserID, &assignment.BikeID, &assignment.AssignedAt, &assignment.UnassignedAt); err != nil {
			http.Error(w, "Failed to scan assignment", http.StatusInternalServerError)
			return
		}
		assignments = append(assignments, assignment)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		http.Error(w, "Error encountered during row iteration", http.StatusInternalServerError)
		return
	}

	// Respond with the list of assignments in JSON format
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(assignments); err != nil {
		http.Error(w, "Failed to encode assignments to JSON", http.StatusInternalServerError)
		return
	}
}
