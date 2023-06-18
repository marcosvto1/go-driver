package folders

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
)

func TestDeleteHTTP(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/{id}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

	expectedQuery := regexp.QuoteMeta(`SELECT
	id,
	name,
	folder_id,
	owner_id,
	type,
	path,
	created_at,
	modified_at,
	deleted
	FROM files
WHERE folder_id = $1 AND deleted = false`)

	folderID := 1

	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"folder_id",
		"owner_id",
		"type",
		"path",
		"created_at",
		"modified_at",
		"deleted_at",
	}).AddRow(
		1,
		"any_name",
		1,
		1,
		"file",
		"/any/path",
		time.Now(),
		time.Now(),
		false,
	)

	mock.ExpectQuery(expectedQuery).
		WithArgs(folderID).
		WillReturnRows(rows)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "files" SET "name"=$1, "modified_at"=$2, "deleted"=$3 where id = $4`)).
		WithArgs("any_name", sqlmock.AnyArg(), true, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	columns := []string{
		"id",
		"name",
		"parent_id",
		"created_at",
		"modified_at",
		"deleted",
	}

	rootrows := sqlmock.NewRows(columns).
		AddRow(
			2,
			"any folder name",
			2,
			time.Now(),
			time.Now(),
			false,
		).
		AddRow(
			2,
			"any folder name 2",
			3,
			time.Now(),
			time.Now(),
			false,
		)

	expectedSQL := regexp.QuoteMeta(`
		SELECT
			id,
			name,
			parent_id,
			created_at,
			modified_at,
			deleted
		FROM "folders" where parent_id=$1`,
	)

	mock.ExpectQuery(expectedSQL).
		WithArgs(1).
		WillReturnRows(rootrows)

	expectedQueryDeleteFolder := regexp.QuoteMeta(`UPDATE "folders" SET "modified_at"=$1, "deleted"=true WHERE id=$2`)

	mock.ExpectExec(expectedQueryDeleteFolder).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	h := handler{
		db: db,
	}

	h.Delete(recorder, request)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	defer db.Close()

	expectedQuery := regexp.QuoteMeta(`UPDATE "folders" SET "modified_at"=$1, "deleted"=true WHERE id=$2`)

	mock.ExpectExec(expectedQuery).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = Delete(db, int64(1))
	if err != nil {
		t.Error(t)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(t)
	}
}
