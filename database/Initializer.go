package database

import (
	"database/sql"
	"os"
	"strings"
)

func ExecuteInitSQL(dbConn *sql.DB) error {
	sqlFileContent, err := os.ReadFile("init.sql")
	if err != nil {
		return err
	}

	// Convert the file content to a string and split into individual statements.
	sqlCommands := strings.Split(string(sqlFileContent), ";")

	// Begin a transaction.
	trans, err := dbConn.Begin()
	if err != nil {
		return err
	}

	// Execute each statement.
	for _, command := range sqlCommands {
		command = strings.TrimSpace(command)
		if command == "" {
			continue // skip empty commands
		}

		_, err := trans.Exec(command)
		if err != nil {
			// If an error occurs, rollback the transaction and return the error.
			_ = trans.Rollback() // Ignore rollback error, focus on the original error.
			return err
		}
	}

	// Commit the transaction.
	if err := trans.Commit(); err != nil {
		return err
	}

	return nil
}
