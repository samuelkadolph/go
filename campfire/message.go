package campfire

import (
	"fmt"
)

type Message struct {
	Body      *string `json:"body,omitempty"`
	CreatedAt string  `json:"created_at,omitempty"`
	ID        int64   `json:"id,omitempty"`
	RoomID    int64   `json:"room_id,omitempty"`
	Starred   bool    `json:"starred,omitempty"`
	Type      string  `json:"type,omitempty"`
	UserID    *int64  `json:"user_id,omitempty"`

	client *Client
}

func (m *Message) Highlight() error {
	return m.client.post(fmt.Sprintf("/messages/%d/star", m.ID), nil)
}

func (m *Message) Unhighlight() error {
	return m.client.delete(fmt.Sprintf("/messages/%d/star", m.ID), nil)
}
