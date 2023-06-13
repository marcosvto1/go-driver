package files

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	expectedQuery := regexp.QuoteMeta(`SELECT
		id,
		name,
		folder_id,
		owner_id,
		type,
		path,
		created_at,
		modified_at,
		deleted
		FROM files
	WHERE folder_id = $1 AND deleted = false`)

	fmt.Println(expectedQuery)

	folderID := 1

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
	}).AddRow(
		1,
		"any_name",
		1,
		1,
		"file",
		"/any/path",
		time.Now(),
		time.Now(),
		false,
	)

	mock.ExpectQuery(expectedQuery).
		WithArgs(folderID).
		WillReturnRows(rows)

	files, err := List(db, int64(folderID))
	if err != nil {
		t.Error(err)
	}

	if len(files) == 0 {
		t.Error(errors.New("invalid result"))
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}

func TestListRoot(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	expectedQuery := regexp.QuoteMeta(`SELECT
	id,
	name,
	folder_id,
	owner_id,
	type,
	path,
	created_at,
	modified_at,
	deleted
	FROM files
	WHERE folder_id IS NULL AND deleted = false`)

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
	}).AddRow(
		1,
		"any_name",
		1,
		1,
		"file",
		"/any/path",
		time.Now(),
		time.Now(),
		false,
	)

	mock.ExpectQuery(expectedQuery).
		WillReturnRows(rows)

	files, err := ListRootFiles(db)
	if err != nil {
		t.Error(err)
	}

	if len(files) == 0 {
		t.Error(errors.New("invalid result"))
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
