package bucket

import (
	"io"
	"os"
)

type MockBucket struct {
	content map[string][]byte
}

func (m *MockBucket) Upload(reader io.Reader, destiny string) error {
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
