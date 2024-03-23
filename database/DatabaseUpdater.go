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
	"time"
)

func UpdateAddon(db *sql.DB, addon models.Addon) (r int64, err error) {
	result, err := db.Exec("UPDATE kaasufouji.addons SET last_downloaded = ? WHERE name = ?",
		time.Now().UTC(), addon.Name)
	if err != nil {
		return 0, fmt.Errorf("database.UpdateAddon: %v", err)
	}

	r, err = result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("database.UpdateAddon: %v", err)
	}

	return
}

func InsertAddon(db *sql.DB, addon models.Addon) (r int64, err error) {
	result, err := db.Exec(("INSERT IGNORE INTO kaasufouji.addons (name, filename, url, download_url) VALUES (?, ?, ?, ?);"),
		addon.Name, addon.Filename, addon.Url, addon.DownloadUrl)
	if err != nil {
		return 0, fmt.Errorf("database.InsertAddon :: %v\n", err)
	}

	r, err = result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("database.InsertAddon :: %v\n", err)
	}
	return
}

func RemoveAddonByID(db *sql.DB, id int) (r int64, err error) {
	result, err := db.Exec("DELETE FROM kaasufouji.addons WHERE id = ?;", id)
	if err != nil {
		return 0, fmt.Errorf("database.RemoveAddonByID :: Error while trying to delete addon by ID [%d] -> %v\n",
			id, err)
	}

	r, err = result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("database.RemoveAddonByID :: Error while trying to delete addon by ID [%d] -> %v\n",
			id, err)
	}
	return
}
