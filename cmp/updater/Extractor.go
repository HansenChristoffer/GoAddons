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
	"goaddons/cmp/system"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Extractor(downloadPath string, addonPath string) error {
	log.Printf("Extractor starting...")

	compressedAddonFiles, found, err := GetCompressedAddonsAtPath(downloadPath)
	if err != nil {
		return err
	}

	if found {
		err = extractFiles(compressedAddonFiles, downloadPath, addonPath)
		if err != nil {
			return err
		}
	}

	time.Sleep(2 * time.Second)
	log.Printf("Extractor is done!")
	return nil
}

// extractFiles processes and extracts the filtered addon files
func extractFiles(files []os.DirEntry, downloadPath string, addonPath string) error {
	log.Printf("Will now commence extraction of a total of %d addons...\n", len(files))
	for idx, file := range files {
		absPath := filepath.Join(downloadPath, file.Name())

		log.Printf("[%d/%d] System extract of %s!\n", idx+1, len(files), file.Name())
		if err := system.Extract(absPath, addonPath); err != nil {
			return fmt.Errorf("extractor -> extractFiles(): %w\n", err)
		}

		log.Printf("[%d/%d] Removing compressed addon file now that we're done with it: %s\n",
			idx+1, len(files), file.Name())
		if err := os.Remove(absPath); err != nil {
			return fmt.Errorf("updater:Extractor.extractFiles():os.Remove(%s) -> %w", absPath, err)
		}
	}
	return nil
}

// GetCompressedAddonsAtPath retrieves the list of directory entries at the specified path that contains .zip
func GetCompressedAddonsAtPath(path string) (files []os.DirEntry, found bool, err error) {
	readFiles, err := os.ReadDir(path)
	if err != nil {
		return nil, false, fmt.Errorf("extractor -> GetCompressedAddonsAtPath(): %w\n", err)
	}
	for _, file := range readFiles {
		if strings.HasSuffix(file.Name(), ".zip") {
			files = append(files, file)
		}
	}
	return files, len(files) > 0, nil
}
