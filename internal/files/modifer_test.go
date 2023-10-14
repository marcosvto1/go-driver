package files

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestModify() {
	defer ts.conn.Close()

	tcs := []struct {
		Desc               string
		Id                 string
		WithMockUpdate     bool
		WithMockUpdateErr  bool
		WithMockGet        bool
		WithMockGetErr     bool
		ExpectedStatusCode int
	}{
		{
			Desc:               "Should update file and return status no content",
			Id:                 "1",
			WithMockUpdate:     true,
			WithMockUpdateErr:  false,
			WithMockGet:        true,
			WithMockGetErr:     false,
			ExpectedStatusCode: http.StatusNoContent,
		},
		{
			Desc:               "Should return 500 when id is invalid",
			Id:                 "",
			WithMockUpdate:     false,
			WithMockUpdateErr:  false,
			WithMockGet:        false,
			WithMockGetErr:     false,
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		{
			Desc:               "Should return 400 when not found file to update",
			Id:                 "1",
			WithMockUpdate:     false,
			WithMockUpdateErr:  false,
			WithMockGet:        true,
			WithMockGetErr:     true,
			ExpectedStatusCode: http.StatusNotFound,
		},

		{
			Desc:               "Should return internal server error when failed update",
			Id:                 "1",
			WithMockUpdate:     true,
			WithMockUpdateErr:  true,
			WithMockGet:        true,
			WithMockGetErr:     false,
			ExpectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tcs {
		ts.T().Log(tc.Desc)

		// ts.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "files" SET "name"=$1, "modified_at"=$2, "deleted"=$3 where id = $4`)).
		// 	WithArgs("any_name2", sqlmock.AnyArg(), false, 1).
		// 	WillReturnResult(sqlmock.NewResult(1, 1))

		if tc.WithMockGet {
			setMockGet(ts.mock, tc.WithMockGetErr)
		}

		if tc.WithMockUpdate {
			setUpdateMock(ts.mock, "any_name2", tc.WithMockUpdateErr)
		}

		var b bytes.Buffer
		f := File{
			Name:    "any_name2",
			OwnerId: 1,
			ID:      1,
			Type:    "file",
		}

		err := json.NewEncoder(&b).Encode(f)
		assert.NoError(ts.T(), err)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPut, "/{id}", &b)

		// SET CONTEXT
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.Id)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

		ts.handler.Modify(recorder, request)

		fmt.Printf(recorder.Body.String())
		assert.Equal(ts.T(), tc.ExpectedStatusCode, recorder.Result().StatusCode)
	}

}

func (ts *TransactionSuite) TestUpdate() {
	defer ts.conn.Close()

	setUpdateMock(ts.mock, ts.entity.Name, false)

	err := Update(ts.conn, int64(1), ts.entity)
	assert.NoError(ts.T(), err)
}

func setUpdateMock(m sqlmock.Sqlmock, name string, withMockErr bool) {
	expect := m.ExpectExec(regexp.QuoteMeta(`UPDATE "files" SET "name"=$1, "modified_at"=$2, "deleted"=$3 where id = $4`)).
		WithArgs(name, sqlmock.AnyArg(), false, 1)

	if withMockErr {
		expect.WillReturnError(sql.ErrConnDone)
	} else {
		expect.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func setMockGet(m sqlmock.Sqlmock, withMockErr bool) {
	expectSQL := regexp.QuoteMeta(`SELECT
	id,
	name,
	folder_id,
	owner_id,
	type,
	path,
	created_at,
	modified_at,
	deleted
	FROM "files" WHERE "id"=$1`)

	if withMockErr {
		rows := sqlmock.NewRows([]string{
			"id",
		})

		rows.AddRow(nil).RowError(0, fmt.Errorf("row error"))
		expect := m.ExpectQuery(expectSQL).WithArgs(1)
		expect.WillReturnRows(rows)
	} else {
		rows := sqlmock.NewRows([]string{
			"id",
			"name",
			"folder_id",
			"owner_id",
			"type",
			"path",
			"created_at",
			"modified_at",
			"deleted",
		})

		rows.AddRow(1, "any_name", 1, 1, "file", "/any/path", time.Now(), time.Now(), false)

		expect := m.ExpectQuery(expectSQL).WithArgs(1)
		expect.WillReturnRows(rows)
	}
}
