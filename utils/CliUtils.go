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

package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type function func()

// ClearScreen clears the CLI screen based on the operating system.
func ClearScreen() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		// Assuming Unix-like system
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		// Handle the error if needed
		fmt.Println("utils.ClearScreen :: Error clearing screen -> %w\n", err)
		return
	}
}

// PressEnterToReturnToFunction Waits for user to press the return key then calls function
func PressEnterToReturnToFunction(fn function) {
	fmt.Printf("\n\n Press ENTER to return...\n")
	_, err := fmt.Scanln()
	if err != nil {
		return
	}
	fn()
}

// PressEnterToReturn Waits for user to press the return key
func PressEnterToReturn() {
	fmt.Printf("\n\n Press ENTER to return...\n")
	_, err := fmt.Scanln()
	if err != nil {
		return
	}
}
