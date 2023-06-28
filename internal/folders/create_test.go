package folders

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (ts *TransactionSuite) TestCreate() {
	defer ts.conn.Close()

	expectedQuery := `INSERT INTO "folders" ("name", "parent_id", "modified_at")*`
	ts.mock.ExpectExec(expectedQuery).
		WithArgs(
			"any_name_folder",
			1,
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(ts.entity)
	assert.NoError(ts.T(), err)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", &b)

	ts.handler.Create(recorder, request)

	assert.Equal(ts.T(), http.StatusCreated, recorder.Result().StatusCode)
}

func (ts *TransactionSuite) TestInsert() {
	defer ts.conn.Close()

	expectedSQL := regexp.QuoteMeta(`INSERT INTO "folders" ("name", "parent_id", "modified_at") VALUES ($1, $2, $3)`)
	ts.mock.ExpectExec(expectedSQL).
		WithArgs(
			"any_name_folder",
			1,
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := Insert(ts.conn, ts.entity)
	assert.NoError(ts.T(), err)
}
