package files

import (
	"database/sql"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestList() {
	defer ts.conn.Close()

	tcs := []struct {
		Desc        string
		WithMockErr bool
	}{
		{
			Desc:        "Should returns files from filters",
			WithMockErr: false,
		},
		{
			Desc:        "Should return error if database failed",
			WithMockErr: true,
		},
	}

	for _, tc := range tcs {
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

		expect := ts.mock.ExpectQuery(expectedQuery).
			WithArgs(folderID)

		if tc.WithMockErr {
			expect.WillReturnError(sql.ErrConnDone)
		} else {
			expect.WillReturnRows(rows)
		}

		files, err := List(ts.conn, int64(folderID))
		if !tc.WithMockErr {
			assert.NoError(ts.T(), err)
			assert.GreaterOrEqual(ts.T(), 1, len(files))
		} else {
			assert.EqualError(ts.T(), sql.ErrConnDone, err.Error())
		}
	}
}

func (ts *TransactionSuite) TestListRoot() {
	defer ts.conn.Close()

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

	ts.mock.ExpectQuery(expectedSQL).
		WillReturnRows(rows)

	files, err := ListRootFiles(ts.conn)

	assert.NoError(ts.T(), err)
	assert.GreaterOrEqual(ts.T(), 1, len(files))
}

func SetMockListRootFiles(mock sqlmock.Sqlmock, withMockErr bool) {
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
