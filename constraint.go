package dbmodel

// Constraint is constraint's meta data.
// This struct using unique constraint and check constraint.
type Constraint struct {
	schema    string
	tableName string
	name      string
	kind      string
	content   string
}

// Schema returns constraint's schema.
func (c Constraint) Schema() string {
	return c.schema
}

// TableName returns table name that having this constraint.
func (c Constraint) TableName() string {
	return c.tableName
}

// Name returns constraint's name.
func (c Constraint) Name() string {
	return c.name
}

// Kind returns constraint's kind.
// If Constraint is unique constraint, Kind returns 'UNIQUE'.
// If Constraint is check constraint, Kind returns 'CHECK'.
func (c Constraint) Kind() string {
	return c.kind
}

// Content returns constraint's content.
func (c Constraint) Content() string {
	return c.content
}

// NewConstraint returns new Constraint initialized with arguments.
func NewConstraint(schema string, tableName string, name string, kind string, content string) Constraint {
	return Constraint{
		schema:    schema,
		tableName: tableName,
		name:      name,
		kind:      kind,
		content:   content,
	}
}
