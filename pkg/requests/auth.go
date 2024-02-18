package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Auth(path, user, pass string) error {
	creds := credentials{
		Username: user,
		Password: pass,
	}

	var body bytes.Buffer

	err := json.NewEncoder(&body).Encode(creds)
	if err != nil {
		return err
	}

	res, err := doRequest("POST", path, &body, nil, false)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return createTokenCache(res.Body)

}

type cacheToken struct {
	Token string `json:"token"`
}

func createTokenCache(body io.ReadCloser) error {
	token, err := io.ReadAll(body)
	if err != nil {
		return nil
	}

	file, err := os.Create(".cacheToken")
	if err != nil {
		return err
	}

	var jsonToken struct {
		Token string `json:"access_token"`
	}

	json.Unmarshal(token, &jsonToken)

	cache := cacheToken{
		Token: string(jsonToken.Token),
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	_, err = file.Write(data)

	return err
}

func readTokenCache() (string, error) {
	data, err := os.ReadFile(".cacheToken")
	if err != nil {
		return "", err
	}

	var cache cacheToken
	err = json.Unmarshal(data, &cache)
	if err != nil {
		return "", err
	}
	return cache.Token, nil
}
