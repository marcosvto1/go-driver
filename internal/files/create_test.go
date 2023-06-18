package files

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/marcosvto1/go-driver/internal/bucket"
	"github.com/marcosvto1/go-driver/internal/queue"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	mock.ExpectExec(`INSERT INTO files (folder_id, owner_id, name, type, path, modified_at)*`).
		WithArgs(
			1,
			1,
			"file.csv",
			"application/octet-stream",
			"/file.csv",
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Alternative Upload with local file
	file, err := os.Open("./testdata/file.csv")
	if err != nil {
		t.Error(err)
	}
	wf, err := writer.CreateFormFile("file", "file.csv")
	if err != nil {
		t.Error(err)
	}

	_, err = io.Copy(wf, file)
	if err != nil {
		t.Error(err)
	}

	writer.WriteField("folder_id", "1")

	// Alternative with simulate buffer file in memory
	// part, _ := writer.CreateFormFile("file", "file.csv")
	// part.Write([]byte(`sample`))
	writer.Close() // Close multipart writer

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", body)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	// Configure Mocks Dependecy
	mQueue, err := queue.New(queue.MockQueue, nil)
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	mBucket, err := bucket.New(bucket.MockProvider, nil)
	if err != nil {
		t.Error(err)
	}

	h := handler{
		db:     db,
		queue:  mQueue,
		bucket: mBucket,
	}

	h.Create(recorder, request)

	t.Log("It should respond with an HTTP status code of 201")
	assert.Equal(t, http.StatusCreated, recorder.Result().StatusCode)
}
func TestInsertOne(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	file, err := New(int64(1), "file", "jpg", "/any/path")
	if err != nil {
		t.Error(err)
	}

	file.FolderId = 1

	mock.ExpectExec(`INSERT INTO files (folder_id, owner_id, name, type, path, modified_at)*`).
		WithArgs(
			1,
			1,
			"file",
			"jpg",
			"/any/path",
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = InsertOne(db, file)
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
