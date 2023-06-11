package folders

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	folder, err := New("any_name_folder", 1)
	if err != nil {
		t.Error(err)
	}

	expectedQuery := `INSERT INTO "folders" ("name", "parent_id", "modified_at")*`
	mock.ExpectExec(expectedQuery).
		WithArgs(
			folder.Name,
			folder.ParentID,
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = Insert(db, folder)
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
