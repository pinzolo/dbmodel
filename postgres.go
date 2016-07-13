package dbmodel

import (
	"database/sql"
	"strconv"
	"strings"
)

// postgres is Provider implementation for PostgreSQL.
type postgres struct{}

// Connect to PostgreSQL using given DataSource setting.
func (p postgres) Connect(ds DataSource) (*sql.DB, error) {
	return sql.Open("postgres", p.dataSourceName(ds))
}

func (p postgres) AllTableNamesSQL() string {
	return `
SELECT t.tablename AS table_name
     , d.description AS comment
FROM pg_catalog.pg_tables t
LEFT JOIN pg_catalog.pg_class c1
ON c1.relname = t.tablename
LEFT JOIN pg_catalog.pg_description d
ON d.objoid = c1.oid
AND d.objsubid = 0
WHERE t.schemaname  = $1
ORDER BY t.tablename`
}

func (p postgres) TableNamesSQL() string {
	return `
SELECT t.tablename AS table_name
     , d.description AS comment
FROM pg_catalog.pg_tables t
LEFT JOIN pg_catalog.pg_class c1
ON c1.relname = t.tablename
LEFT JOIN pg_catalog.pg_description d
ON d.objoid = c1.oid
AND d.objsubid = 0
WHERE t.schemaname  = $1
AND t.tablename LIKE '%' || $2 || '%'
ORDER BY t.tablename`
}

func (p postgres) dataSourceName(ds DataSource) string {
	parts := make([]string, 0, 10)
	if ds.Host != "" {
		parts = append(parts, "host="+ds.Host)
	}
	if ds.Port != 0 {
		parts = append(parts, "port="+strconv.Itoa(ds.Port))
	}
	if ds.User != "" {
		parts = append(parts, "user="+ds.User)
	}
	if ds.Password != "" {
		parts = append(parts, "password="+ds.Password)
	}
	if ds.Database != "" {
		parts = append(parts, "dbname="+ds.Database)
	}
	for k, v := range ds.Options {
		parts = append(parts, k+"="+v)
	}
	return strings.Join(parts, " ")
}
