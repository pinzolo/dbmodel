package dbmodel

import "fmt"

// Table stores table meta data.
type Table struct {
	schema      string
	name        string
	comment     string
	columns     []*Column
	indices     []*Index
	foreignKeys []*ForeignKey
	refKeys     []*ForeignKey
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

// ForeignKeys returns having foreign keys.
func (t Table) ForeignKeys() []*ForeignKey {
	return t.foreignKeys
}

// ReferencedKeys returns having referenced keys.
func (t Table) ReferencedKeys() []*ForeignKey {
	return t.refKeys
}

// NewTable returns new Table initialized with arguments.
func NewTable(schema string, tableName string, comment string) Table {
	return Table{
		schema:      schema,
		name:        tableName,
		comment:     comment,
		columns:     make([]*Column, 0, 10),
		indices:     make([]*Index, 0, 5),
		foreignKeys: make([]*ForeignKey, 0, 5),
		refKeys:     make([]*ForeignKey, 0, 5),
	}
}

// AddColumn appends column to Columns.
func (t *Table) AddColumn(col *Column) {
	col.schema = t.schema
	col.tableName = t.name
	t.columns = append(t.columns, col)
}

// AddIndex appends index to Indecis.
func (t *Table) AddIndex(idx *Index) {
	idx.schema = t.schema
	idx.tableName = t.name
	t.indices = append(t.indices, idx)
}

// AddForeignKey appends foreign key to ForeignKeys.
func (t *Table) AddForeignKey(fk *ForeignKey) {
	fk.schema = t.schema
	fk.tableName = t.name
	t.foreignKeys = append(t.foreignKeys, fk)
}

// AddReferencedKey appends other table's foreign key that reference this table's column to ReferencedKeys.
func (t *Table) AddReferencedKey(rk *ForeignKey) {
	t.refKeys = append(t.refKeys, rk)
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
	return t.columns[len(t.columns)-1]
}

func (t *Table) lastIndex() *Index {
	return t.indices[len(t.indices)-1]
}

func (t *Table) lastForeignKey() *ForeignKey {
	return t.foreignKeys[len(t.foreignKeys)-1]
}

func (t *Table) lastRefKey() *ForeignKey {
	return t.refKeys[len(t.refKeys)-1]
}
