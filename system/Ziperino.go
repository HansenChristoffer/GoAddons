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

package system

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Extract(source string, destination string) error {
	if source == "" || destination == "" {
		return fmt.Errorf("source and destination is required")
	}

	return unzipSource(source, destination)
}

func unzipSource(source string, destination string) error {
	read, err := zip.OpenReader(source)
	if err != nil {
		return err
	}

	defer func(read *zip.ReadCloser) {
		err := read.Close()
		if err != nil {
			log.Printf("unzipSource.readClose: %v\n", err)
		}
	}(read)

	destination, err = filepath.Abs(destination)
	if err != nil {
		return err
	}

	log.Printf("Will now try to extract: %s to the destination: %s\n", source, destination)

	for _, f := range read.File {
		err := unzipFile(f, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(f *zip.File, destination string) error {
	if f == nil {
		return fmt.Errorf("zip file is not allowed to be nil")
	}

	sourceFilePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(sourceFilePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", sourceFilePath)
	}

	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(sourceFilePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(sourceFilePath), os.ModePerm); err != nil {
		return err
	}

	destinationFile, err := os.OpenFile(sourceFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}

	defer func(destinationFile *os.File) {
		err := destinationFile.Close()
		if err != nil {
			log.Printf("Failed to close destinationFile -> %v\n", err)
		}
	}(destinationFile)

	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer func(zippedFile io.ReadCloser) {
		err := zippedFile.Close()
		if err != nil {
			log.Printf("Failed to close zippedFile -> %v\n", err)
		}
	}(zippedFile)

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}

	return nil
}
