package dbmodel

import (
	"strings"
	"testing"
)

func TestPostgresDataSourceName(t *testing.T) {
	p := postgres{}
	ds := InitDataSource()
	excepted := ""
	if p.dataSourceName(ds) != excepted {
		t.Errorf("Postgres#dataSourceName returns invalid value.(expected: %v, actually: %v)", excepted, p.dataSourceName(ds))
	}
	ds.User = "postgres"
	excepted = "user=postgres"
	if p.dataSourceName(ds) != excepted {
		t.Errorf("Postgres#dataSourceName returns invalid value.(expected: %v, actually: %v)", excepted, p.dataSourceName(ds))
	}
	ds.Password = "12345"
	excepted = "user=postgres password=12345"
	if p.dataSourceName(ds) != excepted {
		t.Errorf("Postgres#dataSourceName returns invalid value.(expected: %v, actually: %v)", excepted, p.dataSourceName(ds))
	}
	ds.Host = "localhost"
	excepted = "host=localhost user=postgres password=12345"
	if p.dataSourceName(ds) != excepted {
		t.Errorf("Postgres#dataSourceName returns invalid value.(expected: %v, actually: %v)", excepted, p.dataSourceName(ds))
	}
	ds.Port = 5432
	excepted = "host=localhost port=5432 user=postgres password=12345"
	if p.dataSourceName(ds) != excepted {
		t.Errorf("Postgres#dataSourceName returns invalid value.(expected: %v, actually: %v)", excepted, p.dataSourceName(ds))
	}
	ds.Database = "sample"
	excepted = "host=localhost port=5432 user=postgres password=12345 dbname=sample"
	if p.dataSourceName(ds) != excepted {
		t.Errorf("Postgres#dataSourceName returns invalid value.(expected: %v, actually: %v)", excepted, p.dataSourceName(ds))
	}
	ds.Options["sslmode"] = "disable"
	excepted = "host=localhost port=5432 user=postgres password=12345 dbname=sample sslmode=disable"
	if p.dataSourceName(ds) != excepted {
		t.Errorf("Postgres#dataSourceName returns invalid value.(expected: %v, actually: %v)", excepted, p.dataSourceName(ds))
	}
}

func TestPostgresAllTableNames(t *testing.T) {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	ts, err := c.AllTableNames("sales")
	if err != nil {
		t.Error(err)
	}
	if len(ts) != 2 {
		t.Errorf("AllTableNames should return 2 table names. but actual %v", len(ts))
		return
	}
	if ts[0].Name() != "country_region_currency" {
		t.Errorf("AllTableNames returns invalid table name. expected 'country_region_currency', but actual '%v'", ts[0].Name())
	}
	if ts[0].Comment() != "" {
		t.Errorf("Table comment is null, Comment() should return empty")
	}
	if ts[1].Name() != "currency" {
		t.Errorf("AllTableNames returns invalid table name. expected 'currency', but actual '%v'", ts[1].Name())
	}
	if ts[1].Comment() != "Lookup table containing standard ISO currencies." {
		t.Errorf("AllTableNames should pick up table comment. %+v", ts[1])
	}
}

func TestPostgresAllTableNamesOtherSchema(t *testing.T) {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	ts, err := c.AllTableNames("person")
	if err != nil {
		t.Error(err)
	}
	if len(ts) != 1 {
		t.Errorf("AllTableNames should return 1 table name. but actual %v", len(ts))
		return
	}
	if ts[0].Name() != "country_region" {
		t.Errorf("AllTableNames returns invalid table name. expected 'country_region', but actual '%v'", ts[0].Name())
	}
}

func TestPostgresTableNames(t *testing.T) {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	ts, err := c.TableNames("sales", "region")
	if err != nil {
		t.Error(err)
	}
	if len(ts) != 1 {
		t.Errorf("TableNames with region should return 1 table names. but actual %v", len(ts))
		return
	}
	if ts[0].Name() != "country_region_currency" {
		t.Errorf("TableNames returns invalid table name. expected 'country_region_currency', but actual '%v'", ts[0].Name())
	}
	if ts[0].Comment() != "" {
		t.Errorf("Table comment is null, Comment() should return empty")
	}
}

