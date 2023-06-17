package users

import (
	"context"
	"database/sql/driver"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestDeleteOne(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")

	// Create a new request with a dummy URL parameter

	mock.ExpectExec(`UPDATE users SET *`).
		WithArgs(AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create a new response recorder
	recorder := httptest.NewRecorder()

	// Create a new request
	request := httptest.NewRequest("DELETE", "/delete/123", nil)

	// Set context for request
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

	// Create a new instance of the handler
	h := &handler{
		db,
	}

	// Call the Delete method of the handler
	h.Delete(recorder, request)

	// Check the response status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
	}

	// Check the response content type header
	expectedContentType := "application/json"
	if recorder.Header().Get("Content-Type") != expectedContentType {
		t.Errorf("Expected content type %s, but got %s", expectedContentType, recorder.Header().Get("Content-Type"))
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	mock.ExpectExec(`UPDATE users SET *`).
		WithArgs(AnyTime{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = Delete(db, int64(1))
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
