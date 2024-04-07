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

package net

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func DownloadFile(source string, destination string) (done bool, err error) {
	// Create the file
	out, err := os.Create(destination)
	if err != nil {
		return false, err
	}

	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Printf("DownloadFile.outClose: %v\n", err)
		}
	}(out)

	// Get the data
	resp, err := http.Get(source)
	if err != nil {
		return false, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("DownloadFile.BodyClose: %v\n", err)
		}
	}(resp.Body)

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("bad status: %s\n", resp.Status)
	}

	// Writer the body to file
	i, err := io.Copy(out, resp.Body)
	if err != nil {
		return false, err
	}

	if i == 0 {
		return false, nil
	}

	return true, nil
}
