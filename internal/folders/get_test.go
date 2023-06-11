package folders

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetFolder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	columns := []string{
		"id",
		"name",
		"parent_id",
		"created_at",
		"modified_at",
		"deleted",
	}

	rows := sqlmock.NewRows(columns).
		AddRow(
			1,
			"any folder name",
			1,
			time.Now(),
			time.Now(),
			false,
		)

	expectedSQL := `SELECT
	id,
	name,
	parent_id,
	created_at,
	modified_at,
	deleted
	FROM "folders" *`

	mock.ExpectQuery(expectedSQL).
		WithArgs(1).
		WillReturnRows(rows)

	folder, err := GetFolder(db, 1)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, folder.Name, "any folder name")

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}

func TestGetSubFolder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	columns := []string{
		"id",
		"name",
		"parent_id",
		"created_at",
		"modified_at",
		"deleted",
	}

	rows := sqlmock.NewRows(columns).
		AddRow(
			1,
			"any folder name",
			1,
			time.Now(),
			time.Now(),
			false,
		)

	expectedSQL := `SELECT
	id,
	name,
	parent_id,
	created_at,
	modified_at,
	deleted
	FROM "folders" where parent_id *`

	mock.ExpectQuery(expectedSQL).
		WithArgs(1).
		WillReturnRows(rows)

	folders, err := getSubFolder(db, 1)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, len(folders), 1)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
