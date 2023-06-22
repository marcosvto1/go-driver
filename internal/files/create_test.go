package files

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestCreate() {
	defer ts.conn.Close()

	setMockInsert(ts.mock, ts.entity)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Alternative Upload with local file
	file, err := os.Open("./testdata/file.csv")
	assert.NoError(ts.T(), err)

	wf, err := writer.CreateFormFile("file", "file.csv")
	assert.NoError(ts.T(), err)

	_, err = io.Copy(wf, file)
	assert.NoError(ts.T(), err)

	writer.WriteField("folder_id", "1")

	// Alternative with simulate buffer file in memory
	// part, _ := writer.CreateFormFile("file", "file.csv")
	// part.Write([]byte(`sample`))
	writer.Close() // Close multipart writer

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", body)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	ts.handler.Create(recorder, request)

	ts.T().Log("It should respond with an HTTP status code of 201")
	assert.Equal(ts.T(), http.StatusCreated, recorder.Result().StatusCode)
}

func (ts *TransactionSuite) TestInsertOne() {
	defer ts.conn.Close()

	setMockInsert(ts.mock, ts.entity)

	_, err := InsertOne(ts.conn, ts.entity)
	assert.NoError(ts.T(), err)
}

func setMockInsert(mock sqlmock.Sqlmock, entity *File) {
	expectedSQL := regexp.QuoteMeta(`INSERT INTO files (folder_id, owner_id, name, type, path, modified_at)
	VALUES ($1, $2, $3, $4, $5, $6)`)

	mock.ExpectExec(expectedSQL).
		WithArgs(
			1,
			1,
			"file.csv",
			"application/octet-stream",
			"/file.csv",
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
}
