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
	"log"
	"time"

	_ "modernc.org/sqlite"
)

const (
	databaseType             = "sqlite"
	maxConnLifetimeInMinutes = 30
	maxOpenConns             = 1
	maxIdleConns             = 1
)

func ConnectToServer() (dbConn *sql.DB) {
	// Get a database handle.
	dbConn, err := sql.Open(databaseType, "./bin/acd.db")
	if err != nil {
		log.Fatal(err)
	}

	dbConn.SetConnMaxLifetime(time.Minute * maxConnLifetimeInMinutes)
	dbConn.SetMaxOpenConns(maxOpenConns)
	dbConn.SetMaxIdleConns(maxIdleConns)

	pingErr := dbConn.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("Connected!")

	// Stats
	stats := dbConn.Stats()
	log.Printf("\n### Statistics ###\nConnections: %d/%d\nNumber of active connections: %d\n"+
		"Number of idle connections: %d\n",
		stats.OpenConnections, stats.MaxOpenConnections, stats.InUse, stats.Idle)

	err = ExecuteInitSQL(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	return
}