func TestPostgresTableNamesNoResult(t *testing.T) {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	ts, err := c.TableNames("sales", "sample")
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

	_, err := c.TableNames("", "region")
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

	tbl, err := c.Table("production", "location")
	if err != nil {
		t.Error("Client should not raise error when valid schema and table name given.")
	}
	if tbl.Name() != "location" {
		t.Errorf("Table name is invalid. expected: %v, actual: %v", "location", tbl.Name())
	}
}

func TestPostgresTableNotFound(t *testing.T) {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	_, err := c.Table("production", "xxxxx")
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

	_, err := c.Table("", "location")
	if err == nil {
		t.Errorf("Client should raise error when empty schema given.")
	}
	if err != ErrSchemaEmpty {
		t.Errorf("%v is invalid Error", err)
	}
}

func TestPostgresTableColumnsCount(t *testing.T) {
	tbl := loadPostgresTable("production", "location")
	if len(tbl.Columns()) != 5 {
		t.Errorf("Column count is invalid. expected: %v, actual: %v", 5, len(tbl.Columns()))
	}
}

func TestPostgresTableColumnsOrder(t *testing.T) {
	tbl := loadPostgresTable("production", "location")
	if actual, expected := tbl.Columns()[0].Name(), "location_id"; actual != expected {
		t.Errorf("Column order is invalid. expected: %v, actual: %v", expected, actual)
	}
	if actual, expected := tbl.Columns()[1].Name(), "name"; actual != expected {
		t.Errorf("Column order is invalid. expected: %v, actual: %v", expected, actual)
	}
	if actual, expected := tbl.Columns()[2].Name(), "cost_rate"; actual != expected {
		t.Errorf("Column order is invalid. expected: %v, actual: %v", expected, actual)
	}
	if actual, expected := tbl.Columns()[3].Name(), "availability"; actual != expected {
		t.Errorf("Column order is invalid. expected: %v, actual: %v", expected, actual)
	}
	if actual, expected := tbl.Columns()[4].Name(), "modified_date"; actual != expected {
		t.Errorf("Column order is invalid. expected: %v, actual: %v", expected, actual)
	}
}

func TestPostgresTableColumnComment(t *testing.T) {
	tbl := loadPostgresTable("production", "location")
	if actual, expected := tbl.Columns()[0].Comment(), "Primary key for Location records."; actual != expected {
		t.Errorf("Cannot get valid comment. expected: %v, actual: %v", expected, actual)
	}
	if actual, expected := tbl.Columns()[1].Comment(), ""; actual != expected {
		t.Errorf("Commnet() should return empty when column comment is NULL. actual: %v", actual)
	}
}

func TestPostgresTableColumnDataType(t *testing.T) {
	tbl := loadPostgresTable("production", "location")
	if actual, expected := tbl.Columns()[0].DataType(), "int4"; actual != expected {
		t.Errorf("Cannot get valid data type. expected: %v, actual: %v", expected, actual)
	}
	if actual, expected := tbl.Columns()[1].DataType(), "public.Name"; actual != expected {
		t.Errorf("Cannot get valid custom data type. expected: %v, actual: %v", expected, actual)
	}
}

