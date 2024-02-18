package files

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/marcosvto1/go-driver/internal/bucket"
	"github.com/marcosvto1/go-driver/internal/queue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/guregu/null.v4"
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

	mQueue, err := queue.New(queue.MockQueue, queue.MockQueueConfig{
		PublishWillReturnErr: false,
	})
	assert.NoError(ts.T(), err)

	mBucket, err := bucket.New(bucket.MockProvider, bucket.MockBucketConfig{
		UpdateWillReturnErr: false,
	})
	assert.NoError(ts.T(), err)

	ts.handler = handler{
		db:     ts.conn,
		queue:  mQueue,
		bucket: mBucket,
	}

	ts.entity = &File{
		Name:     "file.csv",
		FolderId: null.NewInt(1, true),
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
