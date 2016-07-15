package dbmodel

import "testing"

func TestNewColumn(t *testing.T) {
	c := NewColumn("name", "comment", "TEXT", NewSize(validInt(9), invalidInt(), invalidInt()), true, "Jone Doe")
	if c.Name() != "name" {
		t.Errorf("Name() returns invalid value. expected: %v, actual: %v", "name", c.Name())
	}
	if c.Comment() != "comment" {
		t.Errorf("Comment() returns invalid value. expected: %v, actual: %v", "comment", c.Comment())
	}
	if c.DataType() != "TEXT" {
		t.Errorf("DataType() returns invalid value. expected: %v, actual: %v", "TEXT", c.DataType())
	}
	if c.Size().String() != "9" {
		t.Errorf("Size() returns invalid value. expected: %v, actual: %v", "9", c.Size())
	}
	if !c.IsNullable() {
		t.Error("Geven 'YES', IsNullable should return true.")
	}
	if c.DefaultValue() != "Jone Doe" {
		t.Errorf("DefaultValue() returns invalid value. expected: %v, actual: %v", "Jone Doe", c.DefaultValue())
	}
}
