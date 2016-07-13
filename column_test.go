package dbmodel

import "testing"

func TestNewColumn(t *testing.T) {
	c := NewColumn("name", "TEXT", NewSize(validInt(9), invalidInt(), invalidInt()), true, "Jone Doe")
	if c.Name() != "name" {
		t.Errorf("Name() returns invalid value. actual is %v", c.Name())
	}
	if c.DataType() != "TEXT" {
		t.Errorf("DataType() returns invalid value. actual is %v", c.DataType())
	}
	if c.Size().String() != "9" {
		t.Errorf("Size() returns invalid value. actual is %v", c.Size())
	}
	if !c.IsNullable() {
		t.Error("Geven 'YES', IsNullable should return true.")
	}
	if c.DefaultValue() != "Jone Doe" {
		t.Errorf("DefaultValue() returns invalid value. actual is %v", c.DefaultValue())
	}
}
