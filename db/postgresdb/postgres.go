package postgresdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Implements Cacher
type PostgresDb struct {
	Name string
	db   *sql.DB
}

func NewDbCache(connStr string, cacheName string) (*PostgresDb, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		key BYTEA NOT NULL,
		value BYTEA NOT NULL,
		PRIMARY KEY(key)
	)`, cacheName))
	if err != nil {
		return nil, err
	}

	return &PostgresDb{db: db, Name: cacheName}, nil
}

func (db PostgresDb) Get(key string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var value string

	// TODO: 0x values not allowed
	stmt := fmt.Sprintf(`SELECT value FROM %s WHERE key = $1`, db.Name)

	err := db.db.QueryRowContext(ctx, stmt, key).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false
	}

	if err != nil {
		fmt.Println("error on select", err)
	}

	return value, true
}

func (db *PostgresDb) Put(key string, value string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Put or update value in the db, simulating a cache
	stmt := fmt.Sprintf(`INSERT INTO %s (key, value) values ($1, $2) ON CONFLICT (key) DO UPDATE SET value = $2`, db.Name)

	_, err := db.db.ExecContext(ctx, stmt, key, value)
	if err != nil {
		fmt.Println("error on insert:", err)
	}
}
