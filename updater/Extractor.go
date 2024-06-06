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

package updater

import (
	"database/sql"
	"fmt"
	"goaddons/database"
	"goaddons/models"
	"goaddons/system"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	sleepDuration = 1000 * time.Millisecond
)

var db *sql.DB
var cachedAddons []models.Addon

// PollingExtractor continuously polls the specified directory for addons and extracts them if found
func PollingExtractor(dbConn *sql.DB, runId string, dp string, ap string, stopChannel <-chan struct{}, errorChannel chan<- error) {
	log.Printf("PollingExtractor started -> RunID: [%s]\n", runId)
	db = dbConn

	for {
		select {
		case <-stopChannel:
			log.Println("PollingExtractor routine stopped!")
			return
		default:
			time.Sleep(1 * time.Second)

			fl, found, err := getAddonsAtPath(dp)
			if err != nil {
				errorChannel <- fmt.Errorf("updater:Extractor.PollingExtractor():getAddonsAtPath(%s) "+
					"-> %w", dp, err)
				continue
			}

			if found {
				err = extractFiles(runId, dp, ap, fl)
				if err != nil {
					errorChannel <- fmt.Errorf("updater:Extractor.PollingExtractor():extractFiles(%s, %s, %s) "+
						"-> %w", runId, dp, ap, err)
				}
			}
			time.Sleep(sleepDuration)
		}
	}
}

// filterFilesForAddons filters the directory entries to include only those present in the database
func filterFilesForAddons(files []os.DirEntry) (filteredAddons []os.DirEntry, err error) {
	if len(cachedAddons) == 0 {
		cachedAddons, err = database.GetAllAddons(db)
		if err != nil {
			return nil, fmt.Errorf("updater:Extractor.filterFilesForAddons():database.GetAllAddons(db) "+
				"-> %w", err)
		}
	}

	for _, dirEntry := range files {
		log.Printf("[DEBUG] DirEntry: %s\n", dirEntry.Name())
		for _, addon := range cachedAddons {
			if strings.Contains(dirEntry.Name(), addon.Filename) {
				log.Printf("[DEBUG] %s contains %s == true\n", dirEntry.Name(), addon.Filename)
				filteredAddons = append(filteredAddons, dirEntry)
			}
		}
	}
	return filteredAddons, nil
}

// extractFiles processes and extracts the filtered addon files
func extractFiles(runId string, dp string, ap string, files []os.DirEntry) error {
	filteredFiles, err := filterFilesForAddons(files)
	if err != nil {
		return fmt.Errorf("updater:Extractor.extractFiles():filterFilesForAddons(files) -> %w", err)
	}

	log.Printf("Will now commence extraction of a total of %d addons...\n", len(filteredFiles))
	for idx, file := range filteredFiles {
		absPath := filepath.Join(dp, file.Name())
		log.Printf("[%d/%d] System extract of %s!\n", idx+1, len(filteredFiles), file.Name())
		err := system.Extract(absPath, ap)
		if err != nil {
			return fmt.Errorf("updater:Extractor.extractFiles():system.Extract(%s, %s) "+
				"-> %w", absPath, ap, err)
		}

		id, err := database.InsertELog(db, models.ELog{
			RunId: strings.Replace(runId, "-", "", -1),
			File:  absPath,
		})
		if err != nil {
			return fmt.Errorf("updater:Extractor.extractFiles():database.InsertELog(db, models.ELog) "+
				"-> %w", err)
		}

		log.Printf("Stored extract logging with ID: %d, for the file: %s\n", id, file.Name())
		err = os.Remove(absPath)
		if err != nil {
			return fmt.Errorf("updater:Extractor.extractFiles():os.Remove(%s) -> %w", absPath, err)
		}
	}
	return nil
}

// getAddonsAtPath retrieves the list of directory entries at the specified path that contains .zip
func getAddonsAtPath(path string) (files []os.DirEntry, found bool, err error) {
	readFiles, err := os.ReadDir(path)
	if err != nil {
		return nil, false, fmt.Errorf("updater:Extractor.getAddonsAtPath():os.ReadDir(%s) "+
			"-> %w", path, err)
	}
	for _, file := range readFiles {
		if strings.HasSuffix(file.Name(), ".zip") {
			files = append(files, file)
		}
	}
	found = len(files) > 0
	return files, found, nil
}
