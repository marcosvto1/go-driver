package users

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestGetById() {
	tcs := []struct {
		Desc             string
		Id               string
		ExpectStatusCode int
		WithMock         bool
		WithMocKError    bool
		Err              error
	}{
		{
			Desc:             "should returns status code 200",
			ExpectStatusCode: http.StatusOK,
			Id:               "1",
			WithMock:         true,
			WithMocKError:    false,
			Err:              nil,
		},
		{
			Desc:             "should return status 500 if id is invalid",
			ExpectStatusCode: http.StatusInternalServerError,
			Id:               "",
			WithMock:         false,
			WithMocKError:    false,
			Err:              nil,
		},
		{
			Desc:             "should return status 404 when not found user",
			ExpectStatusCode: http.StatusNotFound,
			Id:               "1",
			WithMock:         true,
			WithMocKError:    true,
			Err:              sql.ErrNoRows,
		},
		{
			Desc:             "should return status 500 when database failed",
			Id:               "1",
			ExpectStatusCode: http.StatusInternalServerError,
			WithMock:         true,
			WithMocKError:    true,
			Err:              sql.ErrConnDone,
		},
	}

	for _, tc := range tcs {
		ts.T().Log(tc.Desc)

		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.Id)

		rw := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/{id}", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		if tc.WithMock {
			setMockSelect(ts.mock, tc.WithMocKError, tc.Err)
		}

		ts.handler.GetById(rw, req)

		assert.Equal(ts.T(), tc.ExpectStatusCode, rw.Result().StatusCode)
	}
}

func (ts *TransactionSuite) TestGet() {
	setMockSelect(ts.mock, false, nil)

	userResult, err := Get(ts.conn, int64(1))
	assert.NoError(ts.T(), err)

	assert.Equal(ts.T(), "Marcos", userResult.Name)
	assert.Equal(ts.T(), "marcos@email", userResult.Login)
}

func setMockSelect(mock sqlmock.Sqlmock, withError bool, err error) {
	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "modified_at", "deleted", "last_login"}).
		AddRow(1, "Marcos", "marcos@email", "123", time.Now(), time.Now(), false, time.Now())

	expect := mock.ExpectQuery(`SELECT id, name, login, password, created_at, modified_at, deleted, last_login
		FROM users
		WHERE *
	`).WithArgs(int64(1))

	if withError {
		expect.WillReturnError(err)
	} else {
		expect.WillReturnRows(rows)
	}
}
