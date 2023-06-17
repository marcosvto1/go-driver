package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	h := handler{
		db,
	}

	u := User{
		Name:     "Marcos",
		Login:    "marcosvto1@gmail.com",
		Password: "123456",
	}

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(&u)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/users", &b)

	u.SetPassword(u.Password)

	mock.ExpectExec(`INSERT INTO "users" ("name", "login", "password", "modified_at")*`).
		WithArgs("Marcos", "marcosvto1@gmail.com", u.Password, u.ModifiedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	h.Create(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Result().StatusCode)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	u, err := New("Marcos", "marcosvto1@gmail.com", "123456")
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec(`INSERT INTO "users" ("name", "login", "password", "modified_at")*`).
		WithArgs("Marcos", "marcosvto1@gmail.com", u.Password, u.ModifiedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = Insert(db, u)
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
