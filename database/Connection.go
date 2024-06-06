// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

const (
	databasePath             = "./init.sql"
	databaseType             = "sqlite"
	maxConnLifetimeInMinutes = 30
	maxOpenConns             = 1
	maxIdleConns             = 1
)

// ConnectToServer opens a connection to the SQLite database and sets up the connection parameters.
func ConnectToServer(database string) (dbConn *sql.DB, err error) {
	// Open a connection to the SQLite database.
	dbConn, err = sql.Open(databaseType, database)
	if err != nil {
		return nil, fmt.Errorf("database:Connection.ConnectToServer():sql.Open(%s, %s) "+
			"-> %w", databaseType, database, err)
	}

	// Set the connection parameters.
	dbConn.SetConnMaxLifetime(time.Minute * maxConnLifetimeInMinutes)
	dbConn.SetMaxOpenConns(maxOpenConns)
	dbConn.SetMaxIdleConns(maxIdleConns)

	// Ping the database to ensure the connection is valid.
	err = dbConn.Ping()
	if err != nil {
		return nil, fmt.Errorf("database:Connection.ConnectToServer():dbConn.Ping() "+
			"-> %w", err)
	}
	log.Println("Connected!")

	// Log connection statistics.
	stats := dbConn.Stats()
	log.Printf("\n### Statistics ###\nConnections: %d/%d\nNumber of active connections: %d\n"+
		"Number of idle connections: %d\n",
		stats.OpenConnections, stats.MaxOpenConnections, stats.InUse, stats.Idle)

	// Execute any initialization SQL.
	err = ExecuteInitSQL(dbConn, databasePath)
	if err != nil {
		return nil, fmt.Errorf("database:Connection.ConnectToServer():ExecuteInitSQL(dbConn) "+
			"-> %w", err)
	}
	return dbConn, nil
}
