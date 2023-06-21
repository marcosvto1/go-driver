package users

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestGetById() {
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/{id}", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	setMockSelect(ts.mock)

	ts.handler.GetById(rw, req)

	assert.Equal(ts.T(), http.StatusOK, rw.Result().StatusCode)
}

func (ts *TransactionSuite) TestGet() {
	setMockSelect(ts.mock)

	userResult, err := Get(ts.conn, int64(1))
	assert.NoError(ts.T(), err)

	assert.Equal(ts.T(), "Marcos", userResult.Name)
	assert.Equal(ts.T(), "marcos@email", userResult.Login)
}

func setMockSelect(mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "modified_at", "deleted", "last_login"}).
		AddRow(1, "Marcos", "marcos@email", "123", time.Now(), time.Now(), false, time.Now())

	mock.ExpectQuery(`SELECT id, name, login, password, created_at, modified_at, deleted, last_login
		FROM users
		WHERE *
	`).WithArgs(int64(1)).
		WillReturnRows(rows)
}
