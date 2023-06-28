package folders

import (
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestRootSubFolder() {
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

	ts.mock.ExpectQuery(expectedSQL).
		WillReturnRows(rows)

	folders, err := getRootSubFolder(ts.conn)

	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), len(folders), 1)
}
