package models

import (
	"database/sql"
)

type Bike struct {
	ID             string       `json:"id"`
	UsageCount     int          `json:"usage_count"`
	LastUnassigned sql.NullTime `json:"last_unassigned"`
	IsAssigned     bool         `json:"is_assigned"`
}
