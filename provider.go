package dbmodel

import "database/sql"

// Provider is interface to absorbe difference of each database.
type Provider interface {
	Connect(ds DataSource) (*sql.DB, error)
	AllTableNamesSQL() string
	TableNamesSQL() string
	TableSQL() string
	AllTablesSQL() string
	IndicesSQL() string
	AllIndicesSQL() string
	ForeignKeysSQL() string
	AllForeignKeysSQL() string
	ReferencedKeysSQL() string
	AllReferencedKeysSQL() string
}
