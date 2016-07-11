package dbmodel

import (
	"database/sql"
	"testing"
)

func TestLength(t *testing.T) {
	l := NewSize(validInt(9), invalidInt(), invalidInt()).Length()
	if !l.Valid {
		t.Errorf("Length() returns invalid value. actual is %v", l)
	}
	if l.Int64 != 9 {
		t.Errorf("Length() returns invalid value. actual is %v", l)
	}
}

func TestPrecision(t *testing.T) {
	p := NewSize(invalidInt(), validInt(6), invalidInt()).Precision()
	if !p.Valid {
		t.Errorf("Precision() returns invalid value. actual is %v", p)
	}
	if p.Int64 != 6 {
		t.Errorf("Precision() returns invalid value. actual is %v", p)
	}
}

func TestScale(t *testing.T) {
	s := NewSize(invalidInt(), validInt(6), validInt(3)).Scale()
	if !s.Valid {
		t.Errorf("Scale() returns invalid value. actual is %v", s)
	}
	if s.Int64 != 3 {
		t.Errorf("Scale() returns invalid value. actual is %v", s)
	}
}

func TestIsValidOnInvalid(t *testing.T) {
	s := NewSize(invalidInt(), invalidInt(), invalidInt())
	if s.IsValid() {
		t.Error("If all values are NULL, IsValid() should return false")
	}
}

func TestIsValidOnValidLength(t *testing.T) {
	s := NewSize(validInt(9), invalidInt(), invalidInt())
	if !s.IsValid() {
		t.Error("On having valid length, IsValid should return true")
	}
}

func TestIsValidOnValidPrecisionAndInvalidScale(t *testing.T) {
	s := NewSize(invalidInt(), validInt(6), invalidInt())
	if !s.IsValid() {
		t.Error("On having valid precision, IsValid should return true")
	}
}

func TestIsValidOnValidPrecisionAndValidScale(t *testing.T) {
	s := NewSize(invalidInt(), validInt(6), validInt(3))
	if !s.IsValid() {
		t.Error("On having valid precision and scale, IsValid should return true")
	}
}

func TestStringOnInvalid(t *testing.T) {
	s := NewSize(invalidInt(), invalidInt(), invalidInt())
	if s.String() != "" {
		t.Errorf("If all values are NULL, String() should return empty string. (expected: %v, actual: %v)", "", s.String())
	}
}

func TestStringOnValidLength(t *testing.T) {
	s := NewSize(validInt(9), invalidInt(), invalidInt())
	if s.String() != "9" {
		t.Errorf("On having valid length, String() should return its length as string. (expected: %v, actual: %v)", "9", s.String())
	}
}

func TestStringOnValidPrecisionAndInvalidScale(t *testing.T) {
	s := NewSize(invalidInt(), validInt(6), invalidInt())
	if s.String() != "6" {
		t.Errorf("On having valid prescision, String() should return its precision as string. (expected: %v, actual: %v)", "6", s.String())
	}
}

func TestStringOnValidPrecisionAndValidScale(t *testing.T) {
	s := NewSize(invalidInt(), validInt(6), validInt(3))
	if s.String() != "6, 3" {
		t.Errorf("On having valid prescision and scale, String() should return its precision and scale as comma separated string. (expected: %v, actual: %v)", "6, 3", s.String())
	}
}

func validInt(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: true}
}

func invalidInt() sql.NullInt64 {
	return sql.NullInt64{Valid: false}
}
