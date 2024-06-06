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
	"goaddons/net"
	"goaddons/utils"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	serviceName            = "GOADDONS_SERVICE"
	systemsAddonsPathName  = "extract.addon.path"
	browserDownloadDirName = "browser.download.dir"
)

func StartUpdater(db *sql.DB) {
	log.Printf("Addon Updater starting...\n")

	// Connect to MySQL database
	if db == nil {
		var err error
		db, err = database.ConnectToServer("./bin/acd.db")
		if err != nil {
			log.Fatalf("updater:AddonUpdater.StartUpdater():database.ConnectToServer(\"./bin/acd\") -> %v", err)
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				log.Println("Error occurred while trying to close connection to database -> %w", err)
			}
		}(db)
	}

	runId := strings.Replace(uuid.New().String(), "-", "", -1)

	// Log that the software is making a "run"
	rLog, err := database.InsertRLog(db, models.RLog{RunId: runId,
		Service: serviceName})
	if err != nil {
		log.Printf("Failed to store run logging into database... -> %v\n", err)
	}

	log.Printf("Stored run logging with ID: %d\n", rLog)

	// Fetch systems configuration from database
	config, err := database.GetSystemConfigurations(db)
	if err != nil {
		log.Fatalf("Failed to fetch system configurations... -> %v\n", err)
	}

	// Fetch the addons path & download dir from the system configurations
	addonsPath, downloadPath, err := getDownloadPathAndAddonsPath(config)
	if err != nil {
		log.Fatalf("Failed to fetch system configurations... -> %v\n", err)
	}

	startTime := time.Now().UTC()

	// Get all addons from database
	addons, err := database.GetAllAddons(db)
	if err != nil {
		log.Fatalf("Failed to fetch all addons... -> %v\n", err)
	}

	log.Printf("Fetched a total of %d addons!\n", len(addons))

	// Creates stop channel and errorChannel and starts the PollingExtractor
	stopChannel := make(chan struct{})
	errorChannel := make(chan error, 1)
	go PollingExtractor(db, runId, downloadPath, addonsPath, stopChannel, errorChannel)

	go func() {
		for err := range errorChannel {
			log.Println("Error from PollingExtractor:", err)
		}
	}()

	// Starts headless browser and downloads addons. Will return bool for done state
	done, err := net.StartHeadlessAndDownloadAddons(runId, addons, downloadPath, db)
	if err != nil {
		log.Fatalf("Error while navigating... -> %v\n", err)
	}

	if done {
		handleDone(stopChannel, errorChannel)
	} else {
		handleNotDone(stopChannel, errorChannel)
	}

	elapsedTime := time.Since(startTime)
	log.Printf("Elapsed duration: %s\n", elapsedTime)
	utils.PressEnterToReturn()
}

func handleDone(stopChannel chan struct{}, errorChannel chan error) {
	log.Printf("Done, closing down resources...\n")
	time.Sleep(5 * time.Second)

	err := <-errorChannel
	if err != nil {
		log.Fatalf("Error while PollingExtractor running... -> %v\n", err)
	}

	time.Sleep(500 * time.Millisecond)
	close(stopChannel)
}

func handleNotDone(stopChannel chan struct{}, errorChannel chan error) {
	log.Println("For some reason it did not finish properly...")

	err := <-errorChannel
	if err != nil {
		log.Fatalf("Error while PollingExtractor running... -> %v\n", err)
	}

	time.Sleep(500 * time.Millisecond)
	close(stopChannel)
}

func getDownloadPathAndAddonsPath(config []models.SystemConfig) (ap string, dp string, err error) {
	for _, c := range config {
		if c.Name == systemsAddonsPathName {
			if c.Path == "" {
				return "", "", fmt.Errorf("system's addons directory path is empty -> %v\n", err)
			}

			ap = c.Path
			log.Printf("Found the host system's addons directory path at: '%s'\n", ap)
		} else if c.Name == browserDownloadDirName {
			if c.Path == "" {
				return "", "", fmt.Errorf("system's download directory path is empty -> %v\n", err)
			}

			dp = c.Path
			log.Printf("Found the host system's download directory path at: '%s'\n", dp)
		}
	}
	return
}
