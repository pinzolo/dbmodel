package dbmodel

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	err := createPostgresTestResources()
	if err != nil {
		panic(err)
	}
	code := m.Run()
	defer os.Exit(code)
	err = dropPostgresTestResources()
	if err != nil {
		fmt.Println(err)
		code = 2
	}
}

func createPostgresTestResources() error {
	err := createPostgresTestDB()
	if err != nil {
		return err
	}
	db, err := sql.Open("postgres", "host=localhost user=postgres dbname=dbmodel_test sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	bytes, err := readSQLFile("create_postgres_resources")
	_, err = db.Exec(string(bytes))
	if err != nil {
		return err
	}

	return nil
}

func createPostgresTestDB() error {
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

func dropPostgresTestResources() error {
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

func readSQLFile(fileName string) ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(filepath.Join(wd, "testdata", fileName+".sql"))
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
