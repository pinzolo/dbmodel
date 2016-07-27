# dbmodel

[![Build Status](https://travis-ci.org/pinzolo/dbmodel.png)](http://travis-ci.org/pinzolo/dbmodel)
[![Coverage Status](https://coveralls.io/repos/github/pinzolo/dbmodel/badge.svg?branch=master)](https://coveralls.io/github/pinzolo/dbmodel?branch=master)

## Description

dbmodel is simple database modeling.  
dbmodel does not import database driver in production code. You must import database driver in your code.

## Databases

* PostgreSQL: higher 8.3
* MySQL: not supported yet
* Oracle: not supported yet
* SQL Server: not supported yet

## Usage

Install and use in your code.

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
