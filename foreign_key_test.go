package dbmodel

import "testing"

func TestNewForeignKey(t *testing.T) {
	fk := NewForeignKey("foo", "users", "users_pk")
	if expected, actual := "foo", fk.Schema(); actual != expected {
		t.Errorf("Schema() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := "users", fk.TableName(); actual != expected {
		t.Errorf("TableName() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := "users_pk", fk.Name(); actual != expected {
		t.Errorf("Name() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
	if len(fk.ColumnReferences()) != 0 {
		t.Error("ColumnReferences() should be empty when initialized.")
	}
}

func TestAddColumnReference(t *testing.T) {
	fk := NewForeignKey("foo", "users", "users_pk")
	c1 := &Column{tableName: "users", name: "id"}
	c2 := &Column{tableName: "posts", name: "user_id"}
	cr1 := NewColumnReference(c2, c1)
	fk.AddColumnReference(&cr1)
	if len(fk.ColumnReferences()) != 1 {
		t.Error("Failed to add column refrence.")
	}
	if fk.ColumnReferences()[0] != &cr1 {
		t.Error("Failed to add column refrence.")
	}
	c3 := &Column{tableName: "comments", name: "user_id"}
	cr2 := NewColumnReference(c3, c1)
	fk.AddColumnReference(&cr2)
	if len(fk.ColumnReferences()) != 2 {
		t.Error("Failed to add column refrence.")
	}
	if fk.ColumnReferences()[1] != &cr2 {
		t.Error("AddColumnReference shouled add to list tail.")
	}
}
