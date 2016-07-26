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
SELECT t.schemaname AS schema
     , t.tablename AS table_name
     , d.description AS table_comment
FROM pg_catalog.pg_tables t
LEFT OUTER JOIN pg_catalog.pg_class c1
ON c1.relname = t.tablename
LEFT OUTER JOIN pg_catalog.pg_description d
ON  d.objoid = c1.oid
AND d.objsubid = 0
WHERE t.schemaname = $1
ORDER BY t.tablename`
}

func (p postgres) TableNamesSQL() string {
	return `
SELECT t.schemaname AS schema
     , t.tablename AS table_name
     , d.description AS table_comment
FROM pg_catalog.pg_tables t
LEFT OUTER JOIN pg_catalog.pg_class c1
ON c1.relname = t.tablename
LEFT OUTER JOIN pg_catalog.pg_description d
ON  d.objoid = c1.oid
AND d.objsubid = 0
WHERE t.schemaname = $1
AND   t.tablename LIKE '%' || $2 || '%'
ORDER BY t.tablename`
}

func (p postgres) TableSQL() string {
	return `
SELECT ns.nspname AS schema
     , cls.relname AS table_name
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
    FROM information_schema.table_constraints tc
    INNER JOIN information_schema.key_column_usage kcu
    ON  kcu.constraint_catalog = tc.constraint_catalog
    AND kcu.constraint_schema = tc.constraint_schema
    AND kcu.constraint_name = tc.constraint_name
    WHERE tc.constraint_type = 'PRIMARY KEY'
    AND   tc.table_schema = $1
    AND   tc.table_name = $2
    AND   kcu.table_schema = $1
    AND   kcu.table_name = $2
) pk
ON  pk.table_schema = col.table_schema
AND pk.table_name = col.table_name
AND pk.column_name = col.column_name
WHERE cls.relkind = 'r'
AND   ns.nspname = $1
AND   cls.relname = $2
ORDER BY cls.relname, col.ordinal_position`
}

func (p postgres) AllTablesSQL() string {
	return `
SELECT ns.nspname AS schema
     , cls.relname AS table_name
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
    FROM information_schema.table_constraints tc
    INNER JOIN information_schema.key_column_usage kcu
    ON  kcu.constraint_catalog = tc.constraint_catalog
    AND kcu.constraint_schema = tc.constraint_schema
    AND kcu.constraint_name = tc.constraint_name
    WHERE tc.constraint_type = 'PRIMARY KEY'
    AND   tc.table_schema = $1
    AND   kcu.table_schema = $1
) pk
ON  pk.table_schema = col.table_schema
AND pk.table_name = col.table_name
AND pk.column_name = col.column_name
WHERE cls.relkind = 'r'
AND   ns.nspname = $1
ORDER BY cls.relname, col.ordinal_position`
}

func (p postgres) IndicesSQL() string {
	return `
SELECT idxs.schemaname AS schema
     , col.table_name
     , idxs.indexname AS index_name
     , CASE WHEN idx.uniq THEN 'YES' ELSE 'NO' END AS uniq
     , col.column_name
FROM pg_catalog.pg_class cls
INNER JOIN pg_catalog.pg_namespace ns
ON  cls.relnamespace = ns.oid
INNER JOIN pg_catalog.pg_indexes idxs
ON  idxs.schemaname = ns.nspname
AND idxs.indexname = cls.relname
INNER JOIN (
    SELECT indexrelid
         , indisunique AS uniq
         , string_to_array(indkey::text, ' ')::int[] AS column_positions
         , generate_series(1, length(indkey::text) - length(replace(indkey::text, ' ', '')) + 1) AS column_ordinal
    FROM pg_catalog.pg_index
) idx
ON idx.indexrelid = cls.oid
JOIN information_schema.columns col
ON  col.table_schema = idxs.schemaname
AND col.table_name = idxs.tablename
AND col.ordinal_position = idx.column_positions[column_ordinal]
WHERE idxs.schemaname = $1
AND   idxs.tablename = $2
ORDER BY table_name, idxs.indexname, idx.column_ordinal`
}

func (p postgres) AllIndicesSQL() string {
	return `
SELECT idxs.schemaname AS schema
     , col.table_name
     , idxs.indexname AS index_name
     , CASE WHEN idx.uniq THEN 'YES' ELSE 'NO' END AS uniq
     , col.column_name
