package users

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestDeleteOne() {
	// define test case
	tcs := []struct {
		ID               string
		ExpectStatusCode int
		WithMock         bool
		WithMockError    bool
	}{
		{ID: "1", ExpectStatusCode: http.StatusNoContent, WithMock: true, WithMockError: false},
		{ID: "A", ExpectStatusCode: http.StatusInternalServerError, WithMock: false, WithMockError: false},
		{ID: "1", ExpectStatusCode: http.StatusInternalServerError, WithMock: true, WithMockError: true},
	}

	for _, tc := range tcs {
		// Create a new context with a dummy URL parameter
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.ID)

		// Create a new response recorder
		recorder := httptest.NewRecorder()

		// Create a new request
		request := httptest.NewRequest("DELETE", "/{id}", nil)

		// Set context for request
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

		if tc.WithMock {
			setMockDelete(ts.mock, tc.WithMockError)
		}

		// Call the Delete method of the handler
		ts.handler.Delete(recorder, request)

		// Check the response status code
		assert.Equal(ts.T(), tc.ExpectStatusCode, recorder.Result().StatusCode)
	}
}

func (ts *TransactionSuite) TestDelete() {
	setMockDelete(ts.mock, false)

	err := Delete(ts.conn, int64(1))
	assert.NoError(ts.T(), err)
}

func setMockDelete(mock sqlmock.Sqlmock, withError bool) {
	exp := mock.ExpectExec(`UPDATE users SET *`).
		WithArgs(AnyTime{}, 1)

	if withError {
		exp.WillReturnError(sql.ErrNoRows)
	} else {
		exp.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}
