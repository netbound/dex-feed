package postgresdb

import "testing"

func testNew() (*PostgresDb, error) {
	connStr := "postgres://dex-feed:dex-feed@localhost/dex-feed?sslmode=disable"

	pdb, err := New(connStr)
	if err != nil {
		return nil, err
	}

	return pdb, err
}

func TestNew(t *testing.T) {
	pdb, err := testNew()
	if err != nil {
		t.Fatal(err)
	}

	if err := pdb.db.Ping(); err != nil {
		t.Fatal(err)
	}
}
