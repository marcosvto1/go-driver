package folders

import (
	"regexp"
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

	expectedSQL := regexp.QuoteMeta(`
	SELECT
		id,
		name,
		parent_id,
		created_at,
		modified_at,
		deleted
	FROM "folders" where id=$1`)

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
			2,
			"any folder name",
			2,
			time.Now(),
			time.Now(),
			false,
		).
		AddRow(
			2,
			"any folder name 2",
			3,
			time.Now(),
			time.Now(),
			false,
		)

	expectedSQL := regexp.QuoteMeta(`
	SELECT
		id,
		name,
		parent_id,
		created_at,
		modified_at,
		deleted
	FROM "folders" where parent_id=$1`,
	)

	mock.ExpectQuery(expectedSQL).
		WithArgs(1).
		WillReturnRows(rows)

	folders, err := getSubFolder(db, 1)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, len(folders), 2)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
