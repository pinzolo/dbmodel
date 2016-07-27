package dbmodel

// Option is table loding option for Table and AllTables function.
type Option struct {
	Indices        bool
	ForeignKeys    bool
	ReferencedKeys bool
	Constraints    bool
}

var (
	// RequireAll is loading option for loading all meta data.
	RequireAll = Option{
		Indices:        true,
		ForeignKeys:    true,
		ReferencedKeys: true,
		Constraints:    true,
	}
	// RequireNone is loading option for loading only columns.
	RequireNone = Option{
		Indices:        false,
		ForeignKeys:    false,
		ReferencedKeys: false,
		Constraints:    false,
	}
)
