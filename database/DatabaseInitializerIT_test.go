package database

import (
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	sqliteInMemoryOption = "file::memory:?cache=shared"
)

func TestInitDatabase(t *testing.T) {
	dbConn, err := ConnectToServer(sqliteInMemoryOption)
	assert.NotNil(t, dbConn)
	assert.NoError(t, err)

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("TestInitDatabase: Problem occurred while trying to close database -> %v", err)
		}
	}(dbConn)

	err = ExecuteInitSQL(dbConn, "./init.sql")
	assert.NoError(t, err)

	tables := [5]string{"system_config_default", "addon", "run_log", "download_log", "extract_log"}

	for _, table := range tables {
		sqliteTableCheck := "SELECT name FROM sqlite_master WHERE type='table' AND name='" + table + "'"
		assert.NotNil(t, sqliteTableCheck)
	}
}
