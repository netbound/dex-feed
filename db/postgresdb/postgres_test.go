package postgresdb

import (
	"fmt"
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
	if val, ok := pdb.Get("hello"); ok {
		fmt.Println(val)
	}
}
