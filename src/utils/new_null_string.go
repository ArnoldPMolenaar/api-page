package utils

import "database/sql"

// NewNullString creates a sql.NullString from a string pointer.
func NewNullString(value *string) sql.NullString {
	if value != nil && *value != "" {
		return sql.NullString{Valid: true, String: *value}
	}
	return sql.NullString{Valid: false}
}
