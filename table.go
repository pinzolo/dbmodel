package dbmodel

// Table stores table meta data.
type Table struct {
	schema      string
	name        string
	comment     string
	columns     []*Column
	indices     []*Index
	foreignKeys []*ForeignKey
	refKeys     []*ForeignKey
	constraints []*Constraint
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

// Constraints returns having constraints.
func (t Table) Constraints() []*Constraint {
	return t.constraints
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
		constraints: make([]*Constraint, 0, 5),
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

// AddConstraint appends foreign key to Constraints.
func (t *Table) AddConstraint(c *Constraint) {
	c.schema = t.schema
	c.tableName = t.name
	t.constraints = append(t.constraints, c)
}

// FindColumn returns column that has same name as argument.
// If column that has same name does not exist, return false as second value.
func (t *Table) FindColumn(name string) (*Column, bool) {
	for _, col := range t.Columns() {
		if col.Name() == name {
			return col, true
		}
	}
	return nil, false
}

// FindIndex returns index that has same name as argument.
// If index that has same name does not exist, return false as second value.
func (t *Table) FindIndex(name string) (*Index, bool) {
	for _, idx := range t.Indices() {
		if idx.Name() == name {
			return idx, true
		}
	}
	return nil, false
}

// FindForeignKey returns index that has same name as argument.
// If foreign key that has same name does not exist, return false as second value.
func (t *Table) FindForeignKey(name string) (*ForeignKey, bool) {
	for _, fk := range t.ForeignKeys() {
		if fk.Name() == name {
			return fk, true
		}
	}
	return nil, false
}

// FindReferencedKey returns index that has same name as argument.
// If referenced key that has same name does not exist, return false as second value.
func (t *Table) FindReferencedKey(name string) (*ForeignKey, bool) {
	for _, rk := range t.ReferencedKeys() {
		if rk.Name() == name {
			return rk, true
		}
	}
	return nil, false
}
