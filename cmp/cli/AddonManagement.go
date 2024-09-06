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

package cli

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"goaddons/cmp/models"
	"goaddons/cmp/utils"
	"log"
)

func ListAllAddons() {
	utils.ClearScreen()
	fmt.Printf("  »»» All Addons «««\n\n")

	addonsVault, err := models.GetAddonVaultInstance()
	if err != nil {
		log.Fatalf("Error initializing AddonVault: %v", err)
	}

	addonsVault.PrintAddonVault()
}

func SearchForAddonByName() {
	utils.ClearScreen()

	fmt.Printf("  »»» Addon Searching «««\n\n Addon name\n > ")
	name := userInput("cli.addonManagement")

	if !utils.IsValidString(name) {
		log.Printf("cli.SearchForAddonByName :: The 'name' argument is not valid! -> [%s]\n", name)
		return
	}

	// Search for and list addon, all by the argument 'name'
	addonsVault, err := models.GetAddonVaultInstance()
	if err != nil {
		log.Fatalf("Was unable to get instance of AddonVault")
	}
	if addonsVault.Length() == 0 {
		return
	}

	for _, addon := range addonsVault.Addons {
		if addon.Name == name {
			b, err := toml.Marshal(&models.AddonConfig{Elements: []models.Addon{addon}})
			if err != nil {
				log.Printf("Error marshalling addon to TOML: %v", err)
			}

			if b != nil && len(b) > 0 {
				fmt.Println(string(b))
			}
		}
	}
}

func AddNewAddon() {
	const CALLER = "cli.AddNewAddon"
	var addon models.Addon

	fmt.Printf("\n  »»» Insert new addon ««« \n\n Addon name:\n > ")
	addon.Name = userInput(CALLER)

	fmt.Printf("\n Addon about URL\n > ")
	addon.Url = userInput(CALLER)

	fmt.Printf("\n »»» New Addon ««« \n")
	fmt.Printf(" Name: %s\n URL: %s\n", addon.Name, addon.Url)
	fmt.Printf("\n Do you want to commit? [y/N]\n > ")
	input := userInput(CALLER)

	switch input {
	case "Y", "y":
		addonsVault, err := models.GetAddonVaultInstance()
		if err != nil {
			log.Fatalf("Error initializing AddonVault: %v", err)
		}

		if err := addonsVault.Append(&addon); err != nil {
			log.Printf("Error appending addon to vault: %v", err)
			return
		}
		fmt.Printf(" Successfully inserted addon(s)!")
	case "N", "n":
		fmt.Println(" Stopped insertion of new addon(s)!")
	default:
		fmt.Println(" Stopped insertion of new addon(s)!")
	}
}

func RemoveAddon() {
	utils.ClearScreen()

	fmt.Printf("\n  »»» Remove addon «««\n\n Addon name\n > ")
	name := userInput("cli.RemoveAddon")

	if !utils.IsValidString(name) {
		log.Printf("cli.RemoveAddon :: The 'name' argument is not valid! -> [%s]\n", name)
		return
	}

	fmt.Printf("\n Will try to remove addon with the name: %s\n", name)
	fmt.Printf("\n Do you want to commit? [y/N]\n > ")
	input := userInput("cli.RemoveAddon")
	switch input {
	case "Y", "y":
		addonsVault, err := models.GetAddonVaultInstance()
		if err != nil {
			log.Fatalf("Error initializing AddonVault: %v", err)
		}
		beenRemoved, err := addonsVault.Remove(name)
		if beenRemoved {
			fmt.Printf(" Successfully removed addon with name: %s", name)
		} else {
			fmt.Printf(" Failed to remove addon with name: %s -> %v\n", name, err)
		}
	case "N", "n":
		fmt.Println(" Stopped deletion of addon(s)!")
	default:
		fmt.Println(" Stopped deletion of addon(s)!")
	}
}

func userInput(caller string) (input string) {
	_, err := fmt.Scanln(&input)
	if err != nil {
		log.Printf(caller+" :: Error occurred while interpreting input from user -> %v\n", err)
	}
	return
}
