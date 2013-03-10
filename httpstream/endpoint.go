package httpstream

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Endpoint interface {
	Target(*http.Request)
}

type Get struct {
	URL *url.URL
}

type Post struct {
	URL         *url.URL
	Body        []byte
	ContentType string
}

type Put struct {
	Post
}

func (e Get) Target(r *http.Request) {
	r.Method = "GET"
	r.URL = e.URL
}

func (e Post) Target(r *http.Request) {
	r.Body = ioutil.NopCloser(bytes.NewBuffer(e.Body))
	r.ContentLength = int64(len(e.Body))
	r.Header.Set("Content-Type", e.ContentType)
	r.Method = "POST"
	r.URL = e.URL
}

func (e Put) Target(r *http.Request) {
	e.Post.Target(r)
	r.Method = "PUT"
}
