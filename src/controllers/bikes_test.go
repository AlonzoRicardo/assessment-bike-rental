package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/test-go/testify/assert"
	"github.com/yourusername/bike-rental/src/database/models"
)

func TestGetAvailableBikes(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Use a fixed time for testing
	fixedTime := time.Date(2024, 8, 21, 7, 33, 52, 0, time.UTC)
	timeNow = func() time.Time {
		return fixedTime
	}
	defer func() { timeNow = time.Now }() // Restore the original timeNow after the test

	gracePeriod := fixedTime.Add(-5 * time.Minute)

	// Prepare mock data
	mockRows := sqlmock.NewRows([]string{"id", "is_assigned", "usage_count", "last_unassigned"}).
		AddRow("bike-1", false, 10, sql.NullTime{Time: gracePeriod.Add(-10 * time.Minute), Valid: true}).
		AddRow("bike-2", false, 5, sql.NullTime{Time: gracePeriod.Add(-15 * time.Minute), Valid: true})

	// Set up the expectations
	mock.ExpectQuery(`SELECT id, is_assigned, usage_count, last_unassigned FROM bikes WHERE is_assigned = false AND \(last_unassigned IS NULL OR last_unassigned < \$1\)`).
		WithArgs(gracePeriod).
		WillReturnRows(mockRows)

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodGet, "/bikes/available", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the function
	GetAvailableBikes(rr, req, db)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK but got %v", rr.Code)

	// Check if the body contains valid JSON
	var bikes []models.Bike
	err = json.NewDecoder(rr.Body).Decode(&bikes)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v\nResponse body: %v", err, rr.Body.String())
	}

	// Assert the response data
	assert.Len(t, bikes, 2, "Expected 2 bikes but got %v", len(bikes))
	assert.Equal(t, "bike-1", bikes[0].ID)
	assert.Equal(t, "bike-2", bikes[1].ID)

	// Assert that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAllBikes(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Prepare mock data
	mockRows := sqlmock.NewRows([]string{"id", "is_assigned", "usage_count", "last_unassigned"}).
		AddRow("bike-1", false, 10, sql.NullTime{Time: time.Now(), Valid: true}).
		AddRow("bike-2", true, 5, sql.NullTime{Time: time.Now(), Valid: true})

	// Set up the expectations for the SELECT query
	mock.ExpectQuery("SELECT id, is_assigned, usage_count, last_unassigned FROM bikes").
		WillReturnRows(mockRows)

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodGet, "/bikes", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the function to test
	GetAllBikes(rr, req, db)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK but got %v", rr.Code)

	// Parse the response
	var bikes []models.Bike
	err = json.NewDecoder(rr.Body).Decode(&bikes)
	assert.NoError(t, err, "Failed to decode response body: %v\nResponse body: %v", err, rr.Body.String())

	// Assert the response data
	assert.Len(t, bikes, 2, "Expected 2 bikes but got %v", len(bikes))
	assert.Equal(t, "bike-1", bikes[0].ID)
	assert.Equal(t, "bike-2", bikes[1].ID)

	// Assert that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAllBikes_DBError(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Set up the expectations for the SELECT query to return an error
	mock.ExpectQuery("SELECT id, is_assigned, usage_count, last_unassigned FROM bikes").
		WillReturnError(sql.ErrConnDone)

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodGet, "/bikes", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the function to test
	GetAllBikes(rr, req, db)

	// Check the status code
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status Internal Server Error but got %v", rr.Code)

	// Assert the response body contains the expected error message
	assert.Equal(t, "Failed to retrieve bikes\n", rr.Body.String())

	// Assert that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
