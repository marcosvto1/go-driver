package users

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestModify(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	u := User{
		Name: "Marcos",
	}
	var b bytes.Buffer

	err = json.NewEncoder(&b).Encode(&u)
	if err != nil {
		t.Error(err)
	}

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")

	wr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/{id}", &b)

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

	h := handler{
		db,
	}

	expectedSQL := regexp.QuoteMeta(`UPDATE "users" SET "name"=$1, "modified_at"=$4 WHERE id=$5`)
	mock.ExpectExec(expectedSQL).
		WithArgs(u.Name, sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	expectedSelectSQL := regexp.QuoteMeta(`
		SELECT id, name, login, password, created_at, modified_at, deleted, last_login
		FROM users
		WHERE id = $1"
	`)
	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "created_at", "modified_at", "deleted", "last_login"}).
		AddRow(1, "Marcos", "marcos@email", "123", time.Now(), time.Now(), false, time.Now())
	mock.ExpectQuery(expectedSelectSQL).WithArgs(int64(1)).
		WillReturnRows(rows)

	h.Modify(wr, r)

	assert.Equal(t, http.StatusOK, wr.Result().StatusCode)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}

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

	expectedSQL := regexp.QuoteMeta(`UPDATE "users" SET "name"=$1, "modified_at"=$4 WHERE id=$5`)

	mock.ExpectExec(expectedSQL).
		WithArgs(u.Name, AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = Update(db, int64(1), u)
	if err != nil {
		t.Error(err)
	}
}
