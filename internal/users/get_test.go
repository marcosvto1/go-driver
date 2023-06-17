package users

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	h := handler{
		db,
	}

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/{id}", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "modified_at", "deleted", "last_login"}).
		AddRow(1, "Marcos", "marcos@email", "123", time.Now(), time.Now(), false, time.Now())

	mock.ExpectQuery(`SELECT id, name, login, password, created_at, modified_at, deleted, last_login
		FROM users
		WHERE *
	`).WithArgs(int64(1)).
		WillReturnRows(rows)

	h.GetById(rw, req)

	assert.Equal(t, http.StatusOK, rw.Result().StatusCode)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "modified_at", "deleted", "last_login"}).
		AddRow(1, "Marcos", "marcos@email", "123", time.Now(), time.Now(), false, time.Now())

	mock.ExpectQuery(`SELECT id, name, login, password, created_at, modified_at, deleted, last_login
		FROM users
		WHERE *
	`).WithArgs(int64(1)).
		WillReturnRows(rows)

	userResult, err := Get(db, int64(1))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "Marcos", userResult.Name)
	assert.Equal(t, "marcos@email", userResult.Login)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}

}
