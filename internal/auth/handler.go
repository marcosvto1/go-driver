package auth

import "net/http"

type Authenticated interface {
	GetId() int64
	GetName() string
}

type authenticateFunc func(string, string) (Authenticated, error)

type handler struct {
	authenticate authenticateFunc
}

func HandleAuth(fn authenticateFunc) func(http.ResponseWriter, *http.Request) {
	h := handler{
		authenticate: fn,
	}

	return h.auth
}
