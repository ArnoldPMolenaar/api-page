package utils

import "database/sql"

// NewNullUint8 creates a sql.NullUint8 from an uint8 pointer.
func NewNullUint8(value *uint8) sql.Null[uint8] {
	if value != nil {
		return sql.Null[uint8]{Valid: true, V: *value}
	}
	return sql.Null[uint8]{Valid: false}
}
