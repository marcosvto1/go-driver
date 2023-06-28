package folders

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
)

func TestModifier(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	expectedSQL := regexp.QuoteMeta(`UPDATE folders SET name=$1, modified_at=$2 WHERE id=$3`)
	mock.ExpectExec(expectedSQL).
		WithArgs(
			"any_name",
			sqlmock.AnyArg(),
			1,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	columns := []string{
		"id",
		"name",
		"parent_id",
		"created_at",
		"modified_at",
		"deleted",
	}

	rows := sqlmock.NewRows(columns).
		AddRow(
			1,
			"any_name",
			1,
			time.Now(),
			time.Now(),
			false,
		)

	expectedSQLSelect := regexp.QuoteMeta(`
		SELECT
			id,
			name,
			parent_id,
			created_at,
			modified_at,
			deleted
		FROM "folders" where id=$1`)

	mock.ExpectQuery(expectedSQLSelect).
		WithArgs(1).
		WillReturnRows(rows)

	var b bytes.Buffer
	err = json.NewEncoder(&b).Encode(&Folder{
		Name: "any_name",
	})
	if err != nil {
		t.Error(err)
	}

	record := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/{id}", &b)

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")

	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

	h := handler{
		db,
	}

	h.Modify(record, request)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(t)
	}
	defer db.Close()

	folder, err := New("any_name_folder", 1)
	if err != nil {
		t.Error(err)
	}

	expectedSQL := `UPDATE folders SET *`
	mock.ExpectExec(expectedSQL).
		WithArgs(
			folder.Name,
			sqlmock.AnyArg(),
			1,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = Update(db, int64(1), folder)
	if err != nil {
		t.Error(err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}
