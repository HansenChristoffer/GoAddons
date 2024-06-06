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
	"goaddons/models"
)

func InsertRLog(db *sql.DB, rlog models.RLog) (int64, error) {
	result, err := db.Exec("INSERT OR IGNORE INTO run_log (run_id, service) VALUES (?, ?)",
		rlog.RunId, rlog.Service)
	if err != nil {
		return 0, fmt.Errorf("InsertRLog: %v\n", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("InsertRLog: %v\n", err)
	}

	return id, nil
}

func InsertDLog(db *sql.DB, dlog models.DLog) (int64, error) {
	result, err := db.Exec("INSERT OR IGNORE INTO download_log (run_id, url) VALUES (?, ?)",
		dlog.RunId, dlog.Url)
	if err != nil {
		return 0, fmt.Errorf("InsertDLog: %v\n", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("InsertDLog: %v\n", err)
	}

	return id, nil
}

func InsertELog(db *sql.DB, elog models.ELog) (int64, error) {
	result, err := db.Exec("INSERT OR IGNORE INTO extract_log (run_id, file) VALUES (?, ?)",
		elog.RunId, elog.File)
	if err != nil {
		return 0, fmt.Errorf("InsertELog: %v\n", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("InsertELog: %v\n", err)
	}

	return id, nil
}
