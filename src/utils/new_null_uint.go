package utils

import "database/sql"

// NewNullUint creates a sql.Null[uint] from an uint pointer.
func NewNullUint(value *uint) sql.Null[uint] {
	var nullUint sql.Null[uint]

	if value != nil && *value != 0 {
		nullUint.Valid = true
		nullUint.V = *value
	} else {
		nullUint.Valid = false
	}

	return nullUint
}
