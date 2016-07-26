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
     , pk.pos AS primary_key_position
FROM pg_catalog.pg_class cls
INNER JOIN pg_catalog.pg_namespace ns
ON  cls.relnamespace = ns.oid
LEFT OUTER JOIN pg_catalog.pg_description td
ON  cls.oid = td.objoid
AND td.objsubid = 0
INNER JOIN pg_catalog.pg_attribute att
ON  att.attrelid = cls.oid
AND att.attnum > 0
LEFT OUTER JOIN pg_catalog.pg_description cd
ON  cls.oid = cd.objoid
AND att.attnum = cd.objsubid
INNER JOIN information_schema.columns col
ON  col.table_schema = ns.nspname
AND col.table_name = cls.relname
AND col.column_name = att.attname
LEFT OUTER JOIN (
    SELECT conrelid
         , conname
         , conkey AS colnums
         , generate_series(1, length(array_to_string(conkey, ' ')) - length(array_to_string(conkey, '')) + 1) AS pos
    FROM pg_catalog.pg_constraint
    WHERE contype = 'p'
) pk
ON  pk.conrelid = cls.oid
AND att.attnum = pk.colnums[pk.pos]
WHERE cls.relkind = 'r'
AND   ns.nspname = $1
AND   cls.relname = $2
ORDER BY cls.relname, att.attnum`
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
     , pk.pos AS primary_key_position
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
    SELECT conrelid
         , conname
         , conkey AS colnums
         , generate_series(1, length(array_to_string(conkey, ' ')) - length(array_to_string(conkey, '')) + 1) AS pos
    FROM pg_catalog.pg_constraint
    WHERE contype = 'p'
) pk
ON  pk.conrelid = cls.oid
AND att.attnum = pk.colnums[pk.pos]
WHERE cls.relkind = 'r'
AND   ns.nspname = $1
ORDER BY cls.relname, col.ordinal_position`
}

func (p postgres) IndicesSQL() string {
	return `
SELECT ns.nspname AS schema
     , tcls.relname AS table_name
     , icls.relname AS index_name
     , CASE WHEN idx.uniq THEN 'YES' ELSE 'NO' END AS uniq
     , att.attname AS column_name
FROM (
    SELECT indexrelid AS index_oid
         , indrelid AS table_oid
         , indisunique AS uniq
         , string_to_array(indkey::text, ' ')::int[] AS colnums
         , generate_series(1, length(indkey::text) - length(replace(indkey::text, ' ', '')) + 1) AS pos
    FROM pg_catalog.pg_index
) idx
INNER JOIN pg_catalog.pg_class tcls
ON tcls.oid = idx.table_oid
INNER JOIN pg_catalog.pg_namespace ns
ON tcls.relnamespace = ns.oid
INNER JOIN pg_catalog.pg_class icls
ON icls.oid = idx.index_oid
JOIN pg_catalog.pg_attribute att
ON  att.attrelid = tcls.oid
AND att.attnum = idx.colnums[idx.pos]
WHERE ns.nspname = $1
AND   tcls.relname = $2
ORDER BY tcls.relname, icls.relname, idx.pos`
}

