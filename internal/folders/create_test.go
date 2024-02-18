package folders

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func (ts *TransactionSuite) TestCreate() {
	defer ts.conn.Close()

	entityInvalidData := Folder{
		Name:     "",
		ParentID: null.Int{},
	}

	tcs := []struct {
		Desc               string
		ExpectedStatusCode int
		WithMock           bool
		WithMockErr        bool
		Body               any
	}{
		{
			Desc:               "Should return status 500 when body is invalid",
			WithMock:           false,
			WithMockErr:        false,
			ExpectedStatusCode: http.StatusInternalServerError,
			Body:               "",
		},
		{
			Desc:               "Should return 400 if not body valid data",
			WithMock:           false,
			WithMockErr:        false,
			ExpectedStatusCode: http.StatusBadRequest,
			Body:               entityInvalidData,
		},
		{
			Desc:               "Should return 500 when failed to create folder",
			WithMock:           true,
			WithMockErr:        true,
			ExpectedStatusCode: http.StatusInternalServerError,
			Body:               ts.entity,
		},
		{
			Desc:               "Should return 200 when success",
			WithMock:           true,
			WithMockErr:        false,
			ExpectedStatusCode: http.StatusCreated,
			Body:               ts.entity,
		},
	}

	for _, tc := range tcs {
		if tc.WithMock {
			setInsertMock(ts.mock, tc.WithMockErr)
		}

		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(tc.Body)
		assert.NoError(ts.T(), err)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/", &b)

		ts.handler.Create(recorder, request)

		assert.Equal(ts.T(), tc.ExpectedStatusCode, recorder.Result().StatusCode)
	}

}

func (ts *TransactionSuite) TestInsert() {
	defer ts.conn.Close()

	setInsertMock(ts.mock, false)

	_, err := Insert(ts.conn, ts.entity)
	assert.NoError(ts.T(), err)
}

func setInsertMock(mock sqlmock.Sqlmock, withMockErr bool) {
	expectedQuery := `INSERT INTO "folders" ("name", "parent_id", "modified_at")*`
	mExpect := mock.ExpectExec(expectedQuery).
		WithArgs(
			"any_name_folder",
			1,
			sqlmock.AnyArg(),
		)
	if withMockErr {
		mExpect.WillReturnError(sql.ErrNoRows)
	} else {
		mExpect.WillReturnResult(sqlmock.NewResult(1, 1))
	}
}
