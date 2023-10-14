package folders

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestGet() {
	defer ts.conn.Close()

	tcs := []struct {
		Desc                     string
		Id                       string
		ExpectStatusCode         int
		WithMockGetFolder        bool
		WithMockGetFolderErr     bool
		WithMockFolderContent    bool
		WithMockFolderContentErr bool
	}{
		{
			Desc:                     "Should returns status 500 if params invalid",
			ExpectStatusCode:         http.StatusInternalServerError,
			Id:                       "",
			WithMockGetFolder:        false,
			WithMockGetFolderErr:     false,
			WithMockFolderContent:    false,
			WithMockFolderContentErr: false,
		},
		{
			Desc:                 "Should retuns status 500 when get folder in database failed",
			ExpectStatusCode:     http.StatusInternalServerError,
			Id:                   "1",
			WithMockGetFolder:    true,
			WithMockGetFolderErr: true,
		},
		{
			Desc:                     "Should retuns status 500 when get folder content in database failed",
			ExpectStatusCode:         http.StatusInternalServerError,
			Id:                       "1",
			WithMockGetFolder:        true,
			WithMockGetFolderErr:     false,
			WithMockFolderContent:    true,
			WithMockFolderContentErr: true,
		},
		{
			Desc:                     "Should retuns status 200 when get folder content in database failed",
			ExpectStatusCode:         http.StatusOK,
			Id:                       "1",
			WithMockGetFolder:        true,
			WithMockGetFolderErr:     false,
			WithMockFolderContent:    true,
			WithMockFolderContentErr: false,
		},
	}

	for _, tc := range tcs {
		ts.T().Log(tc.Desc)

		if tc.WithMockGetFolder {
			setGetFolderMock(ts.mock, tc.WithMockGetFolderErr)
		}

		if tc.WithMockFolderContent {
			setGetSubFolderMock(ts.mock, tc.WithMockFolderContentErr)
			if !tc.WithMockFolderContentErr {
				setFileListMock(ts.mock, false)
			}
		}

		record := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/{id}", nil)

		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", tc.Id)

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

		h := handler{
			ts.conn,
		}

		h.Get(record, request)

		assert.Equal(ts.T(), tc.ExpectStatusCode, record.Result().StatusCode)
	}

}

func (ts *TransactionSuite) TestGetFolder() {
	defer ts.conn.Close()

	setGetFolderMock(ts.mock, false)

	folder, err := GetFolder(ts.conn, 1)

	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), folder.Name, "any folder name")
}

func (ts *TransactionSuite) TestGetSubFolder() {
	defer ts.conn.Close()

	setGetSubFolderMock(ts.mock, false)

	folders, err := getSubFolder(ts.conn, 1)

	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), len(folders), 2)
}

func setGetFolderMock(mock sqlmock.Sqlmock, withMockErr bool) {
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
			"any folder name",
			1,
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
	FROM "folders" where id=$1`)

	expected := mock.ExpectQuery(expectedSQL).
		WithArgs(1)

	if !withMockErr {
		expected.WillReturnRows(rows)
	} else {
		expected.WillReturnError(sqlmock.ErrCancelled)
	}
}

func setGetSubFolderMock(mock sqlmock.Sqlmock, withMockErr bool) {
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

	expected := mock.ExpectQuery(expectedSQL).
		WithArgs(1)

	if !withMockErr {
		expected.WillReturnRows(rows)

	} else {
		expected.WillReturnError(sqlmock.ErrCancelled)
	}

}

func setFileListMock(mock sqlmock.Sqlmock, withMockErr bool) {
	expectedSQL := regexp.QuoteMeta(`SELECT
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

	expected := mock.ExpectQuery(expectedSQL).WithArgs(1)

	expected.WillReturnRows(rows)
}
