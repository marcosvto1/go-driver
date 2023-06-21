package users

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestDeleteOne() {

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")

	// Create a new request with a dummy URL parameter
	ts.mock.ExpectExec(`UPDATE users SET *`).
		WithArgs(AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create a new response recorder
	recorder := httptest.NewRecorder()

	// Create a new request
	request := httptest.NewRequest("DELETE", "/delete/123", nil)

	// Set context for request
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

	// Call the Delete method of the handler
	ts.handler.Delete(recorder, request)

	// Check the response status code
	assert.Equal(ts.T(), http.StatusNoContent, recorder.Result().StatusCode)
}

func (ts *TransactionSuite) TestDelete() {
	ts.mock.ExpectExec(`UPDATE users SET *`).
		WithArgs(AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := Delete(ts.conn, int64(1))
	assert.NoError(ts.T(), err)
}
