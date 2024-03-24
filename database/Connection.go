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
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

func ConnectToServer() *sql.DB {
	var dbConn *sql.DB

	// Capture connection properties.
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "defcon",
		AllowNativePasswords: true,
	}

	// Get a database handle.
	var err error
	dbConn, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	// See "Important settings" section.
	dbConn.SetConnMaxLifetime(time.Minute * 10)
	dbConn.SetMaxOpenConns(5)
	dbConn.SetMaxIdleConns(5)

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

	return dbConn
}
