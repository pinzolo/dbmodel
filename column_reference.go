package dbmodel

// ColumnReference is reference to column from column.
type ColumnReference struct {
	from *Column
	to   *Column
}

// From returns form column of this reference.
func (cr ColumnReference) From() *Column {
	return cr.from
}

// To returns to column of this reference.
func (cr ColumnReference) To() *Column {
	return cr.to
}

// NewColumnReference returns new ColumnRef initialized with arguments.
func NewColumnReference(from *Column, to *Column) ColumnReference {
	return ColumnReference{
		from: from,
		to:   to,
	}
}
