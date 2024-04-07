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
	"errors"
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

func PollingExtractor(runId string, dp string, ap string, stopChannel <-chan struct{}, errorChannel chan<- error) {
	log.Printf("PollingExtractor started -> RunID: [%s]\n", runId)

	for {
		select {
		case <-stopChannel:
			log.Println("PollingExtractor routine stopped!")
			return
		default:
			time.Sleep(1 * time.Second)

			fl, found, err := getAddonsAtPath(dp)
			if err != nil {
				errorChannel <- errors.New(fmt.Sprintf("%v", err))
			}

			if found {
				err = extractFiles(runId, dp, ap, fl)
				if err != nil {
					errorChannel <- errors.New(fmt.Sprintf("%v", err))
				}
			}
		}
	}
}

// filterFilesForAddons: TODO
func filterFilesForAddons(fl []os.DirEntry) (ff []os.DirEntry, err error) {
	// Get all addons file names from database
	// Filter so that we only accept those that are part of the list of names from database
	return
}

func extractFiles(runId string, dp string, ap string, fl []os.DirEntry) (err error) {
	ff, err := filterFilesForAddons(fl)
	if err != nil {
		return err
	}

	log.Printf("Will now commence extraction of a total of %d addons...\n", len(fl))
	for i, f := range ff {
		// Combine the directory path with the entry name to get the absolute path
		absPath := filepath.Join(dp, f.Name())
		log.Printf("[%d/%d] System extract of %s!\n", i+1, len(fl), f.Name())
		err := system.Extract(absPath, ap)
		if err != nil {
			return err
		}

		id, err := database.InsertELog(db, models.ELog{
			RunId: strings.Replace(runId, "-", "", -1),
			File:  absPath,
		})
		if err != nil {
			log.Fatalf("Failed to insert extract log into database for the f: %s -> %v\n", f.Name(), err)
		}

		log.Printf("Stored extract logging with ID: %d, for the f: %s\n", id, f.Name())

		err = os.Remove(absPath)
		if err != nil {
			return err
		}
	}
	return
}

func getAddonsAtPath(p string) (f []os.DirEntry, found bool, err error) {
	f, err = os.ReadDir(p)
	if err != nil {
		return nil, false, err
	}

	if len(f) == 0 {
		return nil, false, nil
	} else {
		found = true
	}
	return
}
