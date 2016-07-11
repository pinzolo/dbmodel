package dbmodel

import (
	"database/sql"
	"fmt"
	"strconv"
)

// Size is column size
type Size struct {
	length    sql.NullInt64
	precision sql.NullInt64
	scale     sql.NullInt64
	valid     bool
}

// Length returns length.(used when data type is charactor)
func (s Size) Length() sql.NullInt64 {
	return s.length
}

// Precision returns precision.(used when data type is numeric)
func (s Size) Precision() sql.NullInt64 {
	return s.precision
}

// Scale returns scale.(used when data type is numeric)
func (s Size) Scale() sql.NullInt64 {
	return s.scale
}

// IsValid returns has value.
// If length, precision and scale are NULL, IsValid returns false.
func (s Size) IsValid() bool {
	return s.valid
}

// String returns string expression.
func (s Size) String() string {
	if s.IsValid() {
		if s.Length().Valid {
			return strconv.FormatInt(s.Length().Int64, 10)
		} else if s.Scale().Valid {
			return fmt.Sprintf("%v, %v", s.Precision().Int64, s.Scale().Int64)
		} else {
			return strconv.FormatInt(s.Precision().Int64, 10)
		}
	} else {
		return ""
	}
}

// NewSize returns new Size initialized with arguments.
func NewSize(length sql.NullInt64, precision sql.NullInt64, scale sql.NullInt64) Size {
	return Size{
		length:    length,
		precision: precision,
		scale:     scale,
		valid:     length.Valid || precision.Valid || scale.Valid,
	}
}
