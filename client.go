package dbmodel

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	// ErrSchemaEmpty is raised when schema is not given.
	ErrSchemaEmpty = errors.New("Schema is required")
	// ErrConnNotFound is raised when call function before connect to database.
	ErrConnNotFound = errors.New("Database connection is not found")
	// ErrInvalidDriver is raised when given driver is unknown.
	ErrInvalidDriver = errors.New("Invalid driver")
	// ErrTableNameEmpty is raised when table name is not given.
	ErrTableNameEmpty = errors.New("Table name is required.")
)

// Client is table meta data loding client.
type Client struct {
	dataSource DataSource
	provider   Provider
	db         *sql.DB
	err        error
}

// NewClient returns new Client for connecting to given data source.
// Acceptable driver names are 'postgres'.
func NewClient(driver string, ds DataSource) *Client {
	p, err := findProvider(driver)
	return &Client{
		dataSource: ds,
		provider:   p,
		err:        err,
	}
}

// SetProvider sets custom provider.
// If use custom provider, call this before Connect.
func (c *Client) SetProvider(p Provider) {
	c.provider = p
	c.err = nil
}

// Connect to database.
func (c *Client) Connect() {
	if c.err != nil {
		return
	}

	c.db, c.err = c.provider.Connect(c.dataSource)
	if c.err != nil {
		c.db = nil
	}
}

// Disconnect from datasource and close database connection.
func (c *Client) Disconnect() error {
	defer func() {
		c.db = nil
	}()

	if c.db != nil {
		return c.db.Close()
	}
	return c.err
}

