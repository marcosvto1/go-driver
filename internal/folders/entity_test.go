package folders

import "github.com/stretchr/testify/assert"

func (ts *TransactionSuite) TestShouldReturnErrNameRequired() {
	_, err := New("", 1)

	assert.EqualError(ts.T(), err, ErrNameRequired.Error())
}

func (ts *TransactionSuite) TestShouldCreateNewFolderEntity() {
	e, err := New("Any Name", 1)

	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), e.Name, "Any Name")
}
