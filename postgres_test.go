package dbmodel

import (
	"database/sql"
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

func TestAllTableNames(t *testing.T) {
	c := createPgClient()
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

func TestAllTableNamesOtherSchema(t *testing.T) {
	c := createPgClient()
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

func TestTableNames(t *testing.T) {
	c := createPgClient()
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

func TestTableNamesNoResult(t *testing.T) {
	c := createPgClient()
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

func createPgTestDB() error {
	db, err := sql.Open("postgres", "host=localhost user=postgres sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE dbmodel_test")
	if err != nil {
		return err
	}

	return nil
}

func createPgClient() *Client {
	ds := NewDataSource("localhost", 5432, "postgres", "", "dbmodel_test", map[string]string{"sslmode": "disable"})
	return NewClient("postgres", ds)
}
