package users

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestCreate() {
	tcs := []struct {
		Desc             string
		Body             any
		ExpectStatusCode int
		WithMock         bool
		WithMockError    bool
	}{
		{
			Desc:             "should return status 200",
			Body:             ts.entity,
			ExpectStatusCode: http.StatusCreated,
			WithMock:         true,
			WithMockError:    false,
		},
		{
			Desc:             "should return status 500 if failed insert database",
			Body:             ts.entity,
			ExpectStatusCode: http.StatusInternalServerError,
			WithMock:         true,
			WithMockError:    true,
		},
		{
			Desc:             "should return status 500 when invalid body request",
			Body:             "",
			ExpectStatusCode: http.StatusInternalServerError,
			WithMock:         false,
			WithMockError:    false,
		},
		{
			Desc:             "should return status 400 when body is invalid",
			Body:             &User{Name: ""},
			ExpectStatusCode: http.StatusBadRequest,
			WithMock:         false,
			WithMockError:    false,
		},
	}

	for _, tc := range tcs {

		ts.entity = &User{
			Name:     "Marcos",
			Login:    "marcosvto1@gmail.com",
			Password: "123456",
		}

		ts.T().Log(tc.Desc)

		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(&tc.Body)
		assert.NoError(ts.T(), err)

		rr := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodPost, "/users", &b)

		ts.entity.SetPassword(ts.entity.Password)

		if tc.WithMock {
			setMockInsert(ts.mock, ts.entity, tc.WithMockError)
		}

		ts.handler.Create(rr, req)

		assert.Equal(ts.T(), tc.ExpectStatusCode, rr.Result().StatusCode)
	}

}

func (ts *TransactionSuite) TestInsert() {
	setMockInsert(ts.mock, ts.entity, false)

	_, err := Insert(ts.conn, ts.entity)
	assert.NoError(ts.T(), err)
}

func setMockInsert(mock sqlmock.Sqlmock, entity *User, withMockError bool) {
	expect := mock.ExpectExec(`INSERT INTO "users" ("name", "login", "password", "modified_at")*`).
		WithArgs("Marcos", "marcosvto1@gmail.com", entity.Password, entity.ModifiedAt)

	if withMockError {
		expect.WillReturnError(sql.ErrConnDone)
	} else {
		expect.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}
