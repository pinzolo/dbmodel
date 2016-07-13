package dbmodel

import (
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
