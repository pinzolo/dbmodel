package dbmodel

import (
	"strconv"
	"testing"
)

func TestAddColumn(t *testing.T) {
	tbl := newUserTable()
	if len(tbl.Columns()) != 0 {
		t.Error("If table has no column, Columns() should be zero length.")
	}
	col := Column{name: "name"}
	tbl.AddColumn(&col)
	if len(tbl.Columns()) != 1 {
		t.Errorf("If table has a column, Columns() should be 1 length. (%#v)", tbl)
	}
	if tbl.Columns()[0].Name() != "name" {
		t.Errorf("Invalid column added. (%#v)", tbl.Columns())
	}
	for i := 0; i < 10; i++ {
		c := Column{name: "name" + strconv.Itoa(i)}
		tbl.AddColumn(&c)
	}
	if len(tbl.Columns()) != 11 {
		t.Errorf("If table has some columns, Columns() should be valid length. (%#v)", tbl)
	}
}

func newUserTable() *Table {
	table := NewTable("foo", "users", "")
	return &table
}
