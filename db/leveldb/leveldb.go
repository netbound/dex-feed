package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type Database struct {
	ldb *leveldb.DB
}

func NewDatabase(filePath string) (*Database, error) {
	ldb, err := leveldb.OpenFile(filePath, nil)
	if err != nil {
		return nil, err
	}

	return &Database{ldb}, nil

	// TODO: put this somewhere on exit
	// defer db.Close()
}

func (db Database) Get(key string) ([]byte, bool) {
	val, err := db.ldb.Get([]byte(key), nil)
	if err != nil {
		return []byte{}, false
	}

	return val, true
}

func (db *Database) Put(key string, value []byte) {
	db.ldb.Put([]byte(key), value, nil)
}
