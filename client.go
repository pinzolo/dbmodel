package dbmodel

import (
	"database/sql"
	"errors"
)

var (
	// ErrSchemaEmpty is raised when schema no given.
	ErrSchemaEmpty = errors.New("Schema is required")
	// ErrConnNotFound is raised when call function before connect to database.
	ErrConnNotFound = errors.New("Database connection is not found")
	// ErrInvalidDriver is raised when given driver is unknown.
	ErrInvalidDriver = errors.New("Invalid driver")
)

// Client is table meta data loding client.
type Client struct {
	dataSource DataSource
	provider   Provider
	db         *sql.DB
	err        error
}

// NewClient returns new Client for connecting to given data source.
// 'postgres' is acceptable as driver.
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
// If schema is empty, raise error.
func (c *Client) AllTableNames(schema string) ([]*Table, error) {
	if err := c.preCheck(schema); err != nil {
		return nil, err
	}

	rows, err := c.db.Query(c.provider.AllTableNamesSQL(), schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return readTables(rows), nil
}

// TableNames returns table names in given schema.
// If schema is empty, raise error.
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

	return readTables(rows), nil
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

func readTables(rows *sql.Rows) []*Table {
	tables := make([]*Table, 0, 10)
	for rows.Next() {
		var name string
		var comment string
		rows.Scan(&name, &comment)
		t := NewTable(name, comment)
		tables = append(tables, &t)
	}
	return tables
}
