package folders

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestList() {
	defer ts.conn.Close()

	tcs := []struct {
		Desc               string
		ExpectedStatusCode int
		WithMockErr        bool
	}{
		{
			Desc:               "Should returns status OK",
			ExpectedStatusCode: http.StatusOK,
			WithMockErr:        false,
		},
		{
			Desc:               "Should returns status 500",
			ExpectedStatusCode: http.StatusInternalServerError,
			WithMockErr:        true,
		},
	}

	for _, tc := range tcs {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/", nil)

		setMockGetSubFolder(ts.mock)
		setMockListRootFiles(ts.mock, tc.WithMockErr)

		h := handler{
			ts.conn,
		}

		h.List(recorder, request)

		assert.Equal(ts.T(), tc.ExpectedStatusCode, recorder.Result().StatusCode)
	}
}

func (ts *TransactionSuite) TestRootSubFolder() {
	defer ts.conn.Close()

	setMockGetSubFolder(ts.mock)

	folders, err := getRootSubFolder(ts.conn)

	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), len(folders), 1)
}

func setMockGetSubFolder(mock sqlmock.Sqlmock) {
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
			0,
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
	FROM "folders" WHERE "parent_id" IS NULL "deleted"=false`)

	mock.ExpectQuery(expectedSQL).
		WillReturnRows(rows)
}

func setMockListRootFiles(mock sqlmock.Sqlmock, withMockErr bool) {
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
	WHERE folder_id IS NULL AND deleted = false`)

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

	expected := mock.ExpectQuery(expectedSQL)

	if !withMockErr {
		expected.WillReturnRows(rows)
	} else {
		expected.WillReturnError(sql.ErrNoRows)
	}
}
