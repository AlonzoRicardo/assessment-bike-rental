package cronjobs

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Set logger to output nothing during tests
	log.Logger = zerolog.New(nil)
}

func TestAutoUnassignOverdueBikes(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Use a fixed time for testing
	fixedTime := time.Date(2024, 8, 20, 7, 19, 48, 208958572, time.UTC)

	// Prepare mock data for overdue assignments
	mockRows := sqlmock.NewRows([]string{"id", "user_id", "bike_id"}).
		AddRow(1, "user-1", "bike-1").
		AddRow(2, "user-2", "bike-2")

	// Set up the expectations for the SELECT query
	mock.ExpectQuery(`SELECT id, user_id, bike_id FROM assignments WHERE assigned_at < .* AND unassigned_at IS NULL`).
		WithArgs(fixedTime.Add(-24 * time.Hour)).
		WillReturnRows(mockRows)

	// Set up expectations for the UPDATE queries
	// mock.ExpectExec(`UPDATE bikes SET is_assigned = false, last_unassigned = ? WHERE id = ?`).
	// 	WithArgs(sqlmock.AnyArg(), "bike-1").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))

	// mock.ExpectExec(`UPDATE assignments SET unassigned_at = ? WHERE id = ?`).
	// 	WithArgs(sqlmock.AnyArg(), 1).
	// 	WillReturnResult(sqlmock.NewResult(1, 1))

	// mock.ExpectExec(`UPDATE bikes SET is_assigned = false, last_unassigned = ? WHERE id = ?`).
	// 	WithArgs(sqlmock.AnyArg(), "bike-2").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))

	// mock.ExpectExec(`UPDATE assignments SET unassigned_at = ? WHERE id = ?`).
	// 	WithArgs(sqlmock.AnyArg(), 2).
	// 	WillReturnResult(sqlmock.NewResult(1, 1))

	// Override the timeNow function to return the fixed time
	timeNow = func() time.Time {
		return fixedTime
	}
	defer func() {
		timeNow = time.Now
	}()

	// Call the function to test
	AutoUnassignOverdueBikes(db)

	// Assert that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestUnassignBikeByUserID(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Prepare mock data
	mockRows := sqlmock.NewRows([]string{"id", "bike_id"}).
		AddRow(1, "bike-1")

	// Set up the expectations
	mock.ExpectQuery(`SELECT id, bike_id FROM assignments WHERE user_id = .* AND unassigned_at IS NULL LIMIT 1`).
		WithArgs("user-1").
		WillReturnRows(mockRows)

	mock.ExpectExec(`UPDATE bikes SET is_assigned = false, last_unassigned = .* WHERE id = .*`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`UPDATE assignments SET unassigned_at = .* WHERE id = .*`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the function
	err = unassignBikeByUserID(db, "user-1")

	// Assert no errors and all expectations were met
	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}
