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
     , d.description AS table_comment
FROM pg_catalog.pg_tables t
LEFT OUTER JOIN pg_catalog.pg_class c1
ON c1.relname = t.tablename
LEFT OUTER JOIN pg_catalog.pg_description d
ON  d.objoid = c1.oid
AND d.objsubid = 0
WHERE t.schemaname  = $1
ORDER BY t.tablename`
}

func (p postgres) TableNamesSQL() string {
	return `
SELECT t.tablename AS table_name
     , d.description AS table_comment
FROM pg_catalog.pg_tables t
LEFT OUTER JOIN pg_catalog.pg_class c1
ON c1.relname = t.tablename
LEFT OUTER JOIN pg_catalog.pg_description d
ON  d.objoid = c1.oid
AND d.objsubid = 0
WHERE t.schemaname  = $1
AND   t.tablename LIKE '%' || $2 || '%'
ORDER BY t.tablename`
}

func (p postgres) TableSQL() string {
	return `
SELECT cls.relname AS table_name
     , td.description AS table_comment
     , att.attname AS column_name
     , cd.description AS column_comment
     , COALESCE(col.domain_schema || '.' || col.domain_name, col.udt_name) AS data_type
     , col.character_maximum_length AS length
     , COALESCE(col.numeric_precision, col.datetime_precision) AS precision
     , col.numeric_scale AS scale
     , col.is_nullable AS nullable
     , col.column_default AS defaul_value
	 , pk.ordinal_position AS primary_key_position
FROM pg_catalog.pg_class cls
INNER JOIN pg_catalog.pg_namespace ns
ON  cls.relnamespace = ns.oid
LEFT OUTER JOIN pg_catalog.pg_description td
ON  cls.oid = td.objoid
AND td.objsubid = 0
INNER JOIN pg_catalog.pg_attribute att
ON  cls.oid = att.attrelid
AND att.attnum > 0
LEFT OUTER JOIN pg_catalog.pg_description cd
ON  cls.oid = cd.objoid
AND att.attnum = cd.objsubid
INNER JOIN information_schema.columns col
ON  col.table_schema = ns.nspname
AND col.table_name = cls.relname
AND col.column_name = att.attname
LEFT OUTER JOIN (
	SELECT tc.table_schema
	     , tc.table_name
		 , kcu.column_name
		 , kcu.ordinal_position
	FROM information_schema.key_column_usage kcu
	INNER JOIN information_schema.constraint_column_usage ccu
	ON  ccu.table_catalog = kcu.table_catalog
    AND ccu.table_schema = kcu.table_schema
    AND ccu.table_name = kcu.table_name
	AND ccu.column_name = kcu.column_name
	INNER JOIN information_schema.table_constraints tc
	ON  tc.table_name = ccu.table_name
	AND tc.constraint_name = ccu.constraint_name
	AND tc.constraint_name = kcu.constraint_name
	WHERE tc.constraint_type = 'PRIMARY KEY'
) pk
ON  pk.table_schema = col.table_schema
AND pk.table_name = col.table_name
AND pk.column_name = col.column_name
WHERE cls.relkind = 'r'
AND   ns.nspname = $1
AND   cls.relname = $2
ORDER BY col.ordinal_position`
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
