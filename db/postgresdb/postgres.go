package postgresdb

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresDb struct {
	db *sql.DB
}

func New(connStr string) (*PostgresDb, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresDb{db}, nil
}
