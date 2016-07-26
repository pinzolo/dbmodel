package dbmodel

import "database/sql"

// Provider is interface to absorbe difference of each database.
type Provider interface {
	// Connect open connection to DataSouce.
	Connect(ds DataSource) (*sql.DB, error)
	// AllTableNamesSQL should return SQL for loading table names.
	// Parameters:
	//     1. schema
	// Return columns:
	//     1. schema
	//     2. table name
	//     3. table comment
	// Order: table name
	AllTableNamesSQL() string
	// TableNamesSQL should return SQL for loading table names using LIKE table_name.
	// Parameters:
	//     1. schema
	//     2. table name
	// Return columns:
	//     same as AllTableNamesSQL
	// Order: table name
	TableNamesSQL() string
	// AllTableSQL should return SQL for loading all tables contains columns.
	// Parameters:
	//     1. schema
	// Return columns:
	//      1. schema
	//      2. table name
	//      3. table comment
	//      4. column name
	//      5. column comment
	//      6. data type
	//      7. length (using in text)
	//      8. precision (using in numeric or date)
	//      9. scale (using in numeric)
	//     10. nullable ("YES" or "NO")
	//     11. default value (as text)
	//     12. primary key position
	// Order:
	//     1. table name
	//     2. column position
	AllTablesSQL() string
	// TableSQL should return SQL for loading a table contains columns.
	// Parameters:
	//     1. schema
	//     2. table name
	// Return columns:
	//     same as AllTableSQL
	// Order:
	//     1. column position
	TableSQL() string
	// AllIndicesSQL should return SQL for loading all indices.
	// Parameters:
	//     1. schema
	// Return columns:
	//     1. schema
	//     2. table name
	//     3. index name
	//     4. unique ("YES" or "NO")
	//     5. column name
	// Order:
	//     1. table name
	//     2. index name
	//     3. column position
	AllIndicesSQL() string
	// IndicesSQL should return SQL for loading indices in a table.
	// Parameters:
	//     1. schema
	//     2. table name
	// Return columns:
	//     same as AllIndicesSQL
	// Order:
	//     1. index name
	//     2. column position
	IndicesSQL() string
	// AllForeignKeysSQL should return SQL for loading all foreign keys.
	// Parameters:
	//     1. schema
	// Return columns:
	//     1. foreign key name
	//     2. schema      (from)
	//     3. table name  (from)
	//     4. column name (from)
	//     5. schema      (to)
	//     6. table name  (to)
	//     7. column name (to)
	// Order:
	//     1. table name (from)
	//     2. foreign key name
	//     3. column position (from)
	AllForeignKeysSQL() string
	// ForeignKeysSQL should return SQL for loading foreign keys in a table.
	// Parameters:
	//     1. schema
	//     2. table name
	// Return columns:
	//     same as AllForeignKeysSQL
	// Order:
	//     1. foreign key name
	//     2. column position (from)
	ForeignKeysSQL() string
	// AllReferencedKeysSQL should return SQL for loading all foreign keys.
	// Parameters:
	//     1. schema
	// Return columns:
	//     1. foreign key name
	//     2. schema      (from)
	//     3. table name  (from)
	//     4. column name (from)
	//     5. schema      (to)
	//     6. table name  (to)
	//     7. column name (to)
	// Order:
	//     1. table name (to)
	//     2. foreign key name
	//     3. column position (to)
	AllReferencedKeysSQL() string
	// ForeignKeysSQL should return SQL for loading referenced foreign keys in a table.
	// Parameters:
	//     1. schema
	//     2. table name
	// Return columns:
	//     same as AllReferencedKeysSQL
	// Order:
	//     1. foreign key name
	//     2. column position (to)
	ReferencedKeysSQL() string
	AllConstrantsSQL() string
	ConstrantsSQL() string
}
