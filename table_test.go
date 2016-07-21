package dbmodel

import (
	"strconv"
	"testing"
)

func TestAddColumnToTable(t *testing.T) {
	tbl := newUserTable()
	if len(tbl.Columns()) != 0 {
		t.Error("If table has no column, Columns() should be zero length.")
	}
	col := Column{name: "name"}
	tbl.AddColumn(&col)
	if len(tbl.Columns()) != 1 {
		t.Errorf("If table has a column, Columns() should be 1 length. (%+v)", tbl)
	}
	if tbl.Columns()[0].Name() != "name" {
		t.Errorf("Invalid column added. (%+v)", tbl.Columns())
	}
	if tbl.Columns()[0].Schema() != tbl.Schema() {
		t.Errorf("Column's schema should be set by table's schema.")
	}
	if tbl.Columns()[0].TableName() != tbl.Name() {
		t.Errorf("Column's table name should be set by table's name.")
	}
	for i := 0; i < 10; i++ {
		c := Column{name: "name" + strconv.Itoa(i)}
		tbl.AddColumn(&c)
	}
	if len(tbl.Columns()) != 11 {
		t.Errorf("If table has some columns, Columns() should be valid length. (%+v)", tbl)
	}
}

func TestAddIndexToTable(t *testing.T) {
	tbl := newUserTable()
	col := Column{name: "id"}
	tbl.AddColumn(&col)
	idx := NewIndex("", "", "users_pk", true)
	idx.AddColumn(&col)
	if len(tbl.Indices()) != 0 {
		t.Error("If table has no index, Indices() should be zero length.")
	}
	tbl.AddIndex(&idx)
	if len(tbl.Indices()) != 1 {
		t.Errorf("If table has a index, Indices() should be 1 length. (%+v)", tbl)
	}
	if tbl.Indices()[0].Schema() != tbl.Schema() {
		t.Errorf("Index's schema should be set by table's schema.")
	}
	if tbl.Indices()[0].TableName() != tbl.Name() {
		t.Errorf("Index's table name should be set by table's name.")
	}
}

func TestFindColumn(t *testing.T) {
	tbl := newUserTable()
	col := Column{name: "id"}
	tbl.AddColumn(&col)
	col = Column{name: "name"}
	tbl.AddColumn(&col)
	fc, err := tbl.FindColumn("name")
	if err != nil || fc.Name() != "name" {
		t.Error("FindColumn should return name column")
	}
	fc, err = tbl.FindColumn("login")
	if err == nil {
		t.Error("FindColumn should raise error when given not having column name.")
	}
}

func newUserTable() *Table {
	table := NewTable("foo", "users", "")
	return &table
}
