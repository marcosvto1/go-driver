package users

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

	u, err := New("Marcos", "marcosvto1@gmail.com", "1234567")
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec(`UPDATE "users" SET *`).
		WithArgs(u.Name, u.Login, u.Password, AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = Update(db, int64(1), u)
	if err != nil {
		t.Error(err)
	}
}
