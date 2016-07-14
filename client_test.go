package dbmodel

import (
	"testing"
)

func TestUnconnectedClientRaisesError(t *testing.T) {
	c := NewClient("postgres", createPgDataSource())
	_, err := c.AllTableNames("sales")
	if err == nil {
		t.Errorf("Client should raise error when use on unconnected.")
	}
	if err != ErrConnNotFound {
		t.Errorf("%v is invalid Error", err)
	}
}

func TestInvalidDriver(t *testing.T) {
	c := NewClient("foobar", createPgDataSource())
	c.Connect()
	defer c.Disconnect()
	_, err := c.AllTableNames("sales")
	if err == nil {
		t.Errorf("Client should raise error when unknown driver given.")
	}
	if err != ErrInvalidDriver {
		t.Errorf("%v is invalid Error", err)
	}
}

func TestUnconnectedClientDisconnectSafe(t *testing.T) {
	c := NewClient("postgres", createPgDataSource())
	err := c.Disconnect()
	if err != nil {
		t.Errorf("Client should not raise error when disconnect on unconnected.")
	}
}

func TestUseCustomProviderWhenInvalidDriverGiven(t *testing.T) {
	c := NewClient("foobar", createPgDataSource())
	c.SetProvider(postgres{})
	c.Connect()
	defer c.Disconnect()
	_, err := c.AllTableNames("sales")
	if err != nil {
		t.Errorf("Client should not raise error when valid provider and unknown driver given.")
	}
}
