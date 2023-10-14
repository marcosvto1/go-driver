package files

import (
	"database/sql"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestGet() {
	defer ts.conn.Close()

	tcs := []struct {
		Desc        string
		WithMockErr bool
	}{
		{
			Desc:        "Should returns files the filters",
			WithMockErr: false,
		},
		{
			Desc:        "Should return error if database failed",
			WithMockErr: true,
		},
	}

	for _, tc := range tcs {
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
		}).
			AddRow(1, "any_name", 1, 1, "file", "/any/path", time.Now(), time.Now(), false)

		expect := ts.mock.ExpectQuery(
			`SELECT
			id,
			name,
			folder_id,
			owner_id,
			type,
			path,
			created_at,
			modified_at,
			deleted
		FROM "files" * `,
		).WithArgs(1)

		if tc.WithMockErr {
			expect.WillReturnError(sql.ErrConnDone)
		} else {
			expect.WillReturnRows(rows)
		}

		fr, err := Get(ts.conn, int64(1))

		if !tc.WithMockErr {
			assert.NoError(ts.T(), err)
			assert.Equal(ts.T(), fr.Name, "any_name")
		} else {
			assert.EqualError(ts.T(), sql.ErrConnDone, err.Error())
		}
	}
}
