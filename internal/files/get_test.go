package files

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"folder_id",
		"owner_id",
		"type",
		"path",
		"created_at",
		"modified_at",
		"deleted_at",
	}).
		AddRow(1, "any_name", 1, 1, "file", "/any/path", time.Now(), time.Now(), false)

	mock.ExpectQuery(
		`SELECT
		id,
		name,
		folder_id,
		owner_id,
		type,
		path,
		created_at,
		modified_at,
		deleted
	FROM "files" * `,
	).WithArgs(1).WillReturnRows(rows)

	fr, err := Get(db, int64(1))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, fr.Name, "any_name")
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
