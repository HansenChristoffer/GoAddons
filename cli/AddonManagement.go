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
	"encoding/json"
	"fmt"
	"goaddons/database"
	"goaddons/models"
	"goaddons/utils"
	"log"
	"strconv"
)

func ListAllAddons() {
	utils.ClearScreen()
	fmt.Printf("  »»» All Addons «««\n\n")

	addons, err := database.GetAllAddons(db)
	if err != nil {
		log.Printf("cli.ListAllAddons :: Error while trying to get all addons -> %v\n", err)
		return
	}

	if addons == nil || len(addons) == 0 {
		return
	}

	b, err := json.MarshalIndent(addons, "", "  ")
	if err != nil {
		log.Printf("cli.ListAllAddons :: Failed to marshal ->  %v\n", err)
		return
	}

	if b != nil && utils.IsValidString(string(b)) {
		fmt.Println(string(b))
	}
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
	addons, err := database.GetAddonsByName(db, name)
	if err != nil {
		log.Printf("cli.SearchForAddonByName :: Error while getting addon... -> %v\n", err)
		return
	}

	if addons == nil || len(addons) == 0 {
		return
	}

	b, err := json.MarshalIndent(addons, "", "  ")
	if err != nil {
		log.Printf("cli.SearchForAddonByName :: Failed to marshal ->  %v\n", err)
		return
	}

	if b != nil && utils.IsValidString(string(b)) {
		fmt.Println(string(b))
	}
}

func AddNewAddon() {
	var addon models.Addon

	fmt.Printf("\n  »»» Insert new addon ««« \n\n Addon name:\n > ")
	addon.Name = userInput("cli.AddNewAddon")

	fmt.Printf("\n What is the extracted addon directory is called\n > ")
	addon.Filename = userInput("cli.AddNewAddon")

	fmt.Printf("\n Addon about URL\n > ")
	addon.Url = userInput("cli.AddNewAddon")

	fmt.Printf("\n Addon download URL\n > ")
	addon.DownloadUrl = userInput("cli.AddNewAddon")

	fmt.Printf("\n Do you want to commit? [y/N]\n > ")
	input := userInput("cli.AddNewAddon")

	switch input {
	case "Y", "y":
		r, err := database.InsertAddon(db, addon)
		if err != nil {
			log.Printf("cli.addNewAddon :: Failed to insert addon(s)! -> %v\n", err)
		}
		fmt.Printf(" Inserted total of %d addon(s) into TanukiDB!", r)
	case "N", "n":
		fmt.Printf(" Stopped insertion of new addon(s)!")
	default:
		fmt.Printf(" Stopped insertion of new addon(s)!")
	}
}

func RemoveAddon() {
	utils.ClearScreen()

	fmt.Printf("\n  »»» Remove addon «««\n\n Addon ID\n > ")
	id := userInput("cli.RemoveAddon")

	if !utils.IsValidString(id) {
		log.Printf("cli.RemoveAddon :: The 'ID' argument is not valid! -> [%s]\n", id)
		return
	}

	idNum, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("cli.RemoveAddon :: Error while trying to convert ID to its decimal equivalent! -> %v\n", err)
		return
	}

	fmt.Printf("\n Do you want to commit? [y/N]\n > ")
	input := userInput("cli.RemoveAddon")

	switch input {
	case "Y", "y":
		r, err := database.RemoveAddonByID(db, idNum)
		if err != nil {
			log.Printf("cli.RemoveAddon :: Failed to remove addon(s)! -> %v\n", err)
		}
		fmt.Printf(" Removed total of %d addon(s) from TanukiDB!", r)
	case "N", "n":
		fmt.Printf(" Stopped deletion of addon(s)!")
	default:
		fmt.Printf(" Stopped deletion of addon(s)!")
	}
}

func userInput(caller string) (input string) {
	_, err := fmt.Scanln(&input)
	if err != nil {
		log.Printf(caller+" :: Error occurred while interpreting input from user -> %v\n", err)
	}
	return
}
