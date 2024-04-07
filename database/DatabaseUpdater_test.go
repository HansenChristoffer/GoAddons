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
	"goaddons/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateAddon(t *testing.T) {
	// Expected: "UPDATE kaasufouji.addons SET last_downloaded = ? WHERE name = ?", time.Now().UTC(), addon.Name
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	addon := models.Addon{Name: "TestAddon"}

	mock.ExpectExec("UPDATE addon SET last_downloaded = \\? WHERE name = \\?").
		WithArgs(sqlmock.AnyArg(), addon.Name).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Assuming the update affects 1 row

	r, err := UpdateAddon(db, addon)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), r)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestInsertAddon(t *testing.T) {
	// Expected: "INSERT IGNORE INTO kaasufouji.addons (name, filename, url, download_url) VALUES (?, ?, ?, ?);"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	addon := models.Addon{Name: "TestAddon", Filename: "testfile.zip", Url: "http://example.com/test", DownloadUrl: "http://example.com/download/testfile.zip"}

	mock.ExpectExec("INSERT OR IGNORE INTO addon \\(name, filename, url, download_url\\) VALUES \\(\\?, \\?, \\?, \\?\\);").
		WithArgs(addon.Name, addon.Filename, addon.Url, addon.DownloadUrl).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Assuming the insert results in 1 row affected

	r, err := InsertAddon(db, addon)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), r)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestRemoveAddonByID(t *testing.T) {
	// Expected: "DELETE FROM kaasufouji.addons WHERE id = ?;", id
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	addonID := 1

	mock.ExpectExec("DELETE FROM addon WHERE id = \\?;").
		WithArgs(addonID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Assuming the delete affects 1 row

	r, err := RemoveAddonByID(db, addonID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), r)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}
