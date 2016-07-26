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

func TestAddForeignKeyToTable(t *testing.T) {
	usr := newUserTable()
	usr.AddColumn(&Column{name: "id"})
	pst := newPostTable()
	pst.AddColumn(&Column{name: "id"})
	pst.AddColumn(&Column{name: "user_id"})
	if len(pst.ForeignKeys()) != 0 {
		t.Error("If table has no foreign key, ForeignKeys() should be empty.")
	}
	fk := NewForeignKey("foo", "posts", "posts_user_id")
	colUserID, _ := pst.FindColumn("user_id")
	colID, _ := usr.FindColumn("id")
	cr := NewColumnReference(colUserID, colID)
	fk.AddColumnReference(&cr)
	pst.AddForeignKey(&fk)
	if len(pst.ForeignKeys()) != 1 {
		t.Errorf("If table has a foreign key, ForeignKeys() should be 1 length. (%+v)", pst.ForeignKeys())
	}
	if pst.ForeignKeys()[0].Schema() != pst.Schema() {
		t.Error("Foreign key's schema should be set by table's schema.")
	}
	if pst.ForeignKeys()[0].TableName() != pst.Name() {
		t.Error("Foreign key's table name should be set by table's name.")
	}
	if pst.ForeignKeys()[0].Name() != "posts_user_id" {
		t.Error("Invalid forein key is added.")
	}
}

func TestAddReferencedKeyToTable(t *testing.T) {
	usr := newUserTable()
	usr.AddColumn(&Column{name: "id"})
	pst := newPostTable()
	pst.AddColumn(&Column{name: "id"})
	pst.AddColumn(&Column{name: "user_id"})
	if len(usr.ReferencedKeys()) != 0 {
		t.Error("If table has no referenced key, ReferencedKeys() should be empty.")
	}
	fk := NewForeignKey("foo", "posts", "posts_user_id")
	colUserID, _ := pst.FindColumn("user_id")
	colID, _ := usr.FindColumn("id")
	cr := NewColumnReference(colUserID, colID)
	fk.AddColumnReference(&cr)
	usr.AddReferencedKey(&fk)
	if len(usr.ReferencedKeys()) != 1 {
		t.Errorf("If table has a referenced key, ReferencedKeys() should be 1 length. (%+v)", usr.ReferencedKeys())
	}
	if usr.ReferencedKeys()[0].Name() != "posts_user_id" {
		t.Error("Invalid referenced key is added.")
	}
}

func TestFindColumn(t *testing.T) {
	tbl := newUserTable()
	col := Column{name: "id"}
	tbl.AddColumn(&col)
	col = Column{name: "name"}
	tbl.AddColumn(&col)
	fc, ok := tbl.FindColumn("name")
	if !ok || fc.Name() != "name" {
		t.Error("FindColumn should return name column")
	}
	fc, ok = tbl.FindColumn("login")
	if ok {
		t.Error("FindColumn should raise error when given not having column name.")
	}
}

func TestAddConstraintToTable(t *testing.T) {
	tbl := newUserTable()
	col := Column{name: "age"}
	tbl.AddColumn(&col)
	if len(tbl.Constraints()) != 0 {
		t.Error("If table has no constraint, Constraints() should be zero length.")
	}
	con := NewConstraint("", "", "users_age_check", "CHECK", "(age >= 0)")
	tbl.AddConstraint(&con)
	if len(tbl.Constraints()) != 1 {
		t.Errorf("If table has a constraint, Constraints() should be 1 length. (%+v)", tbl.Constraints())
	}
	if tbl.Constraints()[0].Schema() != tbl.Schema() {
		t.Errorf("Constraint's schema should be set by table's schema.")
	}
	if tbl.Constraints()[0].TableName() != tbl.Name() {
		t.Errorf("Constraint's table name should be set by table's name.")
	}
}

func newUserTable() *Table {
	table := NewTable("foo", "users", "")
	return &table
}

func newPostTable() *Table {
	table := NewTable("foo", "posts", "")
	return &table
}
