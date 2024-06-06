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

package main

import (
	"database/sql"
	"fmt"
	"goaddons/cli"
	"goaddons/updater"
	"goaddons/version"
	"log"
	"os"
)

var db *sql.DB

func main() {
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "--updater", "-u":
			updater.StartUpdater(db)
		case "--cli", "c":
			cli.StartCli()
		case "--version", "-v":
			fmt.Printf("GoAddons Version: %s, Build Date: %s, Commit: %s\n",
				version.Version, version.BuildDate, version.Commit)
		case "--help", "-h":
			fmt.Println(`Usage: [command]

Commands:
  --updater, -u    Start the updater process.
  --cli, -c        Start the command line interface.
  --help, -h       Display this help message.

Description:
  This program provides a command line interface and an updater. Use the commands listed above to interact with the program or to start specific components of it. If no command is provided, the CLI will start by default.

Examples:
  To start the CLI:
    program_name --cli

  To start the updater:
    program_name --updater

  To display help:
    program_name --help`)
		default:
			log.Printf("Unknown parameter. Please use --help or -h for help!\n")
		}
	} else {
		cli.StartCli()
	}
}
