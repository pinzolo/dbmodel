# dbmodel

[![Build Status](https://travis-ci.org/pinzolo/dbmodel.png)](http://travis-ci.org/pinzolo/dbmodel)
[![Coverage Status](https://coveralls.io/repos/github/pinzolo/dbmodel/badge.svg?branch=master)](https://coveralls.io/github/pinzolo/dbmodel?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/pinzolo/dbmodel)](https://goreportcard.com/report/github.com/pinzolo/dbmodel)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/pinzolo/dbmodel)
[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/pinzolo/dbmodel/master/LICENSE)

## Description

dbmodel is simple database modeling.  
dbmodel does not import database driver in production code. You must import database driver in your code.

## Databases

* PostgreSQL: higher 8.4
* MySQL: not supported yet
* Oracle: not supported yet
* SQL Server: not supported yet

## Usage

Install and use in your code.

Example (PostgreSQL):  
```go
package main

import (
	"fmt"

	"github.com/pinzolo/dbmodel"
)

func main() {
	// Create DataSouce
	ds := dbmodel.NewDataSource("localhost", 5432, "postgres", "", "sample", map[string]string{"sslmode": "disable"})

	// Create Client
	// First argument is Driver name. this name is used in sql.Open
	client := dbmodel.NewClient("postgres", ds)

	// Connect to Database.
	client.Connect()
	// You must close connection.
	defer client.Disconnect()

	// AllTables returns all table in sample schema.
	// dbmodel.RequireAll is built in option.
	// When dbmodel.RequireAll is given, client loads all metadata of table.(columns, indices, constraints, foreign keys, referenced keys)
	tables, err := client.AllTables("sample", dbmodel.RequireAll)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, tbl := range tables {
		for _, col := range tbl.Columns() {
			fmt.Printf("%#v", col)
		}
		for _, idx := range tbl.Indices() {
			fmt.Printf("%#v", idx)
		}
		for _, cns := range tbl.Constraints() {
			fmt.Printf("%#v", cns)
		}
		for _, fk := range tbl.ForeignKeys() {
			fmt.Printf("%#v", fk)
		}
		for _, rk := range tbl.ReferencedKeys() {
			fmt.Printf("%#v", rk)
		}
	}

	// You can load single table users.
	// You need only columns, you can use dbmodel.RequireNone
	table, err := client.Table("sample", "users", dbmodel.RequireNone)
	fmt.Printf("%#v", table)
	if err != nil {
		fmt.Println(err)
		return
	}

	// You can load table names.
	// Returned *dbmodel.Table contains only table name and comment.
	// As well as using client.TableNames("sample", "users"), you can get tables that contains "users" in its name.
	tables, err = client.AllTableNames("sample")
	if err != nil {
		fmt.Println(err)
		return
	}
}
```

## Install

To install, use `go get`:

```bash
$ go get github.com/pinzolo/dbmodel
```

## Contribution

1. Fork ([https://github.com/pinzolo/dbmodel/fork](https://github.com/pinzolo/dbmodel/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[pinzolo](https://github.com/pinzolo)
