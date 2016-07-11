package dbmodel

import "testing"

func TestInitDataSource(t *testing.T) {
	ds := InitDataSource()
	if ds.Options == nil {
		t.Error("Options is nil. Initialized map is expected.")
	}
}
