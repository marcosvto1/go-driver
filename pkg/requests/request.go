package requests

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func doRequest(method string, path string, body io.Reader, auth bool) (*http.Response, error) {
	url := fmt.Sprintf("http://localhost:8000%s", path)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if auth {
		token, err := readTokenCache()
		if err != nil {
			log.Println("Error reading token")
			return nil, err
		}

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
