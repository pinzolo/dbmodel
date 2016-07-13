package pg_test

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/pinzolo/dbmodel"

	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	err := createTestResources()
	if err != nil {
		panic(err)
	}
	code := m.Run()
	defer os.Exit(code)
	err = dropTestResources()
	if err != nil {
		fmt.Println(err)
		code = 2
	}
}

func TestAllTableNames(t *testing.T) {
	c := createClient()
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
	c := createClient()
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
	c := createClient()
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
	c := createClient()
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

func createTestResources() error {
	err := createPgTestDB()
	if err != nil {
		return err
	}
	db, err := sql.Open("postgres", "host=localhost user=postgres dbname=dbmodel_test sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	bytes, err := createResourceSQL()
	_, err = db.Exec(string(bytes))
	if err != nil {
		return err
	}

	return nil
}

func dropTestResources() error {
	db, err := sql.Open("postgres", "host=localhost user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DROP DATABASE IF EXISTS dbmodel_test")
	if err != nil {
		return err
	}
	return nil
}

func createResourceSQL() ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(filepath.Join(wd, "create_resources.sql"))
	defer f.Close()
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func createClient() *dbmodel.Client {
	return dbmodel.NewClient("postgres", createDataSource())
}

func createDataSource() dbmodel.DataSource {
	return dbmodel.NewDataSource("localhost", 5432, "postgres", "", "dbmodel_test", map[string]string{"sslmode": "disable"})
}
