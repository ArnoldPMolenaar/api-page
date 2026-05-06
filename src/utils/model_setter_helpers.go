package utils

import (
	"database/sql"
	"time"
)

// PtrFromUint creates a uint pointer from a value.
func PtrFromUint(v uint) *uint {
	vv := v
	return &vv
}

// PtrFromNullInt16 creates an int16 pointer from sql.NullInt16.
func PtrFromNullInt16(v sql.NullInt16) *int16 {
	if !v.Valid {
		return nil
	}

	vv := v.Int16
	return &vv
}

// PtrFromNullUint creates a uint pointer from sql.Null[uint].
func PtrFromNullUint(v sql.Null[uint]) *uint {
	if !v.Valid {
		return nil
	}

	vv := v.V
	return &vv
}

// PtrFromNullUint8 creates a uint8 pointer from sql.Null[uint8].
func PtrFromNullUint8(v sql.Null[uint8]) *uint8 {
	if !v.Valid {
		return nil
	}

	vv := v.V
	return &vv
}

// PtrFromNullString creates a string pointer from sql.NullString.
func PtrFromNullString(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}

	vv := v.String
	return &vv
}

// PtrFromNullTime creates a time pointer from sql.NullTime.
func PtrFromNullTime(v sql.NullTime) *time.Time {
	if !v.Valid {
		return nil
	}

	vv := v.Time
	return &vv
}
