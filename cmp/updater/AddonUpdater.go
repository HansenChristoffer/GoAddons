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
	"fmt"
	"goaddons/cmp/models"
	"goaddons/cmp/net"
	"goaddons/cmp/utils"
	"log"
	"time"
)

const (
	SYSTEMS_ADDON_PATH_KEY    = "extract.addon.path"
	BROWSER_DOWNLOAD_PATH_KEY = "browser.download.dir"
)

func StartUpdater() {
	log.Printf("Addon Updater starting...\n")

	// Fetch systems configuration from file
	systemVault, err := models.GetSystemVaultInstance()
	if err != nil {
		log.Fatalf("SystemVault could not be initialized: %v", err)
	}

	if systemVault.ConfigElements == nil || len(systemVault.ConfigElements) == 0 {
		log.Fatalf("Failed to find any system configurations! -> %v\n", err)
	}

	// Fetch the addons path & download dir from the system configurations
	addonsPath, downloadPath, err := GetSystemPaths(systemVault.ConfigElements)
	if err != nil {
		log.Fatalf("Failed to parse system configurations for addon and download paths! -> %v\n", err)
	}
	addonsVault, err := models.GetAddonVaultInstance()
	if err != nil {
		log.Fatalf("Error initializing AddonVault: %v", err)
	}

	startTime := time.Now().UTC()

	log.Printf("Vault has a total of %d addons!\n", addonsVault.Length())

	// Starts headless browser and downloads addons. Will return bool for done state
	done, err := net.StartHeadlessAndDownloadAddons(addonsVault, downloadPath)
	if err != nil {
		log.Fatalf("Error while navigating... -> %v\n", err)
	}

	if done {
		log.Printf("Addons downloaded successfully!\n")
		err = Extractor(downloadPath, addonsPath)
		if err != nil {
			log.Printf("Error extracting addons! -> %v\n", err)
		}
	}

	elapsedTime := time.Since(startTime)
	log.Printf("Elapsed duration: %s\n", elapsedTime)
	utils.PressEnterToReturn()
}

func GetSystemPaths(config []models.ConfigElement) (addonPath string, downloadPath string, err error) {
	for _, element := range config {
		if element.Name == SYSTEMS_ADDON_PATH_KEY {
			if !utils.IsValidString(element.Value) {
				return "", "",
					fmt.Errorf("AddonUpdater -> GetSystemPaths(): system's addons directory path is empty\n")
			}

			if !utils.PathExists(element.Value) {
				return "", "", fmt.Errorf("AddonUpdater -> GetSystemPaths(): " +
					"system's addons directory path is empty")
			}

			addonPath = element.Value
			log.Printf("Found the host system's addons directory path at: '%s'\n", addonPath)
		} else if element.Name == BROWSER_DOWNLOAD_PATH_KEY {
			if element.Value == "" {
				return "", "", fmt.Errorf("AddonUpdater -> GetSystemPaths(): " +
					"browser's download directory path is empty\n")
			}

			if !utils.PathExists(element.Value) {
				return "", "", fmt.Errorf("AddonUpdater -> GetSystemPaths(): " +
					"system's addons directory path is empty\n")
			}
			downloadPath = element.Value
			log.Printf("Found the host system's download directory path at: '%s'\n", downloadPath)
		}
	}
	return
}
