package requests

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func doRequest(method string, path string, body io.Reader, headers map[string]string, auth bool) (*http.Response, error) {
	url := fmt.Sprintf("http://localhost:3000%s", path)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if _, ok := headers["Content-Type"]; !ok {
		req.Header.Add("Content-Type", "application/json")
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
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
