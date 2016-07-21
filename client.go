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

	return readTableNames(schema, rows), nil
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

	return readTableNames(schema, rows), nil
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

	tables := readTables(schema, rows)
	if len(tables) == 0 {
		return nil, fmt.Errorf("Table '%v' is not found.", name)
	}
	return tables[0], nil
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

func readTableNames(schema string, rows *sql.Rows) []*Table {
	tables := make([]*Table, 0, 10)
	for rows.Next() {
		var (
			name    sql.NullString
			comment sql.NullString
		)
		rows.Scan(&name, &comment)
		t := NewTable(schema, name.String, comment.String)
		tables = append(tables, &t)
	}
	return tables
}

func readTables(schema string, rows *sql.Rows) []*Table {
	tables := make([]*Table, 0, 10)
	for rows.Next() {
		var (
			tableName     sql.NullString
			tableComment  sql.NullString
			columnName    sql.NullString
			columnComment sql.NullString
			dataType      sql.NullString
			length        sql.NullInt64
			precision     sql.NullInt64
			scale         sql.NullInt64
			nullable      sql.NullString
			defaultValue  sql.NullString
			pkPosition    sql.NullInt64
		)

		rows.Scan(&tableName, &tableComment, &columnName, &columnComment, &dataType, &length, &precision, &scale, &nullable, &defaultValue, &pkPosition)
		if len(tables) == 0 || tables[len(tables)-1].Name() != tableName.String {
			t := NewTable(schema, tableName.String, tableComment.String)
			tables = append(tables, &t)
		}
		c := NewColumn(
			schema,
			tableName.String,
			columnName.String,
			columnComment.String,
			dataType.String,
			NewSize(length, precision, scale),
			nullable.String == "YES",
			defaultValue.String,
			pkPosition.Int64)
		tables[len(tables)-1].AddColumn(&c)
	}
	return tables
}
