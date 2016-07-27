package dbmodel

import (
	"log"
	"testing"
)

func TestTableWithOptionRequireNone(t *testing.T) {
	tbl := loadPostgresTableWithOpt(RequireNone)
	if len(tbl.Indices()) > 0 {
		t.Error("Indices options is false, but Indices loaded.")
	}
	if len(tbl.ForeignKeys()) > 0 {
		t.Error("ForeignKeys options is false, but ForeignKeys loaded.")
	}
	if len(tbl.ReferencedKeys()) > 0 {
		t.Error("ReferencedKeys options is false, but ReferencedKeys loaded.")
	}
	if len(tbl.Constraints()) > 0 {
		t.Error("Constraints options is false, but Constraints loaded.")
	}
}

func TestTableWithOptionRequireAll(t *testing.T) {
	tbl := loadPostgresTableWithOpt(RequireAll)
	if len(tbl.Indices()) == 0 {
		t.Error("Indices options is true, but Indices not loaded.")
	}
	if len(tbl.ForeignKeys()) == 0 {
		t.Error("ForeignKeys options is true, but ForeignKeys not loaded.")
	}
	if len(tbl.ReferencedKeys()) == 0 {
		t.Error("ReferencedKeys options is true, but ReferencedKeys not loaded.")
	}
	if len(tbl.Constraints()) == 0 {
		t.Error("Constraints options is true, but Constraints not loaded.")
	}
}

func TestTableWithOptionIndices(t *testing.T) {
	tbl := loadPostgresTableWithOpt(Option{Indices: true})
	if len(tbl.Indices()) == 0 {
		t.Error("Indices options is true, but Indices not loaded.")
	}
	if len(tbl.ForeignKeys()) > 0 {
		t.Error("ForeignKeys options is false, but ForeignKeys not empty.")
	}
	if len(tbl.ReferencedKeys()) > 0 {
		t.Error("ReferencedKeys options is false, but ReferencedKeys not empty.")
	}
	if len(tbl.Constraints()) > 0 {
		t.Error("Constraints options is false, but Constraints not empty.")
	}
}

func TestTableWithOptionForeignKeys(t *testing.T) {
	tbl := loadPostgresTableWithOpt(Option{ForeignKeys: true})
	if len(tbl.Indices()) > 0 {
		t.Error("Indices options is false, but Indices loaded.")
	}
	if len(tbl.ForeignKeys()) == 0 {
		t.Error("ForeignKeys options is true, but ForeignKeys not loaded.")
	}
	if len(tbl.ReferencedKeys()) > 0 {
		t.Error("ReferencedKeys options is false, but ReferencedKeys loaded.")
	}
	if len(tbl.Constraints()) > 0 {
		t.Error("Constraints options is false, but Constraints loaded.")
	}
}

func TestTableWithOptionReferencedKeys(t *testing.T) {
	tbl := loadPostgresTableWithOpt(Option{ReferencedKeys: true})
	if len(tbl.Indices()) > 0 {
		t.Error("Indices options is false, but Indices loaded.")
	}
	if len(tbl.ForeignKeys()) > 0 {
		t.Error("ForeignKeys options is false, but ForeignKeys loaded.")
	}
	if len(tbl.ReferencedKeys()) == 0 {
		t.Error("ReferencedKeys options is true, but ReferencedKeys not loaded.")
	}
	if len(tbl.Constraints()) > 0 {
		t.Error("Constraints options is false, but Constraints loaded.")
	}
}

func TestTableWithOptionConstraints(t *testing.T) {
	tbl := loadPostgresTableWithOpt(Option{Constraints: true})
	if len(tbl.Indices()) > 0 {
		t.Error("Indices options is false, but Indices loaded.")
	}
	if len(tbl.ForeignKeys()) > 0 {
		t.Error("ForeignKeys options is false, but ForeignKeys loaded.")
	}
	if len(tbl.ReferencedKeys()) > 0 {
		t.Error("ReferencedKeys options is false, but ReferencedKeys loaded.")
	}
	if len(tbl.Constraints()) == 0 {
		t.Error("Constraints options is true, but Constraints not loaded.")
	}
}

func TestAllTablesWithOptionRequireNone(t *testing.T) {
	tbls := loadPostgresAllTablesWithOpt(RequireNone)
	for _, tbl := range tbls {
		if len(tbl.Indices()) > 0 {
			t.Error("Indices options is false, but Indices loaded.")
		}
		if len(tbl.ForeignKeys()) > 0 {
			t.Error("ForeignKeys options is false, but ForeignKeys loaded.")
		}
		if len(tbl.ReferencedKeys()) > 0 {
			t.Error("ReferencedKeys options is false, but ReferencedKeys loaded.")
		}
		if len(tbl.Constraints()) > 0 {
			t.Error("Constraints options is false, but Constraints loaded.")
		}
	}
}

