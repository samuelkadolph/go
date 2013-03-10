package campfire

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type connection struct {
	Config    Config
	Password  string
	Subdomain string
	Username  string

	http *http.Client
}

type hash map[string]string

func newConnection(subdomain, token string, config Config) *connection {
	return &connection{Subdomain: subdomain, Config: config, Password: "x", Username: token, http: &http.Client{}}
}

func (c *connection) delete(path string, a ...interface{}) error {
	a, p := popHash(a)

	_, err := c.request("DELETE", fmt.Sprintf(path, a...), p, nil)

	return err
}

func (c *connection) get(path string, a ...interface{}) error {
	a, v := popWrapper(path, a)
	a, p := popHash(a)

	res, err := c.request("GET", fmt.Sprintf(path, a...), p, nil)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(v)
}

func (c *connection) post(path string, a ...interface{}) error {
	a, v := popWrapper(path, a)
	a, p := popHash(a)

	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	res, err := c.request("POST", fmt.Sprintf(path, a...), p, body)
	if err != nil {
		return err
	}

	res.Body.Close()

	return nil
}

func (c *connection) put(path string, a ...interface{}) error {
	a, v := popWrapper(path, a)
	a, p := popHash(a)

	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	res, err := c.request("PUT", fmt.Sprintf(path, a...), p, body)
	if err != nil {
		return err
	}

	res.Body.Close()

	return nil
}

func (c *connection) request(method string, path string, params hash, body []byte) (*http.Response, error) {
	u := &url.URL{}
	u.Host = c.Subdomain + ".campfirenow.com"
	u.Path = path
	u.RawQuery = params.toQueryString()

	if c.Config.UseSSL {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(c.Username, c.Password)

	if body != nil {
		req.ContentLength = int64(len(body))
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)

	switch resp.StatusCode {
	case 200:
		return resp, err
	}

	resp.Body.Close()

	return nil, errors.New(resp.Status)
}

func (m hash) toQueryString() string {
	var kvp []string
	for key, value := range m {
		kvp = append(kvp, fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(value)))
	}
	return strings.Join(kvp, "&")
}

func popHash(a []interface{}) ([]interface{}, hash) {
	var h hash

	if len(a) > 0 {
		l := len(a) - 1
		if m, ok := a[l].(map[string]string); ok {
			h = hash(m)
			a = a[:l]
		}
	}

	return a, h
}

func popWrapper(p string, a []interface{}) ([]interface{}, interface{}) {
	var v interface{}

	if len(a) > 0 && strings.Count(p, "%") < len(a) {
		l := len(a) - 1
		v = a[l]
		a = a[:l]
	}

	return a, v
}
