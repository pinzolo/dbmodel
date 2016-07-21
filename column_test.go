package dbmodel

import "testing"

func TestNewColumn(t *testing.T) {
	c := NewColumn("foo", "users", "name", "comment", "TEXT", NewSize(validInt(9), invalidInt(), invalidInt()), true, "Jone Doe", int64(1))
	if expected, actual := "foo", c.Schema(); actual != expected {
		t.Errorf("Schema() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := "users", c.TableName(); actual != expected {
		t.Errorf("TableName() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := "name", c.Name(); actual != expected {
		t.Errorf("Name() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := "comment", c.Comment(); actual != expected {
		t.Errorf("Comment() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := "TEXT", c.DataType(); actual != expected {
		t.Errorf("DataType() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := "9", c.Size().String(); actual != expected {
		t.Errorf("Size() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
	if !c.IsNullable() {
		t.Error("Geven 'YES', IsNullable should return true.")
	}
	if expected, actual := "Jone Doe", c.DefaultValue(); actual != expected {
		t.Errorf("DefaultValue() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := int64(1), c.PrimaryKeyPosition(); actual != expected {
		t.Errorf("PrimaryKeyPosition() returns invalid value. expected: %v, actual: %v", "Jone Doe", c.DefaultValue())
	}
}
