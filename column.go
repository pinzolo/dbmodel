package dbmodel

// Column is database column metadata.
type Column struct {
	schema       string
	tableName    string
	name         string
	comment      string
	dataType     string
	size         Size
	nullable     bool
	defaultValue string
	pkPosition   int64
}

// Schema returns column schema.
func (c Column) Schema() string {
	return c.schema
}

// TableName returns table name that this column belonging.
func (c Column) TableName() string {
	return c.tableName
}

// Name returns column name.
func (c Column) Name() string {
	return c.name
}

// Comment returns columns comment
func (c Column) Comment() string {
	return c.comment
}

// DataType returns column data type.
func (c Column) DataType() string {
	return c.dataType
}

// Size returns column size.
func (c Column) Size() Size {
	return c.size
}

// IsNullable returns true if column can accept NULL.
func (c Column) IsNullable() bool {
	return c.nullable
}

// DefaultValue returns column's default value
func (c Column) DefaultValue() string {
	return c.defaultValue
}

// PrimaryKeyPosition returns this column's position in primary key columns.
// If this column is not primary key, returns 0.
func (c Column) PrimaryKeyPosition() int64 {
	return c.pkPosition
}

// NewColumn returns new Column initialized with arguments.
func NewColumn(schema string, tableName string, name string, comment string, dataType string, size Size, nullable bool, defaultValue string, pkPosition int64) Column {
	return Column{
		schema:       schema,
		tableName:    tableName,
		name:         name,
		comment:      comment,
		dataType:     dataType,
		size:         size,
		nullable:     nullable,
		defaultValue: defaultValue,
		pkPosition:   pkPosition,
	}
}
