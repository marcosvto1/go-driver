package folders

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func (ts *TransactionSuite) TestShouldReturnErrNameRequired() {
	_, err := New("", null.NewInt(1, true))

	assert.EqualError(ts.T(), err, ErrNameRequired.Error())
}

func (ts *TransactionSuite) TestShouldCreateNewFolderEntity() {
	e, err := New("Any Name", null.NewInt(1, true))

	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), e.Name, "Any Name")
}
