package bucket

import (
	"errors"
	"io"
	"os"
)

type MockBucketConfig struct {
	UpdateWillReturnErr bool
}

type MockBucket struct {
	content     map[string][]byte
	mockOptions MockBucketConfig
}

func (m *MockBucket) Upload(reader io.Reader, destiny string) error {

	if m.mockOptions.UpdateWillReturnErr {
		return errors.New("upload err")
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	m.content[destiny] = data

	return nil
}

func (m *MockBucket) Download(source string, destiny string) (*os.File, error) {
	return nil, nil
}

func (m *MockBucket) Delete(string) error {
	return nil
}
