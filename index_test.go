package dbmodel

import (
	"strconv"
	"testing"
)

func TestNewIndex(t *testing.T) {
	idx := NewIndex("foo", "users", "users_pk", true)
	if expected, actual := "foo", idx.Schema(); actual != expected {
		t.Errorf("Schema() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := "users", idx.TableName(); actual != expected {
		t.Errorf("TableName() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := "users_pk", idx.Name(); actual != expected {
		t.Errorf("Name() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := true, idx.IsUnique(); actual != expected {
		t.Errorf("IsUnique() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
}

func TestAddColumnToIndex(t *testing.T) {
	idx := NewIndex("foo", "users", "users_pk", true)
	if idx.Columns() == nil || len(idx.Columns()) != 0 {
		t.Error("Columns() should be initialized.")
	}
	col := Column{name: "name"}
	idx.AddColumn(&col)
	if len(idx.Columns()) != 1 {
		t.Errorf("If table has a column, Columns() should be 1 length. (%+v)", idx)
	}
	if idx.Columns()[0].Name() != "name" {
		t.Errorf("Invalid column added. (%+v)", idx.Columns())
	}
	for i := 0; i < 10; i++ {
		c := Column{name: "name" + strconv.Itoa(i)}
		idx.AddColumn(&c)
	}
	if len(idx.Columns()) != 11 {
		t.Errorf("If table has some columns, Columns() should be valid length. (%+v)", idx)
	}
}
