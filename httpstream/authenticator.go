package httpstream

import (
	"net/http"
)

type Authenticator interface {
	Authenticate(*http.Request)
}

type BasicAuth struct {
	Username string
	Password string
}

type HeaderAuth struct {
	Value string
}

type NoAuth struct {
}

func (a BasicAuth) Authenticate(r *http.Request) {
	r.SetBasicAuth(a.Username, a.Password)
}

func (a HeaderAuth) Authenticate(r *http.Request) {
	r.Header.Set("Authorization", a.Value)
}

func (a NoAuth) Authenticate(r *http.Request) {
}
