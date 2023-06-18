package folders

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	folder, err := New("any_name_folder", 0)
	if err != nil {
		t.Error(err)
	}

	expectedQuery := `INSERT INTO "folders" ("name", "parent_id", "modified_at")*`
	mock.ExpectExec(expectedQuery).
		WithArgs(
			"any_name_folder",
			0,
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	h := handler{
		db,
	}

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(folder)
	if err != nil {
		t.Error(err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", &b)

	h.Create(recorder, request)

	assert.Equal(t, http.StatusCreated, recorder.Result().StatusCode)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}

}

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	folder, err := New("any_name_folder", 0)
	if err != nil {
		t.Error(err)
	}

	expectedQuery := `INSERT INTO "folders" ("name", "parent_id", "modified_at")*`
	mock.ExpectExec(expectedQuery).
		WithArgs(
			"any_name_folder",
			0,
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = Insert(db, folder)
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
