package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestCreate() {
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(&ts.entity)
	assert.NoError(ts.T(), err)

	rr := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/users", &b)

	ts.entity.SetPassword(ts.entity.Password)

	setMockInsert(ts.mock, ts.entity)

	ts.handler.Create(rr, req)

	assert.Equal(ts.T(), http.StatusCreated, rr.Result().StatusCode)
}

func (ts *TransactionSuite) TestInsert() {
	setMockInsert(ts.mock, ts.entity)

	_, err := Insert(ts.conn, ts.entity)
	assert.NoError(ts.T(), err)
}

func setMockInsert(mock sqlmock.Sqlmock, entity *User) {
	mock.ExpectExec(`INSERT INTO "users" ("name", "login", "password", "modified_at")*`).
		WithArgs("Marcos", "marcosvto1@gmail.com", entity.Password, entity.ModifiedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
}
