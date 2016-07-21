package dbmodel

// Index is database index metadata.
type Index struct {
	schema    string
	tableName string
	name      string
	unique    bool
	columns   []*Column
}

// Schema returns index schema.
func (i Index) Schema() string {
	return i.schema
}

// TableName returns table name that having this index.
func (i Index) TableName() string {
	return i.tableName
}

// Name returns index name.
func (i Index) Name() string {
	return i.name
}

// IsUnique returns true if this index is unique index.
func (i Index) IsUnique() bool {
	return i.unique
}

// Columns returns having columns.
func (i Index) Columns() []*Column {
	return i.columns
}

// NewIndex returns new Index initialized with arguments.
func NewIndex(schema string, tableName string, name string, unique bool) Index {
	return Index{
		schema:    schema,
		tableName: tableName,
		name:      name,
		unique:    unique,
		columns:   make([]*Column, 0, 5), // default size 5
	}
}

// AddColumn append column to Columns
func (i *Index) AddColumn(col *Column) {
	i.columns = append(i.columns, col)
}
