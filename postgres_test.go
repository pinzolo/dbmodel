package dbmodel

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestPostgresConnectionString(t *testing.T) {
	p := newPostgres(InitDataSource())
	excepted := ""
	if p.connStr() != excepted {
		t.Errorf("ConnectionString was %v, but %v was expected.", excepted, p.connStr())
	}
	p.ds.User = "postgres"
	excepted = "user=postgres"
	if p.connStr() != excepted {
		t.Errorf("ConnectionString was %v, but %v was expected.", excepted, p.connStr())
	}
	p.ds.Password = "12345"
	excepted = "user=postgres password=12345"
	if p.connStr() != excepted {
		t.Errorf("ConnectionString was %v, but %v was expected.", excepted, p.connStr())
	}
	p.ds.Host = "localhost"
	excepted = "host=localhost user=postgres password=12345"
	if p.connStr() != excepted {
		t.Errorf("ConnectionString was %v, but %v was expected.", excepted, p.connStr())
	}
	p.ds.Port = 5432
	excepted = "host=localhost port=5432 user=postgres password=12345"
	if p.connStr() != excepted {
		t.Errorf("ConnectionString was %v, but %v was expected.", excepted, p.connStr())
	}
	p.ds.Database = "sample"
	excepted = "host=localhost port=5432 user=postgres password=12345 dbname=sample"
	if p.connStr() != excepted {
		t.Errorf("ConnectionString was %v, but %v was expected.", excepted, p.connStr())
	}
	p.ds.Options["sslmode"] = "disable"
	excepted = "host=localhost port=5432 user=postgres password=12345 dbname=sample sslmode=disable"
	if p.connStr() != excepted {
		t.Errorf("ConnectionString was %v, but %v was expected.", excepted, p.connStr())
	}
}

func TestPostgresAllTableNames(t *testing.T) {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	ts, err := c.AllTableNames("schm")
	if err != nil {
		t.Error(err)
	}
	if len(ts) != 3 {
		t.Errorf("AllTableNames should return 3 table names. but actual %v", len(ts))
		return
	}
	if ts[0].Name() != "tbl1" {
		t.Errorf("AllTableNames returns invalid table name. expected 'tbl1', but actual '%v'", ts[0].Name())
	}
	if ts[0].Comment() != "" {
		t.Errorf("Table comment is null, Comment() should return empty")
	}
	if ts[1].Name() != "tbl2" {
		t.Errorf("AllTableNames returns invalid table name. expected 'tbl2', but actual '%v'", ts[1].Name())
	}
	if ts[1].Comment() != "This is table2" {
		t.Errorf("AllTableNames should pick up table comment. %#v", ts[1])
	}
}

func TestPostgresAllTableNamesOtherSchema(t *testing.T) {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	ts, err := c.AllTableNames("other")
	if err != nil {
		t.Error(err)
	}
	if len(ts) != 1 {
		t.Errorf("AllTableNames should return 1 table name. but actual %v", len(ts))
		return
	}
	if ts[0].Name() != "tbl_other" {
		t.Errorf("AllTableNames returns invalid table name. expected 'tbl_other', but actual '%v'", ts[0].Name())
	}
}

func TestPostgresTableNames(t *testing.T) {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	ts, err := c.TableNames("schm", "tbl1")
	if err != nil {
		t.Error(err)
	}
	if len(ts) != 1 {
		t.Errorf("TableNames with tbl1 should return 1 table name. but actual %v", len(ts))
		return
	}
	if ts[0].Name() != "tbl1" {
		t.Errorf("TableNames returns invalid table name. expected 'tbl1', but actual '%v'", ts[0].Name())
	}
	if ts[0].Comment() != "" {
		t.Errorf("Table comment is null, Comment() should return empty")
	}
}

func TestPostgresTableNamesNoResultOnInvalidSchema(t *testing.T) {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	ts, err := c.TableNames("other", "tbl1")
	if err != nil {
		t.Error(err)
	}
	if len(ts) != 0 {
		t.Errorf("TableNames should return 0 table name. but actual %v", len(ts))
		return
	}
}

