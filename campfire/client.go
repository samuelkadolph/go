package campfire

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	Subdomain string
	Token     string

	http *http.Client
}

func NewClient(subdomain, token string) *Client {
	return &Client{Subdomain: subdomain, Token: token, http: &http.Client{}}
}

func (c *Client) Account() (a *Account, err error) {
	var wrapper struct {
		Account Account `json:"account"`
	}
	if err = c.get("/account", nil, &wrapper); err != nil {
		return
	}
	a = &wrapper.Account
	return
}

func (c *Client) RoomByID(id int) (room *Room, err error) {
	var wrapper struct {
		Room Room `json:"room"`
	}
	if err = c.get(fmt.Sprintf("/room/%d", id), nil, &wrapper); err != nil {
		return
	}
	room = &wrapper.Room
	room.client = c
	return
}

func (c *Client) RoomByName(name string) (*Room, error) {
	rooms, err := c.Rooms()
	if err != nil {
		return nil, err
	}
	for _, room := range rooms {
		if room.Name == name {
			return room, nil
		}
	}
	return nil, nil
}

func (c *Client) Rooms() (rooms []*Room, err error) {
	var wrapper struct {
		Rooms []*Room
	}
	if err = c.get("/rooms", nil, &wrapper); err != nil {
		return
	}
	for _, room := range wrapper.Rooms {
		room.client = c
		rooms = append(rooms, room)
	}
	return
}

func (c *Client) Search(term string) (msgs []*Message, err error) {
	var wrapper struct {
		Messages []*Message `json:"messages"`
	}
	if err = c.get("/search", map[string]string{"q": term}, &wrapper); err != nil {
		return
	}
	for _, msg := range wrapper.Messages {
		msg.client = c
		msgs = append(msgs, msg)
	}
	return
}

func (c *Client) delete(path string, params map[string]string) error {
	_, err := c.request("DELETE", path, params, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) get(path string, params map[string]string, v interface{}) error {
	res, err := c.request("GET", path, params, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	dec := json.NewDecoder(res.Body)
	return dec.Decode(v)
}

func (c *Client) post(path string, v interface{}) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	res, err := c.request("POST", path, nil, body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (c *Client) request(method, path string, params map[string]string, body []byte) (*http.Response, error) {
	u := &url.URL{}
	u.Scheme = "https"
	u.Host = c.Subdomain + ".campfirenow.com"
	u.Path = path

	var kvp []string
	for key, value := range params {
		kvp = append(kvp, fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(value)))
	}
	u.RawQuery = strings.Join(kvp, "&")

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(c.Token, "x")

	if body != nil {
		req.ContentLength = int64(len(body))
		req.Header.Add("Content-Type", "application/json")
	}

	return c.http.Do(req)
}
