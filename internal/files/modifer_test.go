package files

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	file, err := New(1, "Name", "ext", "/your/path")
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec(`UPDATE "files" SET`).
		WithArgs(file.Name, sqlmock.AnyArg(), file.Deleted, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = Update(db, int64(1), file)
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
