package folders

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(t)
	}
	defer db.Close()

	folder, err := New("any_name_folder", 1)
	if err != nil {
		t.Error(err)
	}

	expectedSQL := `UPDATE folders SET *`
	mock.ExpectExec(expectedSQL).
		WithArgs(
			folder.Name,
			sqlmock.AnyArg(),
			1,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = Update(db, int64(1), folder)
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