func TestPostgresTableNamesNoResultOnInvalidTalbeName(t *testing.T) {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	ts, err := c.TableNames("schm", "sample")
	if err != nil {
		t.Error(err)
	}
	if len(ts) != 0 {
		t.Errorf("TableNames should return 0 table name. but actual %v", len(ts))
		return
	}
}

func TestPostgresAllTableNamesWithoutSchema(t *testing.T) {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	_, err := c.AllTableNames("")
	if err == nil {
		t.Errorf("Client should raise error when empty schema given.")
	}
	if err != ErrSchemaEmpty {
		t.Errorf("%v is invalid Error", err)
	}
}

func TestPostgresTableNamesWithoutSchema(t *testing.T) {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	_, err := c.TableNames("", "tbl1")
	if err == nil {
		t.Errorf("Client should raise error when empty schema given.")
	}
	if err != ErrSchemaEmpty {
		t.Errorf("%v is invalid Error", err)
	}
}

func TestPostgresTableValid(t *testing.T) {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	tbl, err := c.Table("schm", "tbl1", RequireNone)
	if err != nil {
		t.Error("Client should not raise error when valid schema and table name given.")
	}
	if tbl.Name() != "tbl1" {
		t.Errorf("Table name is invalid. expected: %v, actual: %v", "tbl1", tbl.Name())
	}
}

func TestPostgresTableNotFound(t *testing.T) {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	_, err := c.Table("schm", "xxxxx", RequireNone)
	if err == nil {
		t.Error("Client should raise error when given table name not exist")
	}
	if !strings.Contains(err.Error(), "xxxxx") {
		t.Error("Error message should contains given table name.")
	}
}

func TestPostgresTableWithoutSchema(t *testing.T) {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	_, err := c.Table("", "tbl1", RequireNone)
	if err == nil {
		t.Errorf("Client should raise error when empty schema given.")
	}
	if err != ErrSchemaEmpty {
		t.Errorf("%v is invalid Error", err)
	}
}

func TestPostgresTableColumnsCount(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl1")
	for _, tbl := range tbls {
		if len(tbl.Columns()) != 7 {
			t.Errorf("Column count is invalid. expected: %v, actual: %v", 7, len(tbl.Columns()))
		}
	}
}

func TestPostgresTableColumnsOrder(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl2")
	for _, tbl := range tbls {
		colNames := []string{"id", "tbl1_id", "idx_key", "chk"}
		for i, expected := range colNames {
			if actual := tbl.Columns()[i].Name(); actual != expected {
				t.Errorf("Column order is invalid. expected: %v, actual: %v", expected, actual)
			}
		}
	}
}

func TestPostgresTableColumnComment(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl2")
	for _, tbl := range tbls {
		comments := []string{"ID", "tbl1ID", "", ""}
		for i, expected := range comments {
			if actual := tbl.Columns()[i].Comment(); actual != expected {
				t.Errorf("Cannot get valid comment. expected: %v, actual: %v", expected, actual)
			}
		}
	}
}

func TestPostgresTableColumnDataType(t *testing.T) {
	tbls := loadPostgresTableBy2Way("other", "tbl_other")
	for _, tbl := range tbls {
		dts := []string{"int4", "schm.domain1"}
		for i, expected := range dts {
			if actual := tbl.Columns()[i].DataType(); actual != expected {
				t.Errorf("Cannot get valid data type. expected: %v, actual: %v", expected, actual)
			}
		}
	}
}

