package files

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/marcosvto1/go-driver/internal/bucket"
	"github.com/marcosvto1/go-driver/internal/queue"
	"github.com/stretchr/testify/assert"
)

func TestModify(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	file, err := New(1, "Name", "ext", "/your/path")
	if err != nil {
		t.Error(err)
	}

	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"folder_id",
		"owner_id",
		"type",
		"path",
		"created_at",
		"modified_at",
		"deleted",
	}).
		AddRow(1, "any_name", 1, 1, "file", "/any/path", time.Now(), time.Now(), false)

	mock.ExpectQuery(
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

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "files" SET "name"=$1, "modified_at"=$2, "deleted"=$3 where id = $4`)).
		WithArgs("any_name2", sqlmock.AnyArg(), file.Deleted, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mQueue, err := queue.New(queue.MockQueue, nil)
	if err != nil {
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

	var b bytes.Buffer
	f := File{
		Name:    "any_name2",
		OwnerId: 1,
		ID:      1,
		Type:    "file",
	}

	err = json.NewEncoder(&b).Encode(f)
	if err != nil {
		t.Error(err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/{id}", &b)

	// SET CONTEXT
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

	h.Modify(recorder, request)

	assert.Equal(t, http.StatusNoContent, recorder.Result().StatusCode)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	file, err := New(1, "Name", "ext", "/your/path")
	if err != nil {
		t.Error(err)
	}

	mock.ExpectExec(`UPDATE "files" SET`).
		WithArgs(file.Name, sqlmock.AnyArg(), file.Deleted, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = Update(db, int64(1), file)
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
