package requests

import (
	"errors"
	"io"
	"net/http"
)

func validateResponse(res *http.Response) ([]byte, error) {
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, errors.New(string(data))
	}

	return data, nil
}

func AuthenticatedPost(path string, body io.Reader) ([]byte, error) {
	res, err := doRequest("POST", path, body, nil, true)
	if err != nil {
		return nil, err
	}

	return validateResponse(res)
}

func Post(path string, body io.Reader) ([]byte, error) {
	res, err := doRequest("POST", path, body, nil, false)
	if err != nil {
		return nil, err
	}

	return validateResponse(res)
}

func AuthenticatedWithHeaders(path string, body io.Reader, headers map[string]string) ([]byte, error) {
	res, err := doRequest("POST", path, body, headers, true)
	if err != nil {
		return nil, err
	}

	return validateResponse(res)
}

func AuthenticatedPut(path string, body io.Reader) ([]byte, error) {
	res, err := doRequest("PUT", path, body, nil, true)
	if err != nil {
		return nil, err
	}

	return validateResponse(res)
}

func AuthenticatedDelete(path string, body io.Reader) ([]byte, error) {
	res, err := doRequest("DELETE", path, nil, nil, true)
	if err != nil {
		return nil, err
	}

	return validateResponse(res)
}

func AuthenticatedGet(path string) ([]byte, error) {
	res, err := doRequest("GET", path, nil, nil, true)
	if err != nil {
		return nil, err
	}

	return validateResponse(res)
}