func TestPostgresTableColumnSize(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl1")
	for _, tbl := range tbls {
		textSize := tbl.Columns()[2].Size()
		if !textSize.IsValid() || !textSize.Length().Valid || textSize.Precision().Valid {
			t.Error("Cannot get valid text size.")
		}
		if actual, expected := textSize.String(), "50"; actual != expected {
			t.Errorf("Text size value is invalid. expected: %v, actual: %v", expected, actual)
		}

		nullSize := tbl.Columns()[1].Size()
		if nullSize.IsValid() {
			t.Error("Cannot get valid null size.")
		}
		if actual, expected := nullSize.String(), ""; actual != expected {
			t.Errorf("Null size value is invalid. expected: %v, actual: %v", expected, actual)
		}

		intSize := tbl.Columns()[4].Size()
		if !intSize.IsValid() || intSize.Length().Valid || !intSize.Precision().Valid || !intSize.Scale().Valid {
			t.Error("Cannot get valid integer size.")
		}
		if actual, expected := intSize.String(), "8, 2"; actual != expected {
			t.Errorf("Integer size value is invalid. expected: %v, actual: %v", expected, actual)
		}

		dateSize := tbl.Columns()[6].Size()
		if !dateSize.IsValid() || dateSize.Length().Valid || !dateSize.Precision().Valid {
			t.Error("Cannot get valid date size.")
		}
		if actual, expected := dateSize.String(), "6"; actual != expected {
			t.Errorf("Date size value is invalid. expected: %v, actual: %v", expected, actual)
		}
	}
}

func TestPostgresTableColumnNullable(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl1")
	for _, tbl := range tbls {
		if tbl.Columns()[0].IsNullable() {
			t.Errorf("Column '%v' is not nullable, but IsNullable() returns true", tbl.Columns()[0].Name())
		}
		if !tbl.Columns()[3].IsNullable() {
			t.Errorf("Column '%v' is nullable, but IsNullable() returns false", tbl.Columns()[3].Name())
		}
	}
}

func TestPostgresTableColumnDefaultValue(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl1")
	for _, tbl := range tbls {
		if actual := tbl.Columns()[1].DefaultValue(); actual != "" {
			t.Errorf("Column '%v' do not have default value, but DefaultValue() returns %v", tbl.Columns()[1].Name(), actual)
		}
		if actual, expected := tbl.Columns()[5].DefaultValue(), "1"; actual != expected {
			t.Errorf("Cannot get invalid default value of '%v'. expected: %v, actual: %v", tbl.Columns()[2].Name(), expected, actual)
		}
		if actual, expected := tbl.Columns()[6].DefaultValue(), "now()"; actual != expected {
			t.Errorf("Cannot get invalid default value of '%v'. expected: %v, actual: %v", tbl.Columns()[4].Name(), expected, actual)
		}
	}
}

func TestPostgresTableColumnPrimaryKeyPosition(t *testing.T) {
	// Single column primary key
	tbls := loadPostgresTableBy2Way("schm", "tbl1")
	for _, tbl := range tbls {
		if actual, expected := tbl.Columns()[0].PrimaryKeyPosition(), int64(1); actual != expected {
			t.Errorf("Cannot get invalid primary key position of '%v'. expected: %v, actual: %v", tbl.Columns()[0].Name(), expected, actual)
		}
		if actual, expected := tbl.Columns()[1].PrimaryKeyPosition(), int64(0); actual != expected {
			t.Errorf("Cannot get invalid primary key position of '%v'. expected: %v, actual: %v", tbl.Columns()[0].Name(), expected, actual)
		}
	}

	// Multi columns primary key
	tbls = loadPostgresTableBy2Way("schm", "tbl3")
	for _, tbl := range tbls {
		if actual, expected := tbl.Columns()[0].PrimaryKeyPosition(), int64(1); actual != expected {
			t.Errorf("Cannot get invalid primary key position of '%v'. expected: %v, actual: %v", tbl.Columns()[0].Name(), expected, actual)
		}
		if actual, expected := tbl.Columns()[1].PrimaryKeyPosition(), int64(2); actual != expected {
			t.Errorf("Cannot get invalid primary key position of '%v'. expected: %v, actual: %v", tbl.Columns()[1].Name(), expected, actual)
		}
	}
}

func TestPostgresTableIndicesCount(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl2")
	for _, tbl := range tbls {
		if actual, expected := len(tbl.Indices()), 3; actual != expected {
			t.Errorf("Index count is invalid. expected: %v, actual: %v", expected, actual)
		}
	}
}

