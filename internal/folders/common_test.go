package folders

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/guregu/null.v4"
)

type TransactionSuite struct {
	suite.Suite
	handler handler
	entity  *Folder
	mock    sqlmock.Sqlmock
	conn    *sql.DB
}

func (ts *TransactionSuite) SetupTest() {
	var err error
	ts.conn, ts.mock, err = sqlmock.New()
	assert.NoError(ts.T(), err)

	ts.handler = handler{
		db: ts.conn,
	}

	ts.entity, err = New("any_name_folder", null.NewInt(1, true))
	assert.NoError(ts.T(), err)
}

func (ts *TransactionSuite) AfterTest() {
	err := ts.mock.ExpectationsWereMet()
	assert.NoError(ts.T(), err)
}

func TestSuiteFolder(t *testing.T) {
	suite.Run(t, new(TransactionSuite))
}
