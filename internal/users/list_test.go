package users

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "modified_at", "deleted", "last_login"}).
		AddRow(1, "Marcos", "marcos@email", "123", time.Now(), time.Now(), false, time.Now())

	mock.ExpectQuery(`SELECT id, name, login, password, created_at, modified_at, deleted, last_login FROM users
WHERE *`).WillReturnRows(rows)

	h := handler{
		db,
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	h.List(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "modified_at", "deleted", "last_login"}).
		AddRow(1, "Marcos", "marcos@email", "123", time.Now(), time.Now(), false, time.Now())

	mock.ExpectQuery(`SELECT id, name, login, password, created_at, modified_at, deleted, last_login FROM users
	WHERE *`).WillReturnRows(rows)

	users, err := FindAll(db)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 1, len(users))

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}

}