FROM pg_catalog.pg_class cls
INNER JOIN pg_catalog.pg_namespace ns
ON  cls.relnamespace = ns.oid
INNER JOIN pg_catalog.pg_indexes idxs
ON  idxs.schemaname = ns.nspname
AND idxs.indexname = cls.relname
INNER JOIN (
    SELECT indexrelid
         , indisunique AS uniq
         , string_to_array(indkey::text, ' ')::int[] AS column_positions
         , generate_series(1, length(indkey::text) - length(replace(indkey::text, ' ', '')) + 1) AS column_ordinal
    FROM pg_catalog.pg_index
) idx
ON idx.indexrelid = cls.oid
JOIN information_schema.columns col
ON  col.table_schema = idxs.schemaname
AND col.table_name = idxs.tablename
AND col.ordinal_position = idx.column_positions[column_ordinal]
WHERE idxs.schemaname = $1
ORDER BY table_name, idxs.indexname, idx.column_ordinal`
}

func (p postgres) ForeignKeysSQL() string {
	return `
SELECT cns.conname AS foreign_key_name
     , ns.nspname AS schema
     , cls.relname AS table_name
     , att.attname AS column_name
     , fns.nspname AS foreign_schema
     , fcls.relname AS foreign_table_name
     , fatt.attname AS foreign_column_name
FROM (
    SELECT conname
         , conrelid AS relid
         , conkey AS key
         , generate_series(1, length(array_to_string(conkey, ' ')) - length(array_to_string(conkey, '')) + 1) AS pos
    FROM pg_catalog.pg_constraint
    WHERE contype = 'f'
) AS cns
INNER JOIN pg_catalog.pg_class cls
ON cls.oid = cns.relid
INNER JOIN pg_catalog.pg_namespace ns
ON cls.relnamespace = ns.oid
INNER JOIN pg_catalog.pg_attribute att
ON  att.attrelid = cls.oid
AND att.attnum = cns.key[cns.pos]
JOIN (
    SELECT conname
         , confrelid AS relid
         , confkey AS key
         , generate_series(1, length(array_to_string(confkey, ' ')) - length(array_to_string(confkey, '')) + 1) AS pos
    FROM pg_catalog.pg_constraint
    WHERE contype = 'f'
) AS fcns
ON fcns.conname = cns.conname
AND fcns.pos = cns.pos
INNER JOIN pg_catalog.pg_class fcls
ON fcls.oid = fcns.relid
INNER JOIN pg_catalog.pg_namespace fns
ON fcls.relnamespace = fns.oid
INNER JOIN pg_catalog.pg_attribute fatt
ON  fatt.attrelid = fcls.oid
AND fatt.attnum = fcns.key[fcns.pos]
WHERE ns.nspname = $1
AND   cls.relname = $2
ORDER BY cls.relname, cns.conname, cns.pos`
}

func (p postgres) AllForeignKeysSQL() string {
	return `
SELECT cns.conname AS foreign_key_name
     , ns.nspname AS schema
     , cls.relname AS table_name
     , att.attname AS column_name
     , fns.nspname AS foreign_schema
     , fcls.relname AS foreign_table_name
     , fatt.attname AS foreign_column_name
FROM (
    SELECT conname
         , conrelid AS relid
         , conkey AS key
         , generate_series(1, length(array_to_string(conkey, ' ')) - length(array_to_string(conkey, '')) + 1) AS pos
    FROM pg_catalog.pg_constraint
    WHERE contype = 'f'
) AS cns
INNER JOIN pg_catalog.pg_class cls
ON cls.oid = cns.relid
INNER JOIN pg_catalog.pg_namespace ns
ON cls.relnamespace = ns.oid
INNER JOIN pg_catalog.pg_attribute att
ON  att.attrelid = cls.oid
AND att.attnum = cns.key[cns.pos]
JOIN (
    SELECT conname
         , confrelid AS relid
         , confkey AS key
         , generate_series(1, length(array_to_string(confkey, ' ')) - length(array_to_string(confkey, '')) + 1) AS pos
    FROM pg_catalog.pg_constraint
    WHERE contype = 'f'
) AS fcns
ON fcns.conname = cns.conname
AND fcns.pos = cns.pos
INNER JOIN pg_catalog.pg_class fcls
ON fcls.oid = fcns.relid
INNER JOIN pg_catalog.pg_namespace fns
ON fcls.relnamespace = fns.oid
INNER JOIN pg_catalog.pg_attribute fatt
ON  fatt.attrelid = fcls.oid
AND fatt.attnum = fcns.key[fcns.pos]
WHERE ns.nspname = $1
ORDER BY cls.relname, cns.conname, cns.pos`
}

func (p postgres) ReferencedKeysSQL() string {
	return `
