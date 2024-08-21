package models

import (
	"database/sql"
)

// Assignment represents a record in the assignments table
type Assignment struct {
	ID           uint         `json:"id"`
	UserID       string       `json:"user_id"`
	BikeID       string       `json:"bike_id"`
	AssignedAt   sql.NullTime `json:"assigned_at"`
	UnassignedAt sql.NullTime `json:"unassigned_at"`
}
