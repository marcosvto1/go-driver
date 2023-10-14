package files

import "github.com/stretchr/testify/assert"

func (ts *TransactionSuite) TestShouldReturnErrNameRequired() {
	_, err := New(1, "", "", "")
	assert.EqualError(ts.T(), err, ErrNameRequired.Error())
}

func (ts *TransactionSuite) TestShouldReturnOwnerRequired() {
	_, err := New(0, "Any name", "", "")
	assert.EqualError(ts.T(), err, ErrOwnerRequired.Error())
}

func (ts *TransactionSuite) TestShouldPathRequired() {
	_, err := New(1, "Any name", "", "")
	assert.EqualError(ts.T(), err, ErrPathRequired.Error())
}

func (ts *TransactionSuite) TestShouldTypeRequired() {
	_, err := New(1, "Any name", "", "any/path")
	assert.EqualError(ts.T(), err, ErrTypoRequired.Error())
}
