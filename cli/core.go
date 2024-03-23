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
	"database/sql"
	"fmt"
	"goaddons/database"
	"goaddons/updater"
	"goaddons/utils"
	"os"
)

var db *sql.DB

func StartCli() {
	utils.ClearScreen()

	fmt.Printf(`  

  »»» GoAddons «««

 1. Addon Management
 2. Updater Menu
 3. About
 X. Exit
 > `)

	input := userInput("cli.StartCli")
	if len(input) > 0 {
		switch input {
		case "1":
			addonManagement()
		case "2":
			updaterMenu()
		case "3":
			about()
		case "X", "x":
			terminate()
		default:
			StartCli()
		}
	}

	StartCli()
}

func addonManagement() {
	utils.ClearScreen()
	tryDatabaseConnect()

	fmt.Printf(`

  »»» Addon Management «««

 1. List all addons
 2. Search for addon
 3. Add addon
 4. Remove addon
 X. Go back...
 > `)

	input := userInput("cli.addonManagement")
	if len(input) > 0 {
		switch input {
		case "1":
			ListAllAddons()
			utils.PressEnterToReturn(addonManagement)
		case "2":
			SearchForAddonByName()
			utils.PressEnterToReturn(addonManagement)
		case "3":
			AddNewAddon()
			utils.PressEnterToReturn(addonManagement)
		case "4":
			RemoveAddon()
			utils.PressEnterToReturn(addonManagement)
		case "X", "x":
			StartCli()
		default:
			addonManagement()
		}
	}

	addonManagement()
}

func updaterMenu() {
	utils.ClearScreen()

	fmt.Printf(`  
  »»» Updater Menu «««

 1. Start updater
 X. Go back...
 > `)

	input := userInput("cli.updaterMenu")
	if len(input) > 0 {
		switch input {
		case "1":
			updater.StartUpdater()
		case "X", "x":
			StartCli()
		default:
			updaterMenu()
		}
	}

	updaterMenu()
}

func about() {
	utils.ClearScreen()

	fmt.Printf(`
  »»» About GoAddons «««

GoAddons is a state-of-the-art command-line interface (CLI) application designed to revolutionize the way World of Warcraft (WoW) enthusiasts manage their addons. 
Developed with precision and a deep understanding of the gamer's needs, GoAddons offers an unparalleled user experience, allowing for effortless management, updating, and discovery of WoW addons.

At the heart of GoAddons is the integration with TanukiDB, a comprehensive database of WoW addons. This powerful synergy enables users to seamlessly search, add, and remove addons, ensuring their gaming setup is always optimized for victory.

Key Features:
- Addon Management: Curate your collection of WoW addons with simple commands. List, search, add, or remove addons with ease, tailoring your addon library to your gaming needs.
- Updater Menu: Stay ahead of the game with the Updater feature. GoAddons checks for the latest versions of your addons and updates them automatically, ensuring you're always equipped with the latest tools and enhancements.
- About GoAddons: Learn about the philosophy, features, and the team behind GoAddons. We're committed to providing an exceptional tool that enhances your gaming experience.

GoAddons is more than just a tool; it's a companion for every WoW player who believes in the power of customization and efficiency. Developed by a team of passionate gamers and skilled engineers, GoAddons is dedicated to elevating your WoW experience to new heights.

Thank you for choosing GoAddons. Your adventure awaits.

Press ENTER to return to the main menu...
`)
	_, err := fmt.Scanln()
	if err != nil {
		return
	}
	StartCli()
}

func tryDatabaseConnect() {
	if db == nil {
		db = database.ConnectToServer()
	}
}

func terminate() {
	os.Exit(0)
}
