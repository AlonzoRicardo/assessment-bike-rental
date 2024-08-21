package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/yourusername/bike-rental/src/database/models"
)

var timeNow = time.Now

func GetAvailableBikes(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	gracePeriod := timeNow().Add(-5 * time.Minute)

	// Prepare the query
	query := `SELECT id, is_assigned, usage_count, last_unassigned 
	          FROM bikes 
	          WHERE is_assigned = false 
	          AND (last_unassigned IS NULL OR last_unassigned < $1)`

	// Execute the query
	rows, err := db.Query(query, gracePeriod)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to retrieve available bikes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Initialize the bikes slice
	bikes := []models.Bike{}

	// Iterate over the rows and build the bikes slice
	for rows.Next() {
		var bike models.Bike
		if err := rows.Scan(&bike.ID, &bike.IsAssigned, &bike.UsageCount, &bike.LastUnassigned); err != nil {
			log.Printf("Scan error: %v", err)
			http.Error(w, "Failed to scan bike", http.StatusInternalServerError)
			return
		}
		bikes = append(bikes, bike)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		http.Error(w, "Error encountered during row iteration", http.StatusInternalServerError)
		return
	}

	// Respond with the list of available bikes in JSON format
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bikes); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, "Failed to encode bikes to JSON", http.StatusInternalServerError)
		return
	}
}

func GetAllBikes(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Prepare the query
	query := "SELECT id, is_assigned, usage_count, last_unassigned FROM bikes"

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed to retrieve bikes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Initialize the bikes slice
	bikes := []models.Bike{}

	// Iterate over the rows and build the bikes slice
	for rows.Next() {
		var bike models.Bike
		if err := rows.Scan(&bike.ID, &bike.IsAssigned, &bike.UsageCount, &bike.LastUnassigned); err != nil {
			http.Error(w, "Failed to scan bike", http.StatusInternalServerError)
			return
		}
		bikes = append(bikes, bike)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		http.Error(w, "Error encountered during row iteration", http.StatusInternalServerError)
		return
	}

	// Respond with the list of bikes in JSON format
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bikes); err != nil {
		http.Error(w, "Failed to encode bikes to JSON", http.StatusInternalServerError)
		return
	}
}
