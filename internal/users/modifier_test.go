package users

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestModify() {
	tcs := []struct {
		Desc                string
		Id                  string
		Body                any
		ExpectedStatusCode  int
		WithMock            bool
		WithMocKError       bool
		WithMockSelect      bool
		WithMockSelectError bool
		SelectErr           error
	}{
		{
			Desc:               "Should return 500 if id not provided",
			ExpectedStatusCode: http.StatusInternalServerError,
			WithMock:           false,
			WithMocKError:      false,
			Id:                 "",
			Body: User{
				Name: "Marcos",
			},
		},
		{
			Desc:               "Should return 500 when body is invalid",
			ExpectedStatusCode: http.StatusInternalServerError,
			WithMock:           false,
			WithMocKError:      false,
			Id:                 "1",
			Body:               "",
		},

		{
			Desc:               "Should return 400 when name is empty",
			ExpectedStatusCode: http.StatusBadRequest,
			WithMock:           false,
			WithMocKError:      false,
			Id:                 "1",
			Body: User{
				Name: "",
			},
		},

		{
			Desc:                fmt.Sprintf("should return %v when successfully", http.StatusOK),
			ExpectedStatusCode:  http.StatusOK,
			WithMock:            true,
			WithMocKError:       false,
			WithMockSelect:      true,
			WithMockSelectError: false,
			SelectErr:           nil,
			Id:                  "1",
			Body: User{
				Name: "Marcos",
			},
		},

		{
			Desc:                fmt.Sprintf("should return %v when failed update", http.StatusInternalServerError),
			ExpectedStatusCode:  http.StatusInternalServerError,
			WithMock:            true,
			WithMocKError:       true,
			WithMockSelect:      false,
			WithMockSelectError: false,
			Id:                  "1",
			Body: User{
				Name: "Marcos",
			},
		},

		{
			Desc:                fmt.Sprintf("should return %v if find user return error", http.StatusInternalServerError),
			ExpectedStatusCode:  http.StatusInternalServerError,
			WithMock:            true,
			WithMocKError:       false,
			WithMockSelect:      true,
			WithMockSelectError: true,
			SelectErr:           sql.ErrConnDone,
			Id:                  "1",
			Body: User{
				Name: "Marcos",
			},
		},
	}

	defer ts.conn.Close()

	for _, tc := range tcs {
		ts.T().Log(tc.Desc)

		u := tc.Body

		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(&u)
		assert.NoError(ts.T(), err)

		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.Id)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPut, "/{id}", &b)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

		if tc.WithMock {
			id, _ := strconv.Atoi(tc.Id)
			setMockModifier(ts.mock, ts.entity, tc.WithMocKError)

			if tc.WithMockSelect {
				setMockSelect(ts.mock, id, tc.WithMockSelectError, tc.SelectErr)
			}
		}

		ts.handler.Modify(recorder, request)

		assert.Equal(ts.T(), tc.ExpectedStatusCode, recorder.Result().StatusCode)
	}

}

func (ts *TransactionSuite) TestUpdate() {
	defer ts.conn.Close()

	setMockModifier(ts.mock, ts.entity, false)

	err := Update(ts.conn, int64(1), ts.entity)
	assert.NoError(ts.T(), err)
}

func setMockModifier(mock sqlmock.Sqlmock, entity *User, withMockErr bool) {
	expectedSQL := regexp.QuoteMeta(`UPDATE "users" SET "name"=$1, "modified_at"=$2 WHERE id=$3`)
	expect := mock.ExpectExec(expectedSQL).
		WithArgs(entity.Name, AnyTime{}, 1)
	if withMockErr {
		expect.WillReturnError(sql.ErrConnDone)
	} else {
		expect.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}
