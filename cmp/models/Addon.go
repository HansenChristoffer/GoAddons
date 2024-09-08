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

type AddonConfig struct {
	Elements []Addon `toml:"addon"`
}

type Addon struct {
	Name string `toml:"name"`
	Url  string `toml:"url"`
}

func (a *Addon) Validate() bool {
	if len(a.Name) == 0 || len(a.Url) == 0 {
		return false
	}
	return true
}
