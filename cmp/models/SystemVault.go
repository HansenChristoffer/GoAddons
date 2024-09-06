package models

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"sync"
)

const SYSTEMS_CONFIG_FILENAME = "./.systemsconfig.toml"

var (
	systemVaultInstance *SystemVault
	systemVaultOnce     sync.Once
)

type SystemVault struct {
	ConfigElements      []ConfigElement `toml:"config_elements"`
	systemVaultStructMx sync.RWMutex
	systemConfigFileMx  sync.Mutex
}

func GetSystemVaultInstance() (*SystemVault, error) {
	systemVaultOnce.Do(func() {
		log.Println("SystemVault instance created")
		systemVaultInstance = &SystemVault{}
		err := systemVaultInstance.LoadSystemVault()
		if err != nil {
			log.Printf("Failed to load SystemVault: %v", err)
			systemVaultInstance = nil // reset systemVaultInstance to nil on failure
		}
	})
	if systemVaultInstance == nil {
		return nil, errors.New("failed to initialize SystemVault")
	}
	return systemVaultInstance, nil
}

func (v *SystemVault) AppendManyToSystemConfig(configs []ConfigElement) {
	if configs == nil {
		return
	}

	v.systemVaultStructMx.Lock()
	defer v.systemVaultStructMx.Unlock()
	v.ConfigElements = append(v.ConfigElements, configs...)
}

func (v *SystemVault) GetConfigByName(name string) *ConfigElement {
	v.systemVaultStructMx.RLock()
	defer v.systemVaultStructMx.RUnlock()

	for _, config := range v.ConfigElements {
		if config.Name == name {
			return &config
		}
	}
	return nil
}

func (v *SystemVault) LoadSystemVault() error {
	v.systemConfigFileMx.Lock()
	file, err := os.Open(SYSTEMS_CONFIG_FILENAME)
	if err != nil {
		return fmt.Errorf("SystemVault -> LoadSystemVault(): %v", err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Printf("SystemVault -> LoadSystemVault(): error closing file: %v", err)
		}
	}(file)
	v.systemConfigFileMx.Unlock()

	// Decode into AddonConfig, not []Addon
	var systemConfig SystemConfig
	decoder := toml.NewDecoder(file)
	_, err = decoder.Decode(&systemConfig)
	if err != nil {
		systemConfig.Elements = make([]ConfigElement, 0)
		return nil
	}

	v.AppendManyToSystemConfig(systemConfig.Elements)
	return nil
}

func (v *SystemVault) PrintSystemVault() {
	v.systemVaultStructMx.RLock()
	defer v.systemVaultStructMx.RUnlock()

	if v.ConfigElements == nil {
		return
	}

	bytes, err := toml.Marshal(&SystemConfig{Elements: v.ConfigElements})
	if err != nil {
		log.Printf("Error marshalling addon to TOML: %v", err)
	}

	if len(bytes) > 0 {
		fmt.Println(string(bytes))
	}
}
