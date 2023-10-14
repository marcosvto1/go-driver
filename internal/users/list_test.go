package users

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestList() {
	tcs := []struct {
		Desc               string
		ExpectedStatusCode int
		WithMock           bool
		WithMockErr        bool
	}{
		{
			Desc:               fmt.Sprintf("should return %d if successfully", http.StatusOK),
			ExpectedStatusCode: http.StatusOK,
			WithMock:           true,
			WithMockErr:        false,
		},
		{
			Desc:               fmt.Sprintf("should return %d if find users", http.StatusOK),
			ExpectedStatusCode: http.StatusInternalServerError,
			WithMock:           true,
			WithMockErr:        true,
		},
	}

	defer ts.conn.Close()

	for _, tc := range tcs {
		ts.T().Log(tc.Desc)

		if tc.WithMock {
			setMockList(ts.mock, tc.WithMockErr)
		}

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/", nil)

		ts.handler.List(recorder, request)

		assert.Equal(ts.T(), tc.ExpectedStatusCode, recorder.Result().StatusCode)
	}

}

func (ts *TransactionSuite) TestFindAll() {
	defer ts.conn.Close()
	setMockList(ts.mock, false)

	users, err := FindAll(ts.conn)

	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), 1, len(users))
}

func setMockList(mock sqlmock.Sqlmock, withMockErr bool) {
	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "modified_at", "deleted", "last_login"}).
		AddRow(1, "Marcos", "marcos@email", "123", time.Now(), time.Now(), false, time.Now())

	expect := mock.ExpectQuery(`SELECT id, name, login, password, created_at, modified_at, deleted, last_login FROM users
WHERE *`)

	if withMockErr {
		expect.WillReturnError(sql.ErrConnDone)
	} else {
		expect.WillReturnRows(rows)
	}

}
