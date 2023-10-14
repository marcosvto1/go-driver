package folders

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestModifier() {
	defer ts.conn.Close()

	okBodyToUpdate := &Folder{
		Name: "any_name_folder",
	}

	tsc := []struct {
		Desc               string
		Body               any
		ID                 string
		ExpectedStatusCode int
		WithMock           bool
		WithMockErr        bool
		WithMockGet        bool
		WithMockGetErr     bool
	}{
		{
			Desc:               "Should return status 500 when id param is invalid",
			ExpectedStatusCode: http.StatusInternalServerError,
			ID:                 "",
			WithMock:           false,
			WithMockErr:        false,
			Body:               okBodyToUpdate,
		},
		{
			Desc:               "Should return status 500 when request body is invalid",
			ExpectedStatusCode: http.StatusInternalServerError,
			ID:                 "1",
			WithMock:           false,
			WithMockErr:        false,
			Body:               "",
		},
		{
			Desc:               "Should return status 400 when body contains invalid data",
			ExpectedStatusCode: http.StatusBadRequest,
			ID:                 "1",
			WithMock:           false,
			WithMockErr:        false,
			Body: &Folder{
				Name: "",
			},
		},
		{
			Desc:               "Should return status 500 when failed update in database",
			ExpectedStatusCode: http.StatusInternalServerError,
			ID:                 "1",
			WithMock:           true,
			WithMockErr:        true,
			WithMockGetErr:     false,
			WithMockGet:        false,
			Body:               okBodyToUpdate,
		},
		{
			Desc:               "Should return status 500 when failed get folder in database",
			ExpectedStatusCode: http.StatusInternalServerError,
			ID:                 "1",
			WithMock:           true,
			WithMockErr:        false,
			WithMockGet:        true,
			WithMockGetErr:     true,
			Body:               okBodyToUpdate,
		},
	}

	for _, tc := range tsc {
		ts.T().Log(tc.Desc)

		if tc.WithMock {
			setUpdateMock(ts.mock, tc.WithMockErr)
		}

		if tc.WithMockGet {
			setGetMock(ts.mock, tc.WithMockGetErr)
		}

		var b bytes.Buffer
		json.NewEncoder(&b).Encode(tc.Body)
		record := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPut, "/{id}", &b)

		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.ID)

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

		h := handler{
			ts.conn,
		}

		h.Modify(record, request)

		assert.Equal(ts.T(), tc.ExpectedStatusCode, record.Result().StatusCode)
	}
}

func (ts *TransactionSuite) TestUpdate() {
	defer ts.conn.Close()

	setUpdateMock(ts.mock, false)

	err := Update(ts.conn, int64(1), ts.entity)
	assert.NoError(ts.T(), err)
}

func setUpdateMock(mock sqlmock.Sqlmock, withMockErr bool) {
	expectedSQL := regexp.QuoteMeta(`UPDATE folders SET name=$1, modified_at=$2 WHERE id=$3`)

	expect := mock.ExpectExec(expectedSQL).
		WithArgs(
			"any_name_folder",
			sqlmock.AnyArg(),
			1,
		)

	if withMockErr {
		expect.WillReturnError(sqlmock.ErrCancelled)
	} else {
		expect.WillReturnResult(sqlmock.NewResult(1, 1))
	}

}

func setGetMock(mock sqlmock.Sqlmock, withMockErr bool) {
	columns := []string{
		"id",
		"name",
		"parent_id",
		"created_at",
		"modified_at",
		"deleted",
	}

	rows := sqlmock.NewRows(columns).
		AddRow(
			1,
			"any_name",
			1,
			time.Now(),
			time.Now(),
			false,
		)

	expectedSQLSelect := regexp.QuoteMeta(`
		SELECT
			id,
			name,
			parent_id,
			created_at,
			modified_at,
			deleted
		FROM "folders" where id=$1`)

	expect := mock.ExpectQuery(expectedSQLSelect).WithArgs(1)
	if withMockErr {
		expect.WillReturnError(sqlmock.ErrCancelled)
	} else {
		expect.
			WillReturnRows(rows)
	}
}
