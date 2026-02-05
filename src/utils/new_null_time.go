package utils

import (
	"database/sql"
	"time"
)

// NewNullTime creates a sql.NullTime from a time pointer.
func NewNullTime(value *time.Time) sql.NullTime {
	if value != nil {
		return sql.NullTime{Valid: true, Time: *value}
	}
	return sql.NullTime{Valid: false}
}
