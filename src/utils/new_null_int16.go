package utils

import "database/sql"

// NewNullInt16 creates a sql.NullInt16 from a int16 pointer.
func NewNullInt16(value *int16) sql.NullInt16 {
	if value != nil {
		return sql.NullInt16{Valid: true, Int16: *value}
	}
	return sql.NullInt16{Valid: false}
}
