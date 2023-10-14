package folders

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestDeleteHTTP() {
	defer ts.conn.Close()

	tcs := []struct {
		Desc                       string
		ExpectedStatusCode         int
		Id                         string
		WithMockListFile           bool
		WithMockListFileErr        bool
		WithMockUpdateDelFile      bool
		WithMockSelectSubFolder    bool
		WithMockSelectSubFolderErr bool
		WithMockUpdateDelFolder    bool
		WithMockUpdateDelFolderErr bool
	}{
		{
			Desc:                       "Should return status 500 when params is invalid",
			ExpectedStatusCode:         500,
			Id:                         "",
			WithMockListFile:           false,
			WithMockUpdateDelFile:      false,
			WithMockSelectSubFolder:    false,
			WithMockSelectSubFolderErr: false,
		},
		{
			Desc:                       "Should return 500 when failed delete folder",
			ExpectedStatusCode:         500,
			Id:                         "1",
			WithMockListFile:           true,
			WithMockListFileErr:        true,
			WithMockUpdateDelFile:      false,
			WithMockSelectSubFolder:    false,
			WithMockSelectSubFolderErr: false,
		},
		{
			Desc:                       "Should return 500 if failed to get sub folders",
			ExpectedStatusCode:         500,
			Id:                         "1",
			WithMockListFile:           true,
			WithMockListFileErr:        false,
			WithMockUpdateDelFile:      true,
			WithMockSelectSubFolder:    true,
			WithMockSelectSubFolderErr: true,
		},
		{
			Desc:                       "Should return 500 if failed delete current folder",
			ExpectedStatusCode:         500,
			Id:                         "1",
			WithMockListFile:           true,
			WithMockListFileErr:        false,
			WithMockUpdateDelFile:      true,
			WithMockSelectSubFolder:    true,
			WithMockSelectSubFolderErr: false,
			WithMockUpdateDelFolder:    true,
			WithMockUpdateDelFolderErr: true,
		},

		{
			Desc:                       "Should return 200 when delete with success",
			ExpectedStatusCode:         200,
			Id:                         "1",
			WithMockListFile:           true,
			WithMockListFileErr:        false,
			WithMockUpdateDelFile:      true,
			WithMockSelectSubFolder:    true,
			WithMockSelectSubFolderErr: false,
			WithMockUpdateDelFolder:    true,
			WithMockUpdateDelFolderErr: false,
		},
	}

	for _, tc := range tcs {
		ts.T().Log(tc.Desc)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodDelete, "/{id}", nil)

		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.Id)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

		if tc.WithMockListFile {
			setSelectFilesForDeleteMock(ts.mock, tc.WithMockListFileErr)
		}

		if tc.WithMockUpdateDelFile {
			setMockUpdateDeleteFile(ts.mock, false)
		}

		if tc.WithMockSelectSubFolder {
			setMockSubFolders(ts.mock, tc.WithMockSelectSubFolderErr)
		}

		if tc.WithMockUpdateDelFolder {
			expectedQueryDeleteFolder := regexp.QuoteMeta(`UPDATE "folders" SET "modified_at"=$1, "deleted"=true WHERE id=$2`)

			expect := ts.mock.ExpectExec(expectedQueryDeleteFolder).
				WithArgs(sqlmock.AnyArg(), 1).
				WillReturnResult(sqlmock.NewResult(1, 1))

			if tc.WithMockUpdateDelFolderErr {
				expect.WillReturnError(sqlmock.ErrCancelled)
			} else {
				expect.
					WillReturnResult(sqlmock.NewResult(1, 1))
			}
		}

		ts.handler.Delete(recorder, request)

		assert.Equal(ts.T(), tc.ExpectedStatusCode, recorder.Result().StatusCode)
	}

}

func (ts *TransactionSuite) TestDelete() {
	defer ts.conn.Close()

	expectedQuery := regexp.QuoteMeta(`UPDATE "folders" SET "modified_at"=$1, "deleted"=true WHERE id=$2`)

	ts.mock.ExpectExec(expectedQuery).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := Delete(ts.conn, int64(1))
	assert.NoError(ts.T(), err)
}

func setSelectFilesForDeleteMock(mock sqlmock.Sqlmock, withMockErr bool) {
	expectedQuery := regexp.QuoteMeta(`SELECT
	id,
	name,
	folder_id,
	owner_id,
	type,
	path,
	created_at,
	modified_at,
	deleted
	FROM files
WHERE folder_id = $1 AND deleted = false`)

	folderID := 1

	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"folder_id",
		"owner_id",
		"type",
		"path",
		"created_at",
		"modified_at",
		"deleted_at",
	}).AddRow(
		1,
		"any_name",
		1,
		1,
		"file",
		"/any/path",
		time.Now(),
		time.Now(),
		false,
	)

	expect := mock.ExpectQuery(expectedQuery).
		WithArgs(folderID)

	if withMockErr {
		expect.WillReturnError(sqlmock.ErrCancelled)
	} else {
		expect.WillReturnRows(rows)
	}

}

func setMockUpdateDeleteFile(mock sqlmock.Sqlmock, withMockErr bool) {
	expect := mock.ExpectExec(regexp.QuoteMeta(`UPDATE "files" SET "name"=$1, "modified_at"=$2, "deleted"=$3 where id = $4`)).
		WithArgs("any_name", sqlmock.AnyArg(), true, 1)

	if withMockErr {
		fmt.Println("eere")
		expect.WillReturnError(sqlmock.ErrCancelled)
	} else {
		expect.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func setMockSubFolders(mock sqlmock.Sqlmock, withMockErr bool) {
	columns := []string{
		"id",
		"name",
		"parent_id",
		"created_at",
		"modified_at",
		"deleted",
	}

	rootrows := sqlmock.NewRows(columns).
		AddRow(
			2,
			"any folder name",
			2,
			time.Now(),
			time.Now(),
			false,
		).
		AddRow(
			2,
			"any folder name 2",
			3,
			time.Now(),
			time.Now(),
			false,
		)

	expectedSQL := regexp.QuoteMeta(`
		SELECT
			id,
			name,
			parent_id,
			created_at,
			modified_at,
			deleted
		FROM "folders" where parent_id=$1`,
	)

	expect := mock.ExpectQuery(expectedSQL).
		WithArgs(1)

	if withMockErr {
		expect.WillReturnError(sqlmock.ErrCancelled)
	} else {
		expect.WillReturnRows(rootrows)
	}

}
