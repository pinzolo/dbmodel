package dbmodel

// Column is database column metadata.
type Column struct {
	name         string
	dataType     string
	size         Size
	nullable     bool
	defaultValue string
}

// Name returns column's name.
func (c Column) Name() string {
	return c.name
}

// DataType returns column's data type.
func (c Column) DataType() string {
	return c.dataType
}

// Size returns column's size.
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

// NewColumn returns new Column initialized with arguments.
func NewColumn(name string, dataType string, size Size, nullable bool, defaultValue string) Column {
	return Column{
		name:         name,
		dataType:     dataType,
		size:         size,
		nullable:     nullable,
		defaultValue: defaultValue,
	}
}
