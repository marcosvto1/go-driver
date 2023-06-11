package files

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	file, err := New(int64(1), "file", "jpg", "/any/path")
	if err != nil {
		t.Error(err)
	}

	file.FolderId = 1

	mock.ExpectExec(`INSERT INTO files (folder_id, owner_id, name, type, path, modified_at)*`).
		WithArgs(
			1,
			1,
			"file",
			"jpg",
			"/any/path",
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = InsertOne(db, file)
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
