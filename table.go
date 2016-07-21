package dbmodel

// Table stores table meta data.
type Table struct {
	schema  string
	name    string
	comment string
	columns []*Column
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

// NewTable returns new Table initialized with arguments.
func NewTable(schema string, tableName string, comment string) Table {
	return Table{
		schema:  schema,
		name:    tableName,
		comment: comment,
		columns: make([]*Column, 0, 10), // default size 10
	}
}

// AddColumn append column to Columns.
func (t *Table) AddColumn(col *Column) {
	col.schema = t.schema
	col.tableName = t.name
	t.columns = append(t.columns, col)
}
