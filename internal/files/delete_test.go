package files

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestDeleteHTTP() {
	defer ts.conn.Close()

	tcs := []struct {
		Desc               string
		Id                 string
		ExpectedStatusCode int
		WithMock           bool
		WithMockErr        bool
	}{
		{Desc: "Should returns status Not Content", Id: "1", ExpectedStatusCode: http.StatusNoContent, WithMock: true, WithMockErr: false},
		{Desc: "Should return internal server error when invalid id", Id: "", ExpectedStatusCode: http.StatusInternalServerError, WithMock: false, WithMockErr: false},
		{Desc: "Should return internal server error when failed to delete file", Id: "1", ExpectedStatusCode: http.StatusInternalServerError, WithMock: true, WithMockErr: true},
	}

	for _, tc := range tcs {
		ts.T().Log(tc.Desc)

		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.Id)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodDelete, "/{id}", nil)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

		if tc.WithMock {
			setMockDelete(ts.mock, tc.WithMockErr)
		}

		ts.handler.Delete(recorder, request)

		assert.Equal(ts.T(), tc.ExpectedStatusCode, recorder.Result().StatusCode)
	}
}

func (ts *TransactionSuite) TestDelete() {
	defer ts.conn.Close()

	setMockDelete(ts.mock, false)
	err := Delete(ts.conn, int64(1))
	assert.NoError(ts.T(), err)
}

func setMockDelete(mock sqlmock.Sqlmock, withMockErr bool) {
	expect := mock.ExpectExec(`UPDATE files SET *`).
		WithArgs(true, sqlmock.AnyArg(), 1)

	if withMockErr {
		expect.WillReturnError(sql.ErrNoRows)
	} else {
		expect.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}