func TestAllTablesWithOptionRequireAll(t *testing.T) {
	tbls := loadPostgresAllTablesWithOpt(RequireAll)
	var (
		idxLoaded bool
		fkLoaded  bool
		rkLoaded  bool
		conLoaded bool
	)
	for _, tbl := range tbls {
		if len(tbl.Indices()) > 0 {
			idxLoaded = true
		}
		if len(tbl.ForeignKeys()) > 0 {
			fkLoaded = true
		}
		if len(tbl.ReferencedKeys()) > 0 {
			rkLoaded = true
		}
		if len(tbl.Constraints()) > 0 {
			conLoaded = true
		}
	}
	if !idxLoaded {
		t.Error("Indices options is true, but Indices not loaded.")
	}
	if !fkLoaded {
		t.Error("ForeignKeys options is true, but ForeignKeys not loaded.")
	}
	if !rkLoaded {
		t.Error("ReferencedKeys options is true, but ReferencedKeys not loaded.")
	}
	if !conLoaded {
		t.Error("Constraints options is true, but Constraints not loaded.")
	}
}

func TestAllTablesWithOptionIndices(t *testing.T) {
	tbls := loadPostgresAllTablesWithOpt(Option{Indices: true})
	idxLoaded := false
	for _, tbl := range tbls {
		if len(tbl.Indices()) > 0 {
			idxLoaded = true
		}
		if len(tbl.ForeignKeys()) > 0 {
			t.Error("ForeignKeys options is false, but ForeignKeys not empty.")
		}
		if len(tbl.ReferencedKeys()) > 0 {
			t.Error("ReferencedKeys options is false, but ReferencedKeys not empty.")
		}
		if len(tbl.Constraints()) > 0 {
			t.Error("Constraints options is false, but Constraints not empty.")
		}
	}
	if !idxLoaded {
		t.Error("Indices options is true, but Indices not loaded.")
	}
}

func TestAllTablesWithOptionForeignKeys(t *testing.T) {
	tbls := loadPostgresAllTablesWithOpt(Option{ForeignKeys: true})
	fkLoaded := false
	for _, tbl := range tbls {
		if len(tbl.Indices()) > 0 {
			t.Error("Indices options is false, but Indices loaded.")
		}
		if len(tbl.ForeignKeys()) > 0 {
			fkLoaded = true
		}
		if len(tbl.ReferencedKeys()) > 0 {
			t.Error("ReferencedKeys options is false, but ReferencedKeys loaded.")
		}
		if len(tbl.Constraints()) > 0 {
			t.Error("Constraints options is false, but Constraints loaded.")
		}
	}
	if !fkLoaded {
		t.Error("ForeignKeys options is true, but ForeignKeys not loaded.")
	}
}

func TestAllTablesWithOptionReferencedKeys(t *testing.T) {
	tbls := loadPostgresAllTablesWithOpt(Option{ReferencedKeys: true})
	rkLoaded := false
	for _, tbl := range tbls {
		if len(tbl.Indices()) > 0 {
			t.Error("Indices options is false, but Indices loaded.")
		}
		if len(tbl.ForeignKeys()) > 0 {
			t.Error("ForeignKeys options is false, but ForeignKeys loaded.")
		}
		if len(tbl.ReferencedKeys()) > 0 {
			rkLoaded = true
		}
		if len(tbl.Constraints()) > 0 {
			t.Error("Constraints options is false, but Constraints loaded.")
		}
	}
	if !rkLoaded {
		t.Error("ReferencedKeys options is true, but ReferencedKeys not loaded.")
	}
}

func TestAllTablesWithOptionConstraints(t *testing.T) {
	tbls := loadPostgresAllTablesWithOpt(Option{Constraints: true})
	conLoaded := false
	for _, tbl := range tbls {
		if len(tbl.Indices()) > 0 {
			t.Error("Indices options is false, but Indices loaded.")
		}
		if len(tbl.ForeignKeys()) > 0 {
			t.Error("ForeignKeys options is false, but ForeignKeys loaded.")
		}
		if len(tbl.ReferencedKeys()) > 0 {
			t.Error("ReferencedKeys options is false, but ReferencedKeys loaded.")
		}
		if len(tbl.Constraints()) > 0 {
			conLoaded = true
		}
	}
	if !conLoaded {
		t.Error("Constraints options is true, but Constraints not loaded.")
	}
}

func loadPostgresAllTablesWithOpt(opt Option) []*Table {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	tbls, err := c.AllTables("sales", opt)
	if err != nil {
		log.Fatal(err)
	}
	return tbls
}

func loadPostgresTableWithOpt(opt Option) *Table {
	c := createPostgresClient()
	defer c.Disconnect()
	c.Connect()

	tbl, err := c.Table("human_resources", "employee", opt)
	if err != nil {
		log.Fatal(err)
	}
	return tbl
}
