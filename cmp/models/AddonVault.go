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

package models

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"goaddons/cmp/utils"
	"log"
	"os"
	"sync"
)

const ADDONS_CONFIG_FILENAME = "./.addonsconfig.toml"

var (
	addonVaultInstance *AddonVault
	addonVaultOnce     sync.Once
)

type AddonVault struct {
	Addons   []Addon `toml:"addons"`
	structMx sync.RWMutex
	fileMx   sync.Mutex
}

// GetAddonVaultInstance returns the singleton systemVaultInstance of AddonVault
func GetAddonVaultInstance() (*AddonVault, error) {
	addonVaultOnce.Do(func() {
		log.Println("AddonVault systemVaultInstance created")
		addonVaultInstance = &AddonVault{}
		err := addonVaultInstance.LoadAddonVault()
		if err != nil {
			log.Printf("Failed to load AddonVault: %v", err)
			addonVaultInstance = nil // reset systemVaultInstance to nil on failure
		}
	})
	if addonVaultInstance == nil {
		return nil, errors.New("failed to initialize AddonVault")
	}
	return addonVaultInstance, nil
}

func (v *AddonVault) Append(addon *Addon) error {
	if addon == nil {
		return errors.New("AddonVault -> Append(): argument is nil")
	}

	if !addon.Validate() {
		return errors.New("addonVault -> Append(): addon is invalid")
	}

	v.structMx.Lock()
	v.Addons = append(v.Addons, *addon)
	v.structMx.Unlock()

	if err := v.PersistAddonVault(); err != nil {
		return fmt.Errorf("AddonVault -> Append(): error persisting data: %v", err)
	}
	return nil
}

func (v *AddonVault) AppendMany(addons []Addon) error {
	if addons == nil {
		return nil
	}

	for _, addon := range addons {
		if !addon.Validate() {
			return errors.New("addonVault -> AppendMany(): addon is invalid")
		}
	}

	v.structMx.Lock()
	v.Addons = append(v.Addons, addons...)
	v.structMx.Unlock()

	if err := v.PersistAddonVault(); err != nil {
		return fmt.Errorf("AddonVault -> AppendMany(): error persisting data: %v", err)
	}
	return nil
}

func (v *AddonVault) Remove(id string) (bool, error) {
	if !utils.IsValidString(id) {
		return false, errors.New("AddonVault -> Remove(): argument is invalid")
	}

	v.structMx.Lock()
	var found bool

	// Find and remove the addon
	for i, addon := range v.Addons {
		if addon.Name == id {
			v.Addons = append(v.Addons[:i], v.Addons[i+1:]...)
			found = true
			break
		}
	}
	v.structMx.Unlock() // Unlock before calling PersistAddonVault()

	if !found {
		return false, fmt.Errorf("AddonVault -> Remove(): addon with id '%s' not found", id)
	}

	// PersistAddonVault the updated list outside the lock to avoid deadlock
	if err := v.PersistAddonVault(); err != nil {
		return false, err
	}

	return true, nil
}

func (v *AddonVault) LoadAddonVault() error {
	v.fileMx.Lock()
	file, err := os.Open(ADDONS_CONFIG_FILENAME)
	if err != nil {
		return fmt.Errorf("AddonVault -> LoadAddonVault(): %v", err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Printf("AddonVault -> LoadAddonVault(): error closing file: %v", err)
		}
	}(file)
	v.fileMx.Unlock()

	// Decode into AddonConfig, not []Addon
	var addonConfig AddonConfig
	decoder := toml.NewDecoder(file)
	_, err = decoder.Decode(&addonConfig)
	if err != nil {
		addonConfig.Elements = make([]Addon, 0)
		return nil
	}

	return v.AppendMany(addonConfig.Elements)
}

func (v *AddonVault) PersistAddonVault() error {
	v.fileMx.Lock()
	defer v.fileMx.Unlock()

	v.structMx.RLock()
	addonsCopy := make([]Addon, len(v.Addons))
	copy(addonsCopy, v.Addons)
	v.structMx.RUnlock()

	data, err := toml.Marshal(&AddonConfig{Elements: addonsCopy})
	if err != nil {
		return fmt.Errorf("AddonVault -> SaveToFile(): error marshaling to TOML: %v", err)
	}

	// Open the file for writing, with permissions to create or truncate
	file, err := os.OpenFile(ADDONS_CONFIG_FILENAME, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("AddonVault -> SaveToFile(): error opening file: %v", err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Printf("AddonVault -> SaveToFile(): error closing file: %v", err)
		}
	}(file)

	// Write the TOML data to the file
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("AddonVault -> SaveToFile(): error writing to file: %v", err)
	}

	return nil
}

func (v *AddonVault) Length() int {
	v.structMx.RLock()
	defer v.structMx.RUnlock()
	return len(v.Addons)
}

func (v *AddonVault) PrintAddonVault() {
	v.structMx.RLock()
	defer v.structMx.RUnlock()

	if v.Addons == nil {
		return
	}

	bytes, err := toml.Marshal(&AddonConfig{Elements: v.Addons})
	if err != nil {
		log.Printf("Error marshalling addon to TOML: %v", err)
	}

	if len(bytes) > 0 {
		fmt.Println(string(bytes))
	}
}
