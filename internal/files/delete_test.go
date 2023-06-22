package files

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestDeleteHTTP() {
	defer ts.conn.Close()

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/{id}", nil)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

	setMockDelete(ts.mock)

	ts.handler.Delete(recorder, request)

	assert.Equal(ts.T(), http.StatusNoContent, recorder.Result().StatusCode)
}

func (ts *TransactionSuite) TestDelete() {
	defer ts.conn.Close()

	setMockDelete(ts.mock)

	err := Delete(ts.conn, int64(1))

	assert.NoError(ts.T(), err)
}

func setMockDelete(mock sqlmock.Sqlmock) {
	mock.ExpectExec(`UPDATE files SET *`).
		WithArgs(true, sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
}
