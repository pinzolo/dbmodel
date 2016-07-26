package dbmodel

import "testing"

func TestNewConstraint(t *testing.T) {
	c := NewConstraint("foo", "users", "users_age_check", "CHECK", "(age >= 0)")
	if expected, actual := "foo", c.Schema(); actual != expected {
		t.Errorf("Schema() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := "users", c.TableName(); actual != expected {
		t.Errorf("TableName() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := "users_age_check", c.Name(); actual != expected {
		t.Errorf("Name() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := "CHECK", c.Kind(); actual != expected {
		t.Errorf("Kind() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
	if expected, actual := "(age >= 0)", c.Content(); actual != expected {
		t.Errorf("Content() returns invalid value. expected: %v, actual: %v", expected, actual)
	}
}
