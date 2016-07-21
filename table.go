package dbmodel

import "fmt"

// Table stores table meta data.
type Table struct {
	schema  string
	name    string
	comment string
	columns []*Column
	indices []*Index
}

// Schema returns table schema.
func (t Table) Schema() string {
	return t.schema
}

// Name returns table name.
func (t Table) Name() string {
	return t.name
}

// Comment returns table comment.
func (t Table) Comment() string {
	return t.comment
}

// Columns returns having columns.
func (t Table) Columns() []*Column {
	return t.columns
}

// Indices returns having indices.
func (t Table) Indices() []*Index {
	return t.indices
}

// NewTable returns new Table initialized with arguments.
func NewTable(schema string, tableName string, comment string) Table {
	return Table{
		schema:  schema,
		name:    tableName,
		comment: comment,
		columns: make([]*Column, 0, 10),
		indices: make([]*Index, 0, 5),
	}
}

// AddColumn append column to Columns.
func (t *Table) AddColumn(col *Column) {
	col.schema = t.schema
	col.tableName = t.name
	t.columns = append(t.columns, col)
}

// AddIndex append index to Indecis.
func (t *Table) AddIndex(idx *Index) {
	idx.schema = t.schema
	idx.tableName = t.name
	t.indices = append(t.indices, idx)
}

// FindColumn returns column that name is same as argument.
func (t *Table) FindColumn(name string) (*Column, error) {
	for _, col := range t.Columns() {
		if col.Name() == name {
			return col, nil
		}
	}
	return nil, fmt.Errorf("Column '%v' is not found in '%v' table.", name, t.Name())
}

func (t *Table) lastColumn() *Column {
	return t.Columns()[len(t.Columns())-1]
}

func (t *Table) lastIndex() *Index {
	return t.Indices()[len(t.Indices())-1]
}
