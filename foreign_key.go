package dbmodel

// ForeignKey is foreign key's meta data.
type ForeignKey struct {
	schema    string
	tableName string
	name      string
	colRefs   []*ColumnReference
}

// Schema returns foreign key's schema.
func (fk ForeignKey) Schema() string {
	return fk.schema
}

// TableName returns table name that having this foreign key.
func (fk ForeignKey) TableName() string {
	return fk.tableName
}

// Name returns foreign key's name.
func (fk ForeignKey) Name() string {
	return fk.name
}

// ColumnReferences returns foreign key's references from column to other column.
func (fk ForeignKey) ColumnReferences() []*ColumnReference {
	return fk.colRefs
}

// AddColumnReference appends column reference to ColumnReferences.
func (fk *ForeignKey) AddColumnReference(r *ColumnReference) {
	fk.colRefs = append(fk.colRefs, r)
}

// NewForeignKey returns new ForeignKey initialized with arguments.
func NewForeignKey(schema string, tableName string, name string) ForeignKey {
	return ForeignKey{
		schema:    schema,
		tableName: tableName,
		name:      name,
		colRefs:   make([]*ColumnReference, 0, 2),
	}
}