func TestPostgresTableIndicesOrder(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl2")
	for _, tbl := range tbls {
		idxNames := []string{"tbl2_ex1", "tbl2_idx1", "tbl2_pkey"}
		for i, expected := range idxNames {
			if actual := tbl.Indices()[i].Name(); actual != expected {
				t.Errorf("Index order is invalid. expected: %v, actual: %v", expected, actual)
			}
		}
	}
}

func TestPostgresTableIndexIsUnique(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl2")
	for _, tbl := range tbls {
		if tbl.Indices()[0].IsUnique() {
			t.Errorf("Index '%v' is unique, but IsUnique() returns true", tbl.Indices()[0].Name())
		}
		if tbl.Indices()[1].IsUnique() {
			t.Errorf("Index '%v' is not unique, but IsUnique() returns false", tbl.Indices()[1].Name())
		}
		if !tbl.Indices()[2].IsUnique() {
			t.Errorf("Index '%v' is not unique, but IsUnique() returns false", tbl.Indices()[1].Name())
		}
	}
}

func TestPostgresTableIndexColumns(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl3")
	for _, tbl := range tbls {
		// count
		if actual, expected := len(tbl.Indices()[0].Columns()), 1; actual != expected {
			t.Errorf("Index '%v' should have 1 column, but actually have %v columns.", tbl.Indices()[0].Name(), actual)
		}
		if actual, expected := len(tbl.Indices()[1].Columns()), 2; actual != expected {
			t.Errorf("Index '%v' should have 1 column, but actually have %v columns.", tbl.Indices()[1].Name(), actual)
		}
		// order
		if actual, expected := tbl.Indices()[1].Columns()[0].Name(), "tbl1_id"; actual != expected {
			t.Errorf("Index column order is invalid. expected: %v, actual: %v", expected, actual)
		}
		if actual, expected := tbl.Indices()[1].Columns()[1].Name(), "tbl2_id"; actual != expected {
			t.Errorf("Index column order is invalid. expected: %v, actual: %v", expected, actual)
		}
	}
}

func TestPostgresTableForeignKeysCount(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl3")
	for _, tbl := range tbls {
		if actual, expected := len(tbl.ForeignKeys()), 2; actual != expected {
			t.Errorf("Foreign key count is invalid. expected: %v, actual: %v", expected, actual)
		}
	}
}

func TestPostgresTableForeignKeysOrder(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl3")
	for _, tbl := range tbls {
		fkNames := []string{"tbl3_fk1", "tbl3_fk2"}
		for i, expected := range fkNames {
			if actual := tbl.ForeignKeys()[i].Name(); actual != expected {
				t.Errorf("Foreign key order is invalid. expected: %v, actual: %v", expected, actual)
			}
		}
	}
}

func TestPostgresTableForeignKeysColumnCount(t *testing.T) {
	tbls := loadPostgresTableBy2Way("other", "tbl_other")
	for _, tbl := range tbls {
		cnts := []int{1, 1, 2}
		for i, expected := range cnts {
			if actual := len(tbl.ForeignKeys()[i].ColumnReferences()); actual != expected {
				t.Errorf("Foreign key's column count is invalid. expected: %v, actual: %v", expected, actual)
			}
		}
	}
}

func TestPostgresTableForeignKeysColumnReferences(t *testing.T) {
	tbls := loadPostgresTableBy2Way("other", "tbl_other")
	for _, tbl := range tbls {
		if actual, expected := colRefToString(tbl.ForeignKeys()[0].ColumnReferences()[0]), "other.tbl_other.tbl1_id -> schm.tbl1.id"; actual != expected {
			t.Errorf("Foreign key's column reference is invalid. expected: %v, actual: %v", expected, actual)
		}
		if actual, expected := colRefToString(tbl.ForeignKeys()[2].ColumnReferences()[0]), "other.tbl_other.tbl1_id -> schm.tbl3.tbl1_id"; actual != expected {
			t.Errorf("Foreign key's column reference is invalid. expected: %v, actual: %v", expected, actual)
		}
		if actual, expected := colRefToString(tbl.ForeignKeys()[2].ColumnReferences()[1]), "other.tbl_other.tbl2_id -> schm.tbl3.tbl2_id"; actual != expected {
			t.Errorf("Foreign key's column reference is invalid. expected: %v, actual: %v", expected, actual)
		}
	}
}

