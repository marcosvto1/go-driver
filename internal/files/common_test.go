package files

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TransactionSuite struct {
	suite.Suite
	conn    *sql.DB
	mock    sqlmock.Sqlmock
	entity  *File
	handler handler
}

func (ts *TransactionSuite) SetupTest() {
	var err error

	ts.conn, ts.mock, err = sqlmock.New()
	assert.NoError(ts.T(), err)

	ts.handler = handler{
		db: ts.conn,
	}

	ts.entity = &File{
		Name:     "file.csv",
		FolderId: 1,
		OwnerId:  1,
		Type:     "application/octet-stream",
		Path:     "/file.csv",
	}
}

func (ts *TransactionSuite) AfterTest() {
	err := ts.mock.ExpectationsWereMet()
	assert.NoError(ts.T(), err)
}

func TestSuiteFile(t *testing.T) {
	suite.Run(t, new(TransactionSuite))
}
