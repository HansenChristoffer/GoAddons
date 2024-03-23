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
	"github.com/google/uuid"
	"goaddons/database"
	"goaddons/models"
	"goaddons/net"
	"goaddons/system"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	serviceName            = "GOADDONS_SERVICE"
	systemsAddonsPathName  = "extract.addon.path"
	browserDownloadDirName = "browser.download.dir"
)

var db *sql.DB

func StartUpdater() {
	log.Printf("Addon Updater starting...")

	// Connect to MySQL database
	db = database.ConnectToServer()

	// Log that the software is making a "run"
	rLog, err := database.InsertRLog(db, models.RLog{RunId: strings.Replace(uuid.New().String(), "-", "", -1),
		Service: serviceName})
	if err != nil {
		log.Printf("Failed to store run logging into database... -> %v\n", err)
	}

	log.Printf("Stored run logging with ID: %d\n", rLog)

	// Fetch systems configuration from database
	config, err := database.GetSystemConfigurations(db)
	if err != nil {
		log.Fatalf("Failed to fetch system configurations... -> %v", err)
	}

	// Fetch the addons path & download dir from the system configurations
	addonsPath, downloadPath, err := getDownloadPathAndAddonsPath(config)
	if err != nil {
		log.Fatalf("Failed to fetch system configurations... -> %v", err)
	}

	startTime := time.Now().UTC()

	// Get all addons from database
	addons, err := database.GetAllAddons(db)
	if err != nil {
		log.Fatalf("Failed to fetch all addons... -> %v", err)
	}

	log.Printf("Fetched a total of %d addons!", len(addons))

	done, err := net.StartHeadlessAndDownloadAddons(addons, downloadPath, db)
	if err != nil {
		log.Fatalf("Error while navigating... -> %v", err)
	}

	if done {
		// Download is now done. Time to extract the files within "downloadPath" to "addonsPath"
		err = extractAllAddons(downloadPath, addonsPath)
		if err != nil {
			log.Fatalf("Error while extracting addons... -> %v", err)
		}

		log.Println("Done with extracting all addons!")
	} else {
		log.Println("For some reason it did not finish properly...")
	}

	elapsedTime := time.Since(startTime)
	log.Printf("Elapsed duration: %s\n", elapsedTime)
}

func getDownloadPathAndAddonsPath(config []models.SystemConfig) (ap string, dp string, err error) {
	var addonsPath string
	var downloadPath string

	for _, c := range config {
		if c.Name == systemsAddonsPathName {
			if c.Path == "" {
				return "", "", fmt.Errorf("system's addons directory path is empty -> %v\n", err)
			}

			addonsPath = c.Path
			log.Printf("Found the host system's addons directory path at: '%s'\n", addonsPath)
		} else if c.Name == browserDownloadDirName {
			if c.Path == "" {
				return "", "", fmt.Errorf("system's download directory path is empty -> %v\n", err)
			}

			downloadPath = c.Path
			log.Printf("Found the host system's download directory path at: '%s'\n", downloadPath)
		}
	}

	return addonsPath, downloadPath, nil
}

func extractAllAddons(dp string, ap string) error {
	files, err := getAddonsAtPath(dp)
	if err != nil {
		return err
	}

	log.Printf("Will now commence extraction of a total of %d addons...", len(files))
	for i, file := range files {
		// Combine the directory path with the entry name to get the absolute path
		absPath := filepath.Join(dp, file.Name())
		log.Printf("[%d/%d] System extract of %s!", i+1, len(files), file.Name())
		err := system.Extract(absPath, ap)
		if err != nil {
			return err
		}

		id, err := database.InsertELog(db, models.ELog{
			RunId: strings.Replace(uuid.New().String(), "-", "", -1),
			File:  absPath,
		})
		if err != nil {
			log.Fatalf("Failed to insert extract log into database for the file: %s -> %v", file.Name(), err)
		}

		log.Printf("Stored extract logging with ID: %d, for the file: %s\n", id, file.Name())
	}

	return nil
}

func getAddonsAtPath(p string) (f []os.DirEntry, err error) {
	files, err := os.ReadDir(p)
	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files at: %s", p)
	}

	return files, nil
}