func TestPostgresTableReferencedKeysCount(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl1")
	for _, tbl := range tbls {
		if actual, expected := len(tbl.ReferencedKeys()), 3; actual != expected {
			t.Errorf("Referenced key count is invalid. expected: %v, actual: %v", expected, actual)
		}
	}
}

func TestPostgresTableReferencedKeysOrder(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl1")
	for _, tbl := range tbls {
		names := []string{"tbl2_fk1", "tbl3_fk1", "tbl_other_fk1"}
		for i, expected := range names {
			if actual := tbl.ReferencedKeys()[i].Name(); actual != expected {
				t.Errorf("Referenced key order is invalid. expected: %v, actual: %v", expected, actual)
			}
		}
	}
}

func TestPostgresTableReferencedKeysColumnCount(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl1")
	for _, tbl := range tbls {
		if actual, expected := len(tbl.ReferencedKeys()[0].ColumnReferences()), 1; actual != expected {
			t.Errorf("Referenced key's column count is invalid. expected: %v, actual: %v", expected, actual)
		}
		if actual, expected := len(tbl.ReferencedKeys()[1].ColumnReferences()), 1; actual != expected {
			t.Errorf("Referenced key's column count is invalid. expected: %v, actual: %v", expected, actual)
		}
		if actual, expected := len(tbl.ReferencedKeys()[2].ColumnReferences()), 1; actual != expected {
			t.Errorf("Referenced key's column count is invalid. expected: %v, actual: %v", expected, actual)
		}
	}

	tbls = loadPostgresTableBy2Way("schm", "tbl3")
	for _, tbl := range tbls {
		if actual, expected := len(tbl.ReferencedKeys()[0].ColumnReferences()), 2; actual != expected {
			t.Errorf("Referenced key's column count is invalid. expected: %v, actual: %v", expected, actual)
		}
	}
}

func TestPostgresTableReferencedKeysColumnReferences(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl1")
	for _, tbl := range tbls {
		if actual, expected := colRefToString(tbl.ReferencedKeys()[0].ColumnReferences()[0]), "schm.tbl2.tbl1_id -> schm.tbl1.id"; actual != expected {
			t.Errorf("Referenced key's column reference is invalid. expected: %v, actual: %v", expected, actual)
		}
		if actual, expected := colRefToString(tbl.ReferencedKeys()[1].ColumnReferences()[0]), "schm.tbl3.tbl1_id -> schm.tbl1.id"; actual != expected {
			t.Errorf("Referenced key's column reference is invalid. expected: %v, actual: %v", expected, actual)
		}
		if actual, expected := colRefToString(tbl.ReferencedKeys()[2].ColumnReferences()[0]), "other.tbl_other.tbl1_id -> schm.tbl1.id"; actual != expected {
			t.Errorf("Referenced key's column reference is invalid. expected: %v, actual: %v", expected, actual)
		}
	}

	tbls = loadPostgresTableBy2Way("schm", "tbl3")
	for _, tbl := range tbls {
		if actual, expected := colRefToString(tbl.ReferencedKeys()[0].ColumnReferences()[0]), "other.tbl_other.tbl1_id -> schm.tbl3.tbl1_id"; actual != expected {
			t.Errorf("Referenced key's column reference is invalid. expected: %v, actual: %v", expected, actual)
		}
		if actual, expected := colRefToString(tbl.ReferencedKeys()[0].ColumnReferences()[1]), "other.tbl_other.tbl2_id -> schm.tbl3.tbl2_id"; actual != expected {
			t.Errorf("Referenced key's column reference is invalid. expected: %v, actual: %v", expected, actual)
		}
	}
}

func TestPostgresTableConstraitsCount(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl2")
	for _, tbl := range tbls {
		if len(tbl.Constraints()) != 2 {
			t.Errorf("Constraint count is invalid. expected: %v, actual: %v", 1, len(tbl.Constraints()))
		}
	}
}

