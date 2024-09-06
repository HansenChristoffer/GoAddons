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
	"fmt"
	"github.com/BurntSushi/toml"
	"goaddons/cmp/models"
	"goaddons/cmp/utils"
	"io"
	"log"
	"os"
)

const (
	SYSTEMS_CONFIG_PATH   = "./.systemsconfig.toml"
	ADDONS_CONFIG_PATH    = "./.addonsconfig.toml"
	USERAGENT_CONFIG_PATH = "./.useragentconfig.toml"
)

// GetSystemsConfig fetches the system configurations from the TOML file
// Deprecated: because SystemVault directly loads with file, therefore this won't be required
func GetSystemsConfig() ([]models.ConfigElement, error) {
	var config models.SystemConfig
	configBytes, err := readConfig(SYSTEMS_CONFIG_PATH)
	if err != nil {
		return config.Elements, err
	}

	if err = toml.Unmarshal(configBytes, &config); err != nil {
		return nil, fmt.Errorf("configReader -> GetSystemsConfig(): %w\n", err)
	}
	return config.Elements, nil
}

// GetAddonsConfig fetches the addon configurations from the TOML file
// Deprecated: because AddonVault directly loads and persists with file, therefore this won't be required
func GetAddonsConfig() ([]models.Addon, error) {
	var addons models.AddonConfig
	addonsConfigBytes, err := readConfig(ADDONS_CONFIG_PATH)
	if err != nil {
		return addons.Elements, err
	}

	err = toml.Unmarshal(addonsConfigBytes, &addons)
	if err != nil {
		return nil, fmt.Errorf("configReader -> GetAddonsConfig(): %w\n", err)
	}
	return addons.Elements, nil
}

func GetUserAgentConfig() ([]models.UserAgent, error) {
	var userAgentConfig models.UserAgentConfig
	userAgentConfigBytes, err := readConfig(USERAGENT_CONFIG_PATH)
	if err != nil {
		return userAgentConfig.Elements, err
	}

	err = toml.Unmarshal(userAgentConfigBytes, &userAgentConfig)
	if err != nil {
		return nil, fmt.Errorf("configReader -> GetUserAgentConfig(): %w\n", err)
	}
	return userAgentConfig.Elements, nil
}

// ReadConfig will open, read and return a whole file as a []byte
func readConfig(fp string) ([]byte, error) {
	if !utils.IsValidString(fp) {
		return []byte{}, nil
	}

	f, err := os.OpenFile(fp, os.O_RDONLY, 0666)
	if err != nil {
		return []byte{}, fmt.Errorf("configReader -> ReadConfig(): %w\n", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("configReader -> ReadConfig(): %v\n", err)
			return
		}
	}(f)

	rtnBytes, err := io.ReadAll(f)
	if err != nil {
		return []byte{}, fmt.Errorf("configReader -> ReadConfig(): %w\n", err)
	}
	return rtnBytes, nil
}
