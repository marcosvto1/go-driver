package users

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestModify() {
	defer ts.conn.Close()

	u := User{
		Name: "Marcos",
	}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(&u)
	assert.NoError(ts.T(), err)

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/{id}", &b)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

	setMockModifier(ts.mock, ts.entity)
	setMockSelect(ts.mock, false, nil)

	ts.handler.Modify(recorder, request)

	assert.Equal(ts.T(), http.StatusOK, recorder.Result().StatusCode)
}

func (ts *TransactionSuite) TestUpdate() {
	defer ts.conn.Close()

	setMockModifier(ts.mock, ts.entity)

	err := Update(ts.conn, int64(1), ts.entity)
	assert.NoError(ts.T(), err)
}

func setMockModifier(mock sqlmock.Sqlmock, entity *User) {
	expectedSQL := regexp.QuoteMeta(`UPDATE "users" SET "name"=$1, "modified_at"=$4 WHERE id=$5`)
	mock.ExpectExec(expectedSQL).
		WithArgs(entity.Name, AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
}
