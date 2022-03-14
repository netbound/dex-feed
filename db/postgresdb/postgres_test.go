package postgresdb

import (
	"testing"
)

func testNew(name string) (*PostgresDb, error) {
	connStr := "postgres://dex-feed:dex-feed@localhost/dex-feed?sslmode=disable"

	pdb, err := NewDbCache(connStr, name)
	if err != nil {
		return nil, err
	}

	return pdb, err
}

func TestNew(t *testing.T) {
	pdb, err := testNew("testNew")
	if err != nil {
		t.Fatal(err)
	}

	if err := pdb.db.Ping(); err != nil {
		t.Fatal(err)
	}
}

func TestPut(t *testing.T) {
	pdb, err := testNew("testNew")
	if err != nil {
		t.Fatal(err)
	}

	pdb.Put("hello", "world")
}

func TestGet(t *testing.T) {
	pdb, err := testNew("testNew")
	if err != nil {
		t.Fatal(err)
	}
	if val, ok := pdb.Get("hello"); ok {
		if val != "world" {
			t.Fatalf("Wrong value %s, should be world", val)
		}
	} else {
		t.Fatal("Value not found")
	}
}
