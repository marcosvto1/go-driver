package folders

import (
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestGetFolder() {
	defer ts.conn.Close()

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

	ts.mock.ExpectQuery(expectedSQL).
		WithArgs(1).
		WillReturnRows(rows)

	folder, err := GetFolder(ts.conn, 1)

	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), folder.Name, "any folder name")
}

func (ts *TransactionSuite) TestGetSubFolder() {
	defer ts.conn.Close()

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

	ts.mock.ExpectQuery(expectedSQL).
		WithArgs(1).
		WillReturnRows(rows)

	folders, err := getSubFolder(ts.conn, 1)

	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), len(folders), 2)
}
