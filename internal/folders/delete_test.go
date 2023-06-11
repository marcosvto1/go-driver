package folders

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	expectedQuery := `UPDATE "folders" SET *`
	mock.ExpectExec(expectedQuery).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = Delete(db, int64(1))
	if err != nil {
		t.Error(t)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(t)
	}
}
