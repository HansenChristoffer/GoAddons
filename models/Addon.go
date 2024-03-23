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

type Addon struct {
	Id             int     `json:"id,omitempty"`
	Name           string  `json:"name,omitempty"`
	Filename       string  `json:"filename,omitempty"`
	Url            string  `json:"url,omitempty"`
	DownloadUrl    string  `json:"download_url,omitempty"`
	LastDownloaded []uint8 `json:"last_downloaded,omitempty"`
	LastModifiedAt []uint8 `json:"last_modified_at,omitempty"`
	AddedAt        []uint8 `json:"addedAt,omitempty"`
}
