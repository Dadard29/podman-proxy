package db_test

import (
	"os"

	"github.com/Dadard29/podman-proxy/db"
)

var testDbPath = "test-db.sqlite3"

func NewTestDb() (*db.Db, error) {
	testDb, err := db.NewDb(testDbPath)
	if err != nil {
		return nil, err
	}
	testDb.Init()

	return testDb, nil
}

func CleanTestDb() {
	os.Remove(testDbPath)
}
