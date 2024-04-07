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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var testUUID, _ = uuid.NewUUID()

var (
	expectedRLog1 = models.RLog{
		Id:      1,
		RunId:   testUUID.String(),
		Service: "TestService1",
		AddedAt: nil,
	}
	expectedDLog1 = models.DLog{
		Id:      1,
		RunId:   testUUID.String(),
		Url:     "TestURL1",
		AddedAt: nil,
	}
	expectedELog1 = models.ELog{
		Id:      1,
		RunId:   testUUID.String(),
		File:    "TestFile1",
		AddedAt: nil,
	}
)

func TestInsertRLog(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("database.TestInsertRLog :: Failed to create SQL mock! -> %v\n", err)
	}

	result := sqlmock.NewResult(1, 1)

	mock.ExpectExec("INSERT OR IGNORE run_log \\(run_id, service\\) VALUES \\(\\?, \\?\\)").
		WithArgs(testUUID.String(), "TestService1").
		WillReturnResult(result)

	ra, err := InsertRLog(db, expectedRLog1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), ra)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("database.TestInsertRLog :: there were unfulfilled expectations: %v\n", err)
	}
}

func TestInsertDLog(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("database.TestInsertDLog :: Failed to create SQL mock! -> %v\n", err)
	}

	result := sqlmock.NewResult(1, 1)

	mock.ExpectExec("INSERT OR IGNORE download_log \\(run_id, url\\) VALUES \\(\\?, \\?\\)").
		WithArgs(testUUID.String(), "TestURL1").
		WillReturnResult(result)

	ra, err := InsertDLog(db, expectedDLog1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), ra)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("database.TestInsertDLog :: there were unfulfilled expectations! -> %v\n", err)
	}
}

func TestInsertELog(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("database.TestInsertELog :: Failed to create SQL mock! -> %v\n", err)
	}

	result := sqlmock.NewResult(1, 1)

	mock.ExpectExec("INSERT OR IGNORE extract_log \\(run_id, file\\) VALUES \\(\\?, \\?\\)").
		WithArgs(testUUID.String(), "TestFile1").
		WillReturnResult(result)

	ra, err := InsertELog(db, expectedELog1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), ra)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("database.TestInsertELog :: there were unfulfilled expectations: %v\n", err)
	}
}
