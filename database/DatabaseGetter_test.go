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

package database

import (
	"encoding/json"
	"goaddons/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var (
	expectedTestAddon1 = models.Addon{
		Id:             1,
		Name:           "TestName1",
		Filename:       "TestFilename1",
		Url:            "www.TestUrl1.io",
		DownloadUrl:    "www.TestDownloadUrl1.io/download",
		LastDownloaded: nil,
		LastModifiedAt: nil,
		AddedAt:        nil,
	}
	expectedTestAddon2 = models.Addon{
		Id:             2,
		Name:           "TestName2",
		Filename:       "TestFilename2",
		Url:            "www.TestUrl2.io",
		DownloadUrl:    "www.TestDownloadUrl2.io/download",
		LastDownloaded: nil,
		LastModifiedAt: nil,
		AddedAt:        nil,
	}
	expectedSystemConfig1 = models.SystemConfig{
		Name: "TestConfig1",
		Path: "TestPath1",
	}
	expectedSystemConfig2 = models.SystemConfig{
		Name: "TestConfig2",
		Path: "TestPath2",
	}
)

func TestGetSystemConfigurations(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("database.TestGetSystemConfigurations :: An error was not expected when opening "+
			"a stub database connection -> %v\n", err)
	}

	rows := sqlmock.NewRows([]string{"Name", "Path"}).
		AddRow("TestConfig1", "TestPath1").
		AddRow("TestConfig2", "TestPath2")

	mock.ExpectQuery("SELECT \\* FROM system_config_default").WillReturnRows(rows)

	configs, err := GetSystemConfigurations(db)
	assert.NoError(t, err)
	assert.Len(t, configs, 2)

	expectedSystemConfig1Json, err := json.Marshal(expectedSystemConfig1)
	assert.NoError(t, err)

	expectedSystemConfig2Json, err := json.Marshal(expectedSystemConfig2)
	assert.NoError(t, err)

	actualSystemConfig1Json, err := json.Marshal(configs[0])
	assert.NoError(t, err)

	actualSystemConfig2Json, err := json.Marshal(configs[1])
	assert.NoError(t, err)

	assert.EqualValues(t, string(expectedSystemConfig1Json), string(actualSystemConfig1Json))
	assert.EqualValues(t, string(expectedSystemConfig2Json), string(actualSystemConfig2Json))

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("database.TestGetSystemConfigurations :: There were unfulfilled expectations: %s\n", err)
	}
}

func TestGetAllAddons(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("database.TestGetAllAddons :: An error was not expected when opening "+
			"a stub database connection -> %v\n", err)
	}

	rows := sqlmock.NewRows([]string{"Id", "Name", "Filename", "Url", "DownloadUrl", "LastDownloaded", "LastModifiedAt", "AddedAt"}).
		AddRow(1, "TestName1", "TestFilename1", "www.TestUrl1.io", "www.TestDownloadUrl1.io/download", nil, nil, nil).
		AddRow(2, "TestName2", "TestFilename2", "www.TestUrl2.io", "www.TestDownloadUrl2.io/download", nil, nil, nil)

	mock.ExpectQuery("SELECT id, name, filename, url, download_url, last_downloaded, last_modified_at, " +
		"added_at FROM addon WHERE download_url IS NOT NULL AND download_url != '';").
		WillReturnRows(rows)
	addons, err := GetAllAddons(db)

	assert.NoError(t, err)
	assert.Len(t, addons, 2)

	actualTestAddon1, err := json.Marshal(addons[0])
	assert.NoError(t, err)

	actualTestAddon2, err := json.Marshal(addons[1])
	assert.NoError(t, err)

	expectedTestAddon1Json, err := json.Marshal(expectedTestAddon1)
	assert.NoError(t, err)

	expectedTestAddon2Json, err := json.Marshal(expectedTestAddon2)
	assert.NoError(t, err)

	assert.EqualValues(t, string(expectedTestAddon1Json), string(actualTestAddon1))
	assert.EqualValues(t, string(expectedTestAddon2Json), string(actualTestAddon2))

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("database.TestGetAllAddons :: There were unfulfilled expectations: %s\n", err)
	}
}

func TestGetAddonsByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("database.TestGetAddonsByName :: An error was not expected when opening "+
			"a stub database connection -> %v\n", err)
		return
	}

	rows := sqlmock.NewRows([]string{"Id", "Name", "Filename", "Url", "DownloadUrl", "LastDownloaded", "LastModifiedAt", "AddedAt"}).
		AddRow(1, "TestName2", "TestFilename2", "www.TestUrl2.io", "www.TestDownloadUrl2.io/download", nil, nil, nil)

	mock.ExpectQuery("SELECT id, name, filename, url, download_url, last_downloaded, last_modified_at, added_at " +
		"FROM addon WHERE name LIKE '%" + expectedTestAddon2.Name + "%'").
		WillReturnRows(rows)
	addon, err := GetAddonsByName(db, expectedTestAddon2.Name)
	assert.NoError(t, err)
	assert.Len(t, addon, 1)

	actualTestAddon2, err := json.Marshal(addon[0])
	assert.NoError(t, err)

	expectedTestAddon2.Id = 1
	expectedTestAddon2Json, err := json.Marshal(expectedTestAddon2)
	assert.NoError(t, err)

	assert.EqualValues(t, string(expectedTestAddon2Json), string(actualTestAddon2))

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("database.TestGetAddonsByName :: There were unfulfilled expectations: %s\n", err)
	}
}
