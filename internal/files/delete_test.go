package files

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/marcosvto1/go-driver/internal/bucket"
	"github.com/marcosvto1/go-driver/internal/queue"
	"github.com/stretchr/testify/assert"
)

func TestDeleteHTTP(t *testing.T) {
	db, mock, err := mockDeleteDB()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/{id}", nil)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

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

	h.Delete(recorder, request)

	assert.Equal(t, http.StatusNoContent, recorder.Result().StatusCode)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}

func mockDeleteDB() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	mock.ExpectExec(`UPDATE files SET *`).
		WithArgs(true, sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	return db, mock, nil
}

func TestDelete(t *testing.T) {
	db, mock, err := mockDeleteDB()
	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	err = Delete(db, int64(1))
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