func TestPostgresTableConstraitsOrder(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl2")
	for _, tbl := range tbls {
		if actual, expected := tbl.Constraints()[0].Name(), "tbl2_chk1"; actual != expected {
			t.Errorf("Constraint order is invalid. expected: %v, actual: %v", expected, actual)
		}
		if actual, expected := tbl.Constraints()[1].Name(), "tbl2_ex1"; actual != expected {
			t.Errorf("Constraint order is invalid. expected: %v, actual: %v", expected, actual)
		}
	}
}

func TestPostgresTableConstraintKind(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl2")
	for _, tbl := range tbls {
		if actual, expected := tbl.Constraints()[0].Kind(), "CHECK"; actual != expected {
			t.Errorf("Constraint kind is invalid. expected: %v, actual: %v", expected, actual)
		}
		if actual, expected := tbl.Constraints()[1].Kind(), "EXCLUDE"; actual != expected {
			t.Errorf("Constraint kind is invalid. expected: %v, actual: %v", expected, actual)
		}
	}
}

func TestPostgresTableConstraintContent(t *testing.T) {
	tbls := loadPostgresTableBy2Way("schm", "tbl2")
	for _, tbl := range tbls {
		if actual, expected := tbl.Constraints()[0].Content(), "(chk > 0)"; actual != expected {
			t.Errorf("Constraint kind is invalid. expected: %v, actual: %v", expected, actual)
		}
		if actual, expected := tbl.Constraints()[1].Content(), "range WITH &&"; actual != expected {
			t.Errorf("Constraint kind is invalid. expected: %v, actual: %v", expected, actual)
		}
	}
}

func TestPostgresAllTableCount(t *testing.T) {
	tbls := loadPostgresAllTables("schm")
	if actual, expected := len(tbls), 3; actual != expected {
		t.Errorf("AllTables returns should be all tables in 'schm' schema. expected: %v, actual: %v", expected, actual)
	}
}

func TestPostgresAllTablesOrder(t *testing.T) {
	tbls := loadPostgresAllTables("schm")
	tblNames := []string{"tbl1", "tbl2", "tbl3"}
	for i, expected := range tblNames {
		if actual := tbls[i].Name(); actual != expected {
			t.Errorf("Table order is invalid. expected: %v, actual: %v", expected, actual)
		}
	}
}

func TestPostgresAllTablesTableComment(t *testing.T) {
	tbls := loadPostgresAllTables("schm")
	comments := []string{"", "This is table2", ""}
	for i, expected := range comments {
		if actual := tbls[i].Comment(); actual != expected {
			t.Errorf("Table comment is invalid. expected: %v, actual: %v", expected, actual)
		}
	}
}

func createPostgresClient() *Client {
	return NewClient(createPostgresDataSource("postgres", "9.4"))
}

func createPostgresDataSource(driver string, version string) DataSource {
	return NewDataSource(driver, version, "localhost", 5432, "postgres", "", "dbmodel_test", map[string]string{"sslmode": "disable"})
}

func loadPostgresTableBy2Way(schema string, name string) []*Table {
	tbl1 := loadPostgresTable(schema, name)
	tbl2 := loadPostgresTableByAllTables(schema, name)
	return []*Table{tbl1, tbl2}
}

func loadPostgresTable(schema string, name string) *Table {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	t, err := c.Table(schema, name, RequireAll)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func loadPostgresAllTables(schema string) []*Table {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	tbls, err := c.AllTables(schema, RequireAll)
	if err != nil {
		log.Fatal(err)
	}
	return tbls
}

func loadPostgresTableByAllTables(schema string, name string) *Table {
	tbls := loadPostgresAllTables(schema)
	for _, t := range tbls {
		if t.Name() == name {
			return t
		}
	}
	return nil
}

func colRefToString(colRef *ColumnReference) string {
	return fmt.Sprintf("%v.%v.%v -> %v.%v.%v",
		colRef.From().Schema(), colRef.From().TableName(), colRef.From().Name(),
		colRef.To().Schema(), colRef.To().TableName(), colRef.To().Name())
}