// AllTableNames returns all table names in given schema.
// If schema is empty, raise ErrSchemaEmpty.
func (c *Client) AllTableNames(schema string) ([]*Table, error) {
	if err := c.preCheck(schema); err != nil {
		return nil, err
	}

	rows, err := c.db.Query(c.provider.AllTableNamesSQL(), schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return c.readTableNames(rows), nil
}

// TableNames returns table names in given schema.
// If schema is empty, raise ErrSchemaEmpty.
// If name is empaty, TableNames returns all table names orderd by table names.
// If name is given, TableNames returns table names that matches given name.
func (c *Client) TableNames(schema string, name string) ([]*Table, error) {
	if err := c.preCheck(schema); err != nil {
		return nil, err
	}

	rows, err := c.db.Query(c.provider.TableNamesSQL(), schema, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return c.readTableNames(rows), nil
}

// Table returns table meta data.
// Return value contains table name and column meta data.
// If schema is empty, raise ErrSchemaEmpty.
// If name is empaty, raise ErrTableNameEmpty.
func (c *Client) Table(schema string, name string) (*Table, error) {
	if err := c.preCheck(schema); err != nil {
		return nil, err
	}
	if name == "" {
		return nil, ErrTableNameEmpty
	}

	rows, err := c.db.Query(c.provider.TableSQL(), schema, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tbls := c.readTables(rows)
	if len(tbls) == 0 {
		return nil, fmt.Errorf("Table '%v' is not found.", name)
	}
	tbl := tbls[0]
	c.appendIndices(tbl)
	c.appendForeignKyes(tbl)
	c.appendReferencedKyes(tbl)
	return tbl, nil
}

func (c *Client) preCheck(schema string) error {
	if c.err != nil {
		return c.err
	}
	if schema == "" {
		return ErrSchemaEmpty
	}
	if c.db == nil {
		return ErrConnNotFound
	}

	return nil
}

func findProvider(driver string) (Provider, error) {
	if driver == "postgres" {
		return postgres{}, nil
	}
	return nil, ErrInvalidDriver
}

func (c *Client) readTableNames(rows *sql.Rows) []*Table {
	tables := make([]*Table, 0, 10)
	for rows.Next() {
		var (
			schema  sql.NullString
			name    sql.NullString
			comment sql.NullString
		)
		rows.Scan(&schema, &name, &comment)
		t := NewTable(schema.String, name.String, comment.String)
		tables = append(tables, &t)
	}
	return tables
}

func (c *Client) readTables(rows *sql.Rows) []*Table {
	tbls := make([]*Table, 0, 10)
	for rows.Next() {
		var (
			schema       sql.NullString
			tblName      sql.NullString
			tblComment   sql.NullString
			colName      sql.NullString
			colComment   sql.NullString
			dataType     sql.NullString
			length       sql.NullInt64
			precision    sql.NullInt64
			scale        sql.NullInt64
			nullable     sql.NullString
			defaultValue sql.NullString
			pkPosition   sql.NullInt64
		)

		rows.Scan(&schema, &tblName, &tblComment, &colName, &colComment, &dataType, &length, &precision, &scale, &nullable, &defaultValue, &pkPosition)
		if len(tbls) == 0 || tbls[len(tbls)-1].Name() != tblName.String {
			tbl := NewTable(schema.String, tblName.String, tblComment.String)
			tbls = append(tbls, &tbl)
		}
		col := NewColumn(
			schema.String,
			tblName.String,
			colName.String,
			colComment.String,
			dataType.String,
			NewSize(length, precision, scale),
			nullable.String == "YES",
			defaultValue.String,
			pkPosition.Int64)
		tbls[len(tbls)-1].AddColumn(&col)
	}
	return tbls
}

func (c *Client) appendIndices(tbl *Table) error {
	rows, err := c.db.Query(c.provider.IndicesSQL(), tbl.Schema(), tbl.Name())
	if err != nil {
		return err
	}
	for rows.Next() {
		var (
			schema  sql.NullString
			tblName sql.NullString
			name    sql.NullString
			uniq    sql.NullString
			colName sql.NullString
		)
		rows.Scan(&schema, &tblName, &name, &uniq, &colName)
		if len(tbl.Indices()) == 0 || tbl.lastIndex().Name() != name.String {
			idx := NewIndex(schema.String, tblName.String, name.String, uniq.String == "YES")
			tbl.AddIndex(&idx)
		}
		col, err := tbl.FindColumn(colName.String)
		if err != nil {
			return err
		}
		tbl.lastIndex().AddColumn(col)
	}
	return nil
}

func (c *Client) appendForeignKyes(tbl *Table) error {
	rows, err := c.db.Query(c.provider.ForeignKeysSQL(), tbl.Schema(), tbl.Name())
	if err != nil {
		return err
	}
	for rows.Next() {
		var (
			name     sql.NullString
			schema   sql.NullString
			tblName  sql.NullString
			colName  sql.NullString
			fSchema  sql.NullString
			fTblName sql.NullString
			fColName sql.NullString
		)
		rows.Scan(&name, &schema, &tblName, &colName, &fSchema, &fTblName, &fColName)
		if len(tbl.ForeignKeys()) == 0 || tbl.lastForeignKey().Name() != name.String {
			fk := NewForeignKey(schema.String, tblName.String, name.String)
			tbl.AddForeignKey(&fk)
		}
		col, err := tbl.FindColumn(colName.String)
		if err != nil {
			return err
		}
		fCol := &Column{
			schema:    fSchema.String,
			tableName: fTblName.String,
			name:      fColName.String,
		}
		cr := NewColumnReference(col, fCol)
		tbl.lastForeignKey().AddColumnReference(&cr)
	}
	return nil
}

func (c *Client) appendReferencedKyes(tbl *Table) error {
	rows, err := c.db.Query(c.provider.ReferencedKeysSQL(), tbl.Schema(), tbl.Name())
	if err != nil {
		return err
	}
	for rows.Next() {
		var (
			name     sql.NullString
			schema   sql.NullString
			tblName  sql.NullString
			colName  sql.NullString
			fSchema  sql.NullString
			fTblName sql.NullString
			fColName sql.NullString
		)
		rows.Scan(&name, &schema, &tblName, &colName, &fSchema, &fTblName, &fColName)
		if len(tbl.ReferencedKeys()) == 0 || tbl.lastRefKey().Name() != name.String {
			fk := NewForeignKey(schema.String, tblName.String, name.String)
			tbl.AddReferencedKey(&fk)
		}
		col := &Column{
			schema:    schema.String,
			tableName: tblName.String,
			name:      colName.String,
		}
		fCol, err := tbl.FindColumn(fColName.String)
		if err != nil {
			return err
		}
		cr := NewColumnReference(col, fCol)
		tbl.lastRefKey().AddColumnReference(&cr)
	}
	return nil
}