SELECT cns.conname AS referenced_key_name
     , ns.nspname AS schema
     , cls.relname AS table_name
     , att.attname AS column_name
     , fns.nspname AS foreign_schema
     , fcls.relname AS foreign_table_name
     , fatt.attname AS foreign_column_name
FROM (
    SELECT conname
         , conrelid AS relid
         , conkey AS key
         , generate_series(1, length(array_to_string(conkey, ' ')) - length(array_to_string(conkey, '')) + 1) AS pos
    FROM pg_catalog.pg_constraint
    WHERE contype = 'f'
) AS cns
INNER JOIN pg_catalog.pg_class cls
ON cls.oid = cns.relid
INNER JOIN pg_catalog.pg_namespace ns
ON cls.relnamespace = ns.oid
INNER JOIN pg_catalog.pg_attribute att
ON  att.attrelid = cls.oid
AND att.attnum = cns.key[cns.pos]
JOIN (
    SELECT conname
         , confrelid AS relid
         , confkey AS key
         , generate_series(1, length(array_to_string(confkey, ' ')) - length(array_to_string(confkey, '')) + 1) AS pos
    FROM pg_catalog.pg_constraint
    WHERE contype = 'f'
) AS fcns
ON fcns.conname = cns.conname
AND fcns.pos = cns.pos
INNER JOIN pg_catalog.pg_class fcls
ON fcls.oid = fcns.relid
INNER JOIN pg_catalog.pg_namespace fns
ON fcls.relnamespace = fns.oid
INNER JOIN pg_catalog.pg_attribute fatt
ON  fatt.attrelid = fcls.oid
AND fatt.attnum = fcns.key[fcns.pos]
WHERE fns.nspname = $1
AND   fcls.relname = $2
ORDER BY fcls.relname, fcns.conname, fcns.pos`
}

func (p postgres) AllReferencedKeysSQL() string {
	return `
SELECT cns.conname AS referenced_key_name
     , ns.nspname AS schema
     , cls.relname AS table_name
     , att.attname AS column_name
     , fns.nspname AS foreign_schema
     , fcls.relname AS foreign_table_name
     , fatt.attname AS foreign_column_name
FROM (
    SELECT conname
         , conrelid AS relid
         , conkey AS key
         , generate_series(1, length(array_to_string(conkey, ' ')) - length(array_to_string(conkey, '')) + 1) AS pos
    FROM pg_catalog.pg_constraint
    WHERE contype = 'f'
) AS cns
INNER JOIN pg_catalog.pg_class cls
ON cls.oid = cns.relid
INNER JOIN pg_catalog.pg_namespace ns
ON cls.relnamespace = ns.oid
INNER JOIN pg_catalog.pg_attribute att
ON  att.attrelid = cls.oid
AND att.attnum = cns.key[cns.pos]
JOIN (
    SELECT conname
         , confrelid AS relid
         , confkey AS key
         , generate_series(1, length(array_to_string(confkey, ' ')) - length(array_to_string(confkey, '')) + 1) AS pos
    FROM pg_catalog.pg_constraint
    WHERE contype = 'f'
) AS fcns
ON fcns.conname = cns.conname
AND fcns.pos = cns.pos
INNER JOIN pg_catalog.pg_class fcls
ON fcls.oid = fcns.relid
INNER JOIN pg_catalog.pg_namespace fns
ON fcls.relnamespace = fns.oid
INNER JOIN pg_catalog.pg_attribute fatt
ON  fatt.attrelid = fcls.oid
AND fatt.attnum = fcns.key[fcns.pos]
WHERE fns.nspname = $1
ORDER BY fcls.relname, fcns.conname, fcns.pos`
}

func (p postgres) ConstrantsSQL() string {
	return ``
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