func TestPostgresTableColumnSize(t *testing.T) {
	tbl := loadPostgresTable("production", "location")
	textSize := tbl.Columns()[1].Size()
	if !textSize.IsValid() || !textSize.Length().Valid || textSize.Precision().Valid {
		t.Error("Cannot get valid text size.")
	}
	if actual, expected := textSize.String(), "50"; actual != expected {
		t.Errorf("Text size value is invalid. expected: %v, actual: %v", expected, actual)
	}

	nullSize := tbl.Columns()[2].Size()
	if nullSize.IsValid() {
		t.Error("Cannot get valid null size.")
	}
	if actual, expected := nullSize.String(), ""; actual != expected {
		t.Errorf("Null size value is invalid. expected: %v, actual: %v", expected, actual)
	}

	intSize := tbl.Columns()[3].Size()
	if !intSize.IsValid() || intSize.Length().Valid || !intSize.Precision().Valid || !intSize.Scale().Valid {
		t.Error("Cannot get valid integer size.")
	}
	if actual, expected := intSize.String(), "8, 2"; actual != expected {
		t.Errorf("Integer size value is invalid. expected: %v, actual: %v", expected, actual)
	}

	dateSize := tbl.Columns()[4].Size()
	if !dateSize.IsValid() || dateSize.Length().Valid || !dateSize.Precision().Valid {
		t.Error("Cannot get valid date size.")
	}
	if actual, expected := dateSize.String(), "6"; actual != expected {
		t.Errorf("Date size value is invalid. expected: %v, actual: %v", expected, actual)
	}
}

func TestPostgresTableColumnNullable(t *testing.T) {
	tbl := loadPostgresTable("sales", "currency")
	if tbl.Columns()[0].IsNullable() {
		t.Errorf("Column '%v' is not nullable, but IsNullable() returns true", tbl.Columns()[0].Name())
	}
	if !tbl.Columns()[2].IsNullable() {
		t.Errorf("Column '%v' is nullable, but IsNullable() returns false", tbl.Columns()[2].Name())
	}
}

func TestPostgresTableColumnDefaultValue(t *testing.T) {
	tbl := loadPostgresTable("production", "location")
	if actual := tbl.Columns()[1].DefaultValue(); actual != "" {
		t.Errorf("Column '%v' do not have default value, but DefaultValue() returns %v", tbl.Columns()[1].Name(), actual)
	}
	if actual, expected := tbl.Columns()[2].DefaultValue(), "0.00"; actual != expected {
		t.Errorf("Cannot get invalid default value of '%v'. expected: %v, actual: %v", tbl.Columns()[2].Name(), expected, actual)
	}
	if actual, expected := tbl.Columns()[4].DefaultValue(), "now()"; actual != expected {
		t.Errorf("Cannot get invalid default value of '%v'. expected: %v, actual: %v", tbl.Columns()[4].Name(), expected, actual)
	}
}

func TestPostgresTableColumnPrimaryKeyPosition(t *testing.T) {
	// Single column primary key
	tbl := loadPostgresTable("production", "location")
	if actual, expected := tbl.Columns()[0].PrimaryKeyPosition(), int64(1); actual != expected {
		t.Errorf("Cannot get invalid primary key position of '%v'. expected: %v, actual: %v", tbl.Columns()[0].Name(), expected, actual)
	}
	if actual, expected := tbl.Columns()[1].PrimaryKeyPosition(), int64(0); actual != expected {
		t.Errorf("Cannot get invalid primary key position of '%v'. expected: %v, actual: %v", tbl.Columns()[0].Name(), expected, actual)
	}

	// Multi columns primary key
	tbl = loadPostgresTable("sales", "country_region_currency")
	if actual, expected := tbl.Columns()[0].PrimaryKeyPosition(), int64(1); actual != expected {
		t.Errorf("Cannot get invalid primary key position of '%v'. expected: %v, actual: %v", tbl.Columns()[0].Name(), expected, actual)
	}
	if actual, expected := tbl.Columns()[1].PrimaryKeyPosition(), int64(2); actual != expected {
		t.Errorf("Cannot get invalid primary key position of '%v'. expected: %v, actual: %v", tbl.Columns()[1].Name(), expected, actual)
	}
}

func createPostgresClient() *Client {
	return NewClient("postgres", createPostgresDataSource())
}

func createPostgresDataSource() DataSource {
	return NewDataSource("localhost", 5432, "postgres", "", "dbmodel_test", map[string]string{"sslmode": "disable"})
}

func loadPostgresTable(schema string, name string) *Table {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	t, _ := c.Table(schema, name)
	return t
}
