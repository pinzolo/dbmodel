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
func NewClient(ds DataSource) *Client {
	p, err := findProvider(ds)
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

	c.db, c.err = c.provider.Connect()
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
// If schema is empty, raise ErrSchemaEmpty.
// If name is empaty, raise ErrTableNameEmpty.
func (c *Client) Table(schema string, name string, opt Option) (*Table, error) {
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
	if opt.Indices {
		idxs, err := c.loadIndices(tbl.Schema(), tbl.Name())
		if err != nil {
			return nil, err
		}
		tbl.indices = idxs
	}
	if opt.ForeignKeys {
		fks, err := c.loadForeignKeys(tbl.Schema(), tbl.Name())
		if err != nil {
			return nil, err
		}
		tbl.foreignKeys = fks
	}
	if opt.ReferencedKeys {
		rks, err := c.loadReferencedKeys(tbl.Schema(), tbl.Name())
		if err != nil {
			return nil, err
		}
		tbl.refKeys = rks
	}
	if opt.Constraints {
		cnss, err := c.loadConstraints(tbl.Schema(), tbl.Name())
		if err != nil {
			return nil, err
		}
		tbl.constraints = cnss
	}
	return tbl, nil
}

// AllTables returns table meta data list that are contained in given schema.
// If schema is empty, raise ErrSchemaEmpty.
func (c *Client) AllTables(schema string, opt Option) ([]*Table, error) {
	if err := c.preCheck(schema); err != nil {
		return nil, err
	}

	rows, err := c.db.Query(c.provider.AllTablesSQL(), schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tbls := c.readTables(rows)
	tblMap := make(map[string]*Table)
	for _, tbl := range tbls {
		tblMap[tbl.Name()] = tbl
	}
	if opt.Indices {
		idxMap, err := c.loadIndicesMap(schema)
		if err != nil {
			return nil, err
		}
		for tblName, idxs := range idxMap {
			tbl, ok := tblMap[tblName]
			if ok {
				tbl.indices = idxs
			}
		}
	}
	if opt.ForeignKeys {
		fkMap, err := c.loadForeignKeysMap(schema)
		if err != nil {
			return nil, err
		}
		for tblName, fks := range fkMap {
			tbl, ok := tblMap[tblName]
			if ok {
				tbl.foreignKeys = fks
			}
		}
	}
	if opt.ReferencedKeys {
		rkMap, err := c.loadReferencedKeysMap(schema)
		if err != nil {
			return nil, err
		}
		for tblName, rks := range rkMap {
			tbl, ok := tblMap[tblName]
			if ok {
				tbl.refKeys = rks
			}
		}
	}
	if opt.Constraints {
		cnsMap, err := c.loadConstraintsMap(schema)
		if err != nil {
			return nil, err
		}
		for tblName, cnss := range cnsMap {
			tbl, ok := tblMap[tblName]
			if ok {
				tbl.constraints = cnss
			}
		}
	}
	return tbls, nil
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

func findProvider(ds DataSource) (Provider, error) {
	if ds.Driver == "postgres" {
		return postgres{ds: ds}, nil
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

func (c *Client) loadIndices(schema string, tblName string) ([]*Index, error) {
	rows, err := c.db.Query(c.provider.IndicesSQL(), schema, tblName)
	if err != nil {
		return nil, err
	}

	return c.readIndices(rows), nil
}

func (c *Client) loadIndicesMap(schema string) (map[string][]*Index, error) {
	rows, err := c.db.Query(c.provider.AllIndicesSQL(), schema)
	if err != nil {
		return nil, err
	}

	idxMap := make(map[string][]*Index)
	for _, idx := range c.readIndices(rows) {
		idxs, ok := idxMap[idx.TableName()]
		if ok {
			idxMap[idx.TableName()] = append(idxs, idx)
		} else {
			idxMap[idx.TableName()] = []*Index{idx}
		}
	}
	return idxMap, nil
}

func (c *Client) readIndices(rows *sql.Rows) []*Index {
	idxs := make([]*Index, 0, 10)
	for rows.Next() {
		var (
			schema  sql.NullString
			tblName sql.NullString
			name    sql.NullString
			uniq    sql.NullString
			colName sql.NullString
		)
		rows.Scan(&schema, &tblName, &name, &uniq, &colName)
		if len(idxs) == 0 || idxs[len(idxs)-1].Name() != name.String {
			idx := NewIndex(schema.String, tblName.String, name.String, uniq.String == "YES")
			idxs = append(idxs, &idx)
		}
		col := &Column{
			schema:    schema.String,
			tableName: tblName.String,
			name:      colName.String,
		}
		idxs[len(idxs)-1].AddColumn(col)
	}
	return idxs
}

func (c *Client) loadConstraints(schema string, tblName string) ([]*Constraint, error) {
	rows, err := c.db.Query(c.provider.ConstraintsSQL(), schema, tblName)
	if err != nil {
		return nil, err
	}

	return c.readConstraints(rows), nil
}

func (c *Client) loadConstraintsMap(schema string) (map[string][]*Constraint, error) {
	rows, err := c.db.Query(c.provider.AllConstraintsSQL(), schema)
	if err != nil {
		return nil, err
	}

	cnsMap := make(map[string][]*Constraint)
	for _, cns := range c.readConstraints(rows) {
		fks, ok := cnsMap[cns.TableName()]
		if ok {
			cnsMap[cns.TableName()] = append(fks, cns)
		} else {
			cnsMap[cns.TableName()] = []*Constraint{cns}
		}
	}
	return cnsMap, nil
}

func (c *Client) readConstraints(rows *sql.Rows) []*Constraint {
	cons := make([]*Constraint, 0, 10)
	for rows.Next() {
		var (
			schema  sql.NullString
			tblName sql.NullString
			name    sql.NullString
			kind    sql.NullString
			content sql.NullString
		)
		rows.Scan(&schema, &tblName, &name, &kind, &content)
		if len(cons) == 0 || cons[len(cons)-1].Name() != name.String {
			con := NewConstraint(schema.String, tblName.String, name.String, kind.String, content.String)
			cons = append(cons, &con)
		}
	}
	return cons
}

func (c *Client) loadForeignKeys(schema string, tblName string) ([]*ForeignKey, error) {
	rows, err := c.db.Query(c.provider.ForeignKeysSQL(), schema, tblName)
	if err != nil {
		return nil, err
	}

	return c.readForeignKeys(rows), nil
}

func (c *Client) loadForeignKeysMap(schema string) (map[string][]*ForeignKey, error) {
	rows, err := c.db.Query(c.provider.AllForeignKeysSQL(), schema)
	if err != nil {
		return nil, err
	}

	fkMap := make(map[string][]*ForeignKey)
	for _, fk := range c.readForeignKeys(rows) {
		fks, ok := fkMap[fk.TableName()]
		if ok {
			fkMap[fk.TableName()] = append(fks, fk)
		} else {
			fkMap[fk.TableName()] = []*ForeignKey{fk}
		}
	}
	return fkMap, nil
}

func (c *Client) readForeignKeys(rows *sql.Rows) []*ForeignKey {
	fks := make([]*ForeignKey, 0, 10)
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
		if len(fks) == 0 || fks[len(fks)-1].Name() != name.String {
			fk := NewForeignKey(schema.String, tblName.String, name.String)
			fks = append(fks, &fk)
		}
		col := &Column{
			schema:    schema.String,
			tableName: tblName.String,
			name:      colName.String,
		}
		fCol := &Column{
			schema:    fSchema.String,
			tableName: fTblName.String,
			name:      fColName.String,
		}
		cr := NewColumnReference(col, fCol)
		fks[len(fks)-1].AddColumnReference(&cr)
	}
	return fks
}

func (c *Client) loadReferencedKeys(schema string, tblName string) ([]*ForeignKey, error) {
	rows, err := c.db.Query(c.provider.ReferencedKeysSQL(), schema, tblName)
	if err != nil {
		return nil, err
	}

	return c.readForeignKeys(rows), nil
}

func (c *Client) loadReferencedKeysMap(schema string) (map[string][]*ForeignKey, error) {
	rows, err := c.db.Query(c.provider.AllReferencedKeysSQL(), schema)
	if err != nil {
		return nil, err
	}

	rkMap := make(map[string][]*ForeignKey)
	for _, rk := range c.readForeignKeys(rows) {
		tblName := rk.ColumnReferences()[0].To().TableName()
		rks, ok := rkMap[tblName]
		if ok {
			rkMap[tblName] = append(rks, rk)
		} else {
			rkMap[tblName] = []*ForeignKey{rk}
		}
	}
	return rkMap, nil
}