func (p postgres) AllIndicesSQL() string {
	return `
SELECT ns.nspname AS schema
     , tcls.relname AS table_name
     , icls.relname AS index_name
     , CASE WHEN idx.uniq THEN 'YES' ELSE 'NO' END AS uniq
     , att.attname AS column_name
FROM (
    SELECT indexrelid AS index_oid
         , indrelid AS table_oid
         , indisunique AS uniq
         , string_to_array(indkey::text, ' ')::int[] AS colnums
         , generate_series(1, length(indkey::text) - length(replace(indkey::text, ' ', '')) + 1) AS pos
    FROM pg_catalog.pg_index
) idx
INNER JOIN pg_catalog.pg_class tcls
ON tcls.oid = idx.table_oid
INNER JOIN pg_catalog.pg_namespace ns
ON tcls.relnamespace = ns.oid
INNER JOIN pg_catalog.pg_class icls
ON icls.oid = idx.index_oid
JOIN pg_catalog.pg_attribute att
ON  att.attrelid = tcls.oid
AND att.attnum = idx.colnums[idx.pos]
WHERE ns.nspname = $1
ORDER BY tcls.relname, icls.relname, idx.pos`
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
         , conkey AS colnums
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
AND att.attnum = cns.colnums[cns.pos]
JOIN (
    SELECT conname
         , confrelid AS relid
         , confkey AS colnums
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
AND fatt.attnum = fcns.colnums[fcns.pos]
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
         , conkey AS colnums
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
AND att.attnum = cns.colnums[cns.pos]
JOIN (
    SELECT conname
         , confrelid AS relid
         , confkey AS colnums
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
AND fatt.attnum = fcns.colnums[fcns.pos]
WHERE fns.nspname = $1
ORDER BY fcls.relname, fcns.conname, fcns.pos`
}

func (p postgres) ConstrantsSQL() string {
	return `
SELECT ns.nspname AS schema
     , cls.relname AS table_name
     , cns.conname AS constraint_name
     , 'CHECK' AS constraint_kind
     , cns.consrc AS constraint_content
FROM pg_catalog.pg_constraint cns
JOIN pg_catalog.pg_class cls
ON cls.oid = cns.conrelid
JOIN pg_catalog.pg_namespace ns
ON ns.oid = cls.relnamespace
WHERE cns.contype = 'c'
AND   ns.nspname = $1
AND   cls.relname = $2
UNION
SELECT ns.nspname AS schema
     , cls.relname AS table_name
     , cns.conname AS constraint_name
     , 'EXCLUDE' AS constraint_kind
     , array_to_string(array_agg(att.attname || ' WITH ' || op.oprname), ', ') AS constraint_content
FROM (
    SELECT conrelid
         , conname
         , conkey AS colnums
         , conexclop AS opids
         , generate_series(1, length(array_to_string(conkey, ' ')) - length(array_to_string(conkey, '')) + 1) AS pos
    FROM pg_catalog.pg_constraint
    WHERE contype = 'x'
) cns
JOIN pg_catalog.pg_class cls
ON cls.oid = cns.conrelid
JOIN pg_catalog.pg_namespace ns
ON ns.oid = cls.relnamespace
JOIN pg_catalog.pg_attribute att
ON att.attrelid = cls.oid
AND att.attnum = cns.colnums[pos]
JOIN pg_catalog.pg_operator op
ON op.oid = cns.opids[pos]
WHERE ns.nspname = $1
AND   cls.relname = $2
GROUP BY 1, 2, 3
UNION
SELECT ns.nspname AS schema
     , cls.relname AS table_name
     , cns.conname AS constraint_name
     , 'UNIQUE' AS constraint_kind
     , array_to_string(array_agg(att.attname), ', ') AS constraint_content
FROM (
    SELECT conrelid
         , conname
         , conkey AS colnums
         , generate_series(1, length(array_to_string(conkey, ' ')) - length(array_to_string(conkey, '')) + 1) AS pos
    FROM pg_catalog.pg_constraint
    WHERE contype = 'u'
) cns
JOIN pg_catalog.pg_class cls
ON cls.oid = cns.conrelid
JOIN pg_catalog.pg_namespace ns
ON ns.oid = cls.relnamespace
JOIN pg_catalog.pg_attribute att
ON att.attrelid = cls.oid
AND att.attnum = cns.colnums[cns.pos]
WHERE ns.nspname = $1
AND   cls.relname = $2
GROUP BY 1, 2, 3
ORDER BY table_name, constraint_kind, constraint_name`
}

func (p postgres) AllConstrantsSQL() string {
	return `
SELECT ns.nspname AS schema
     , cls.relname AS table_name
     , cns.conname AS constraint_name
     , 'CHECK' AS constraint_kind
     , cns.consrc AS constraint_content
FROM pg_catalog.pg_constraint cns
JOIN pg_catalog.pg_class cls
ON cls.oid = cns.conrelid
JOIN pg_catalog.pg_namespace ns
ON ns.oid = cls.relnamespace
WHERE cns.contype = 'c'
AND   ns.nspname = $1
UNION
SELECT ns.nspname AS schema
     , cls.relname AS table_name
     , cns.conname AS constraint_name
     , 'EXCLUDE' AS constraint_kind
     , array_to_string(array_agg(att.attname || ' WITH ' || op.oprname), ', ') AS constraint_content
FROM (
    SELECT conrelid
         , conname
         , conkey AS colnums
         , conexclop AS opids
         , generate_series(1, length(array_to_string(conkey, ' ')) - length(array_to_string(conkey, '')) + 1) AS pos
    FROM pg_catalog.pg_constraint
    WHERE contype = 'x'
) cns
JOIN pg_catalog.pg_class cls
ON cls.oid = cns.conrelid
JOIN pg_catalog.pg_namespace ns
ON ns.oid = cls.relnamespace
JOIN pg_catalog.pg_attribute att
ON att.attrelid = cls.oid
AND att.attnum = cns.colnums[cns.pos]
JOIN pg_catalog.pg_operator op
ON op.oid = cns.opids[pos]
WHERE ns.nspname = $1
GROUP BY 1, 2, 3
UNION
SELECT ns.nspname AS schema
     , cls.relname AS table_name
     , cns.conname AS constraint_name
     , 'UNIQUE' AS constraint_kind
     , array_to_string(array_agg(att.attname), ', ') AS constraint_content
FROM (
    SELECT conrelid
         , conname
         , conkey AS colnums
         , generate_series(1, length(array_to_string(conkey, ' ')) - length(array_to_string(conkey, '')) + 1) AS pos
    FROM pg_catalog.pg_constraint
    WHERE contype = 'u'
) cns
JOIN pg_catalog.pg_class cls
ON cls.oid = cns.conrelid
JOIN pg_catalog.pg_namespace ns
ON ns.oid = cls.relnamespace
JOIN pg_catalog.pg_attribute att
ON att.attrelid = cls.oid
AND att.attnum = cns.colnums[cns.pos]
WHERE ns.nspname = $1
GROUP BY 1, 2, 3
ORDER BY table_name, constraint_kind, constraint_name`
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
