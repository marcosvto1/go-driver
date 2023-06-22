package files

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestModify() {
	defer ts.conn.Close()

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

	ts.mock.ExpectQuery(
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

	ts.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "files" SET "name"=$1, "modified_at"=$2, "deleted"=$3 where id = $4`)).
		WithArgs("any_name2", sqlmock.AnyArg(), false, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var b bytes.Buffer
	f := File{
		Name:    "any_name2",
		OwnerId: 1,
		ID:      1,
		Type:    "file",
	}

	err := json.NewEncoder(&b).Encode(f)
	assert.NoError(ts.T(), err)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/{id}", &b)

	// SET CONTEXT
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))

	ts.handler.Modify(recorder, request)

	assert.Equal(ts.T(), http.StatusNoContent, recorder.Result().StatusCode)
}

func (ts *TransactionSuite) TestUpdate() {
	defer ts.conn.Close()

	ts.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "files" SET "name"=$1, "modified_at"=$2, "deleted"=$3 where id = $4`)).
		WithArgs(ts.entity.Name, sqlmock.AnyArg(), false, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := Update(ts.conn, int64(1), ts.entity)
	assert.NoError(ts.T(), err)
}
