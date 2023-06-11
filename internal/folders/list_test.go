package folders

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRootSubFolder(t *testing.T) {
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
	FROM "folders" WHERE "parent_id" IS NULL (.+)`

	mock.ExpectQuery(expectedSQL).
		WillReturnRows(rows)

	folders, err := getRootSubFolder(db)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, len(folders), 1)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
