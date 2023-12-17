package requests

import (
	"io"
)

func AuthenticatedPost(path string, body io.Reader) ([]byte, error) {
	res, err := doRequest("POST", path, body, true)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(res.Body)
}

func AuthenticatedPut(path string, body io.Reader) ([]byte, error) {
	res, err := doRequest("PUT", path, body, true)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(res.Body)
}

func AuthenticatedDelete(path string, body io.Reader) ([]byte, error) {
	res, err := doRequest("DELETE", path, nil, true)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(res.Body)
}

func AuthenticatedGet(path string) ([]byte, error) {
	res, err := doRequest("GET", path, nil, true)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(res.Body)
}
