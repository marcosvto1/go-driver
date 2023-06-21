package users

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestList() {
	defer ts.conn.Close()
	setMockList(ts.mock)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	ts.handler.List(recorder, request)

	assert.Equal(ts.T(), http.StatusOK, recorder.Result().StatusCode)
}

func (ts *TransactionSuite) TestFindAll() {
	defer ts.conn.Close()
	setMockList(ts.mock)

	users, err := FindAll(ts.conn)

	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), 1, len(users))
}

func setMockList(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "modified_at", "deleted", "last_login"}).
		AddRow(1, "Marcos", "marcos@email", "123", time.Now(), time.Now(), false, time.Now())

	mock.ExpectQuery(`SELECT id, name, login, password, created_at, modified_at, deleted, last_login FROM users
WHERE *`).WillReturnRows(rows)
}
