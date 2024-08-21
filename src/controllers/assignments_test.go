package controllers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAssignBike_Success(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create a mock user and bike
	userUUID := "user-uuid-1"
	bikeID := "bike-uuid-1"

	// Prepare mock expectations
	mock.ExpectQuery("SELECT id, role FROM users WHERE id = \\$1").
		WithArgs(userUUID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "role"}).AddRow(userUUID, "Customer"))

	mock.ExpectQuery("SELECT id FROM assignments WHERE user_id = \\$1 AND unassigned_at IS NULL").
		WithArgs(userUUID).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery("SELECT id, is_assigned, usage_count, last_unassigned FROM bikes WHERE is_assigned = false AND \\(last_unassigned IS NULL OR last_unassigned < \\$1\\) ORDER BY usage_count ASC LIMIT 1").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "is_assigned", "usage_count", "last_unassigned"}).AddRow(bikeID, false, 0, time.Now().Add(-10*time.Minute)))

	mock.ExpectExec("UPDATE bikes SET is_assigned = true, usage_count = \\$1 WHERE id = \\$2").
		WithArgs(1, bikeID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO assignments \\(user_id, bike_id, assigned_at\\) VALUES \\(\\$1, \\$2, \\$3\\)").
		WithArgs(userUUID, bikeID, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create a new HTTP request
	reqBody := `{"user_uuid":"user-uuid-1"}`
	req, err := http.NewRequest(http.MethodPost, "/assign-bike", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the function to test
	AssignBike(rr, req, db)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK but got %v", rr.Code)

	// Check the response body
	assert.Equal(t, "Bike assigned successfully", rr.Body.String())

	// Assert that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestAssignBike_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Prepare mock expectations for user not found
	userUUID := "user-uuid-1"
	mock.ExpectQuery("SELECT id, role FROM users WHERE id = \\$1").
		WithArgs(userUUID).
		WillReturnError(sql.ErrNoRows)

	// Create a new HTTP request
	reqBody := `{"user_uuid":"user-uuid-1"}`
	req, err := http.NewRequest(http.MethodPost, "/assign-bike", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the function to test
	AssignBike(rr, req, db)

	// Check the status code
	assert.Equal(t, http.StatusNotFound, rr.Code, "Expected status Not Found but got %v", rr.Code)

	// Check the response body
	assert.Equal(t, "User not found\n", rr.Body.String())

	// Assert that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestAssignBike_ActiveAssignmentExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Prepare mock expectations for active assignment
	userUUID := "user-uuid-1"
	mock.ExpectQuery("SELECT id, role FROM users WHERE id = \\$1").
		WithArgs(userUUID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "role"}).AddRow(userUUID, "Customer"))

	mock.ExpectQuery("SELECT id FROM assignments WHERE user_id = \\$1 AND unassigned_at IS NULL").
		WithArgs(userUUID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Create a new HTTP request
	reqBody := `{"user_uuid":"user-uuid-1"}`
	req, err := http.NewRequest(http.MethodPost, "/assign-bike", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the function to test
	AssignBike(rr, req, db)

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status Bad Request but got %v", rr.Code)

	// Check the response body
	assert.Equal(t, "User already has an active bike assignment\n", rr.Body.String())

	// Assert that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestAssignBike_NoAvailableBikes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Prepare mock expectations for no available bikes
	userUUID := "user-uuid-1"
	mock.ExpectQuery("SELECT id, role FROM users WHERE id = \\$1").
		WithArgs(userUUID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "role"}).AddRow(userUUID, "Customer"))

	mock.ExpectQuery("SELECT id FROM assignments WHERE user_id = \\$1 AND unassigned_at IS NULL").
		WithArgs(userUUID).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery("SELECT id, is_assigned, usage_count, last_unassigned FROM bikes WHERE is_assigned = false AND \\(last_unassigned IS NULL OR last_unassigned < \\$1\\) ORDER BY usage_count ASC LIMIT 1").
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(sql.ErrNoRows)

	// Create a new HTTP request
	reqBody := `{"user_uuid":"user-uuid-1"}`
	req, err := http.NewRequest(http.MethodPost, "/assign-bike", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the function to test
	AssignBike(rr, req, db)

	// Check the status code
	assert.Equal(t, http.StatusNotFound, rr.Code, "Expected status Not Found but got %v", rr.Code)

	// Check the response body
	assert.Equal(t, "No available bikes\n", rr.Body.String())

	// Assert that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
