package files

import (
	"bytes"
	"database/sql"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/marcosvto1/go-driver/internal/bucket"
	"github.com/marcosvto1/go-driver/internal/queue"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestCreate() {
	bucketWithoutErr := ts.handler.bucket
	queueWithoutErr := ts.handler.queue

	tcs := []struct {
		Desc               string
		ExpectedStatusCode int
		FolderId           string
		Filename           string
		WithMock           bool
		WithMockErr        bool
		WithNoFile         bool
		WithBucketErr      bool
		WithQueueErr       bool
	}{
		{
			Desc:               "it should respond with Http status 200",
			ExpectedStatusCode: http.StatusCreated,
			FolderId:           "1",
			WithMock:           true,
			WithMockErr:        false,
			WithNoFile:         false,
			WithBucketErr:      false,
			WithQueueErr:       false,
			Filename:           "file.csv",
		},
		{
			Desc:               "it should respond with Http status 500 when form file not found",
			ExpectedStatusCode: http.StatusInternalServerError,
			FolderId:           "1",
			WithMock:           false,
			WithMockErr:        false,
			WithNoFile:         true,
			WithBucketErr:      false,
			WithQueueErr:       false,
			Filename:           "file.csv",
		},
		{
			Desc:               "it should respond with Http status 500 when form file not found",
			ExpectedStatusCode: http.StatusInternalServerError,
			FolderId:           "1",
			WithMock:           false,
			WithMockErr:        false,
			WithNoFile:         false,
			WithBucketErr:      true,
			WithQueueErr:       false,
			Filename:           "file.csv",
		},

		{
			Desc:               "it should respond with Http status 500 if folder_id is invalid",
			ExpectedStatusCode: http.StatusInternalServerError,
			FolderId:           "A",
			WithMock:           false,
			WithMockErr:        false,
			WithNoFile:         false,
			WithBucketErr:      false,
			WithQueueErr:       false,
			Filename:           "file.csv",
		},

		{
			Desc:               "it should respond with Http status 500 when InsertOne returns err",
			ExpectedStatusCode: http.StatusInternalServerError,
			FolderId:           "1",
			WithMock:           true,
			WithMockErr:        true,
			WithNoFile:         false,
			WithBucketErr:      false,
			WithQueueErr:       false,
			Filename:           "file.csv",
		},

		{
			Desc:               "it should respond with Http status 500 when queue publish returns err",
			ExpectedStatusCode: http.StatusInternalServerError,
			FolderId:           "1",
			WithMock:           true,
			WithMockErr:        false,
			WithNoFile:         false,
			WithBucketErr:      false,
			WithQueueErr:       true,
			Filename:           "file.csv",
		},
	}
	defer ts.conn.Close()

	for _, tc := range tcs {
		ts.T().Log(tc.Desc)

		if tc.WithMock {
			setMockInsert(ts.mock, ts.entity, tc.WithMockErr)
		}

		if tc.WithBucketErr {
			buckWithErr, err := bucket.New(bucket.MockProvider, bucket.MockBucketConfig{
				UpdateWillReturnErr: true,
			})
			assert.NoError(ts.T(), err)

			ts.handler.bucket = buckWithErr
		} else {
			ts.handler.bucket = bucketWithoutErr
		}

		if tc.WithQueueErr {
			mQueueWithErr, err := queue.New(queue.MockQueue, queue.MockQueueConfig{
				PublishWillReturnErr: true,
			})
			assert.NoError(ts.T(), err)

			ts.handler.queue = mQueueWithErr
		} else {
			ts.handler.queue = queueWithoutErr
		}

		var body *bytes.Buffer
		var writer *multipart.Writer

		body = new(bytes.Buffer)
		writer = multipart.NewWriter(body)

		if !tc.WithNoFile {
			// Alternative Upload with local file
			file, err := os.Open("./testdata/file.csv")
			assert.NoError(ts.T(), err)

			wf, err := writer.CreateFormFile("file", tc.Filename)
			assert.NoError(ts.T(), err)

			_, err = io.Copy(wf, file)
			assert.NoError(ts.T(), err)

			writer.WriteField("folder_id", tc.FolderId)

			// Alternative with simulate buffer file in memory
			// part, _ := writer.CreateFormFile("file", "file.csv")
			// part.Write([]byte(`sample`))
			writer.Close() // Close multipart writer
		}

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/", body)

		if !tc.WithNoFile {
			request.Header.Add("Content-Type", writer.FormDataContentType())
		}

		ts.handler.Create(recorder, request)

		assert.Equal(ts.T(), tc.ExpectedStatusCode, recorder.Result().StatusCode)
	}

}

func (ts *TransactionSuite) TestInsertOne() {
	defer ts.conn.Close()

	setMockInsert(ts.mock, ts.entity, false)

	_, err := InsertOne(ts.conn, ts.entity)
	assert.NoError(ts.T(), err)
}

func setMockInsert(mock sqlmock.Sqlmock, entity *File, withMockErr bool) {
	expectedSQL := regexp.QuoteMeta(`INSERT INTO files (folder_id, owner_id, name, type, path, modified_at)
	VALUES ($1, $2, $3, $4, $5, $6)`)

	expect := mock.ExpectExec(expectedSQL).
		WithArgs(
			1,
			1,
			"file.csv",
			"application/octet-stream",
			"/file.csv",
			sqlmock.AnyArg(),
		)

	if !withMockErr {
		expect.WillReturnResult(sqlmock.NewResult(1, 1))
	} else {
		expect.WillReturnError(sql.ErrConnDone)
	}
}
