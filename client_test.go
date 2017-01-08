package dbmodel

import (
	"testing"
)

func BenchmarkLoadAllTables(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := NewClient(createPostgresDataSource("postgres", "9.4"))
		c.AllTables("schm", RequireAll)
	}
}

func TestUnconnectedClientRaisesError(t *testing.T) {
	c := NewClient(createPostgresDataSource("postgres", "9.4"))
	_, err := c.AllTableNames("schm")
	if err == nil {
		t.Errorf("Client should raise error when use on unconnected.")
	}
	if err != ErrConnNotFound {
		t.Errorf("%v is invalid Error", err)
	}
}

func TestInvalidDriver(t *testing.T) {
	c := NewClient(createPostgresDataSource("foobar", "9.4"))
	c.Connect()
	defer c.Disconnect()
	_, err := c.AllTableNames("schm")
	if err == nil {
		t.Errorf("Client should raise error when unknown driver given.")
	}
	if err != ErrInvalidDriver {
		t.Errorf("%v is invalid Error", err)
	}
}

func TestUnconnectedClientDisconnectSafe(t *testing.T) {
	c := NewClient(createPostgresDataSource("postgres", "9.4"))
	err := c.Disconnect()
	if err != nil {
		t.Errorf("Client should not raise error when disconnect on unconnected.")
	}
}

func TestUseCustomProviderWhenInvalidDriverGiven(t *testing.T) {
	c := NewClient(createPostgresDataSource("foobar", "9.4"))
	c.SetProvider(postgres{ds: createPostgresDataSource("postgres", "9.4")})
	c.Connect()
	defer c.Disconnect()
	_, err := c.AllTableNames("schm")
	if err != nil {
		t.Errorf("Client should not raise error when valid provider and unknown driver given.")
	}
}
