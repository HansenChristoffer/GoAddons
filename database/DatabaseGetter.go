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

	"goaddons/models"
)

func GetSystemConfigurations(db *sql.DB) (systemConfig []models.SystemConfig, err error) {
	rows, err := db.Query("SELECT * FROM defcon.system_config_default")
	if err != nil {
		return nil, fmt.Errorf("database.GetSystemConfigurations :: %v", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("database.GetSystemConfigurations:rows.Close() :: %v\n", err)
			return
		}
	}(rows)

	for rows.Next() {
		var config models.SystemConfig
		if err := rows.Scan(&config.Name, &config.Path); err != nil {
			return nil, fmt.Errorf("database.GetSystemConfigurations :: %v", err)
		}
		systemConfig = append(systemConfig, config)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("database.GetSystemConfigurations :: %v", err)
	}
	return
}

func GetAllAddons(db *sql.DB) (addons []models.Addon, err error) {
	rows, err := db.Query("SELECT id, name, filename, url, download_url, last_downloaded, last_modified_at, added_at " +
		"FROM kaasufouji.addons WHERE download_url IS NOT NULL AND download_url != '';")
	if err != nil {
		return nil, fmt.Errorf("database.GetAllAddons :: %v", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("GetAllAddons.GetAllAddons:rows.Close() :: %v\n", err)
			return
		}
	}(rows)

	for rows.Next() {
		var addon models.Addon
		if err := rows.Scan(&addon.Id, &addon.Name, &addon.Filename, &addon.Url, &addon.DownloadUrl,
			&addon.LastDownloaded, &addon.LastModifiedAt, &addon.AddedAt); err != nil {
			return nil, fmt.Errorf("database.GetAllAddons :: %v", err)
		}
		addons = append(addons, addon)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("database.GetAllAddons :: %v", err)
	}
	return
}

func GetAddonsByName(db *sql.DB, name string) (addons []models.Addon, err error) {
	rows, err := db.Query("SELECT id, name, filename, url, download_url, last_downloaded, last_modified_at, added_at " +
		"FROM kaasufouji.addons WHERE name LIKE '%" + name + "%'")
	if err != nil {
		return nil, fmt.Errorf("database.GetAddonsByName :: Error while searching for addon! -> %v\n", err)
	}

	defer func(rows *sql.Rows) {
		rowsErr := rows.Close()
		if rowsErr != nil {
			log.Printf("database.GetAddonsByName:rows.Close() :: Error while trying to close rows! -> %v\n", rowsErr)
			return
		}
	}(rows)

	for rows.Next() {
		var addon models.Addon
		if err = rows.Scan(&addon.Id, &addon.Name, &addon.Filename, &addon.Url, &addon.DownloadUrl,
			&addon.LastDownloaded, &addon.LastModifiedAt, &addon.AddedAt); err != nil {
			return nil, fmt.Errorf("database.GetAddonsByName :: %v\n", err)
		}
		addons = append(addons, addon)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("database.GetAddonsByName :: %v", err)
	}
	return
}
