package dbmodel

// Table stores table meta data.
type Table struct {
	name    string
	comment string
	columns []*Column
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
func NewTable(tableName string, comment string) Table {
	return Table{
		name:    tableName,
		comment: comment,
		columns: make([]*Column, 0, 10), // default size 10
	}
}

// AddColumn append column to Columns
func (t *Table) AddColumn(col *Column) {
	t.columns = append(t.columns, col)
}
