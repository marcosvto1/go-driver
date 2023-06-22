package files

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestGet() {

	defer ts.conn.Close()

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

	ts.mock.ExpectQuery(
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
	).WithArgs(1).WillReturnRows(rows)

	fr, err := Get(ts.conn, int64(1))
	assert.NoError(ts.T(), err)

	assert.Equal(ts.T(), fr.Name, "any_name")
	assert.Nil(ts.T(), err)
}
