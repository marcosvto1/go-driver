package requests

import (
	"io"
	"net/http"
)

func AuthenticatedPost(path string, body io.Reader) (*http.Response, error) {
	return doRequest("POST", path, body, true)
}

func AuthenticatedGet(path string) (*http.Response, error) {
	return doRequest("GET", path, nil, true)
}
