package dbmodel

import "testing"

func TestNewColumnReference(t *testing.T) {
	c1 := &Column{tableName: "users", name: "id"}
	c2 := &Column{tableName: "posts", name: "user_id"}
	cr := NewColumnReference(c2, c1)
	if expected, actual := c2, cr.From(); actual != expected {
		t.Errorf("From() returns invalid value. expected: %+v, actual: %+v", expected, actual)
	}
	if expected, actual := c1, cr.To(); actual != expected {
		t.Errorf("To() returns invalid value. expected: %+v, actual: %+v", expected, actual)
	}
}
