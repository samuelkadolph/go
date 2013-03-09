package campfire

import (
	n "github.com/samuelkadolph/go/nullable"
)

type Message struct {
	Body      n.String `json:"body"`
	CreatedAt n.String `json:"created_at,omitempty"`
	ID        n.Int64  `json:"id,omitempty"`
	RoomID    n.Int64  `json:"room_id,omitempty"`
	Starred   n.Bool   `json:"starred,omitempty"`
	Type      n.String `json:"type"`
	UserID    n.Int64  `json:"user_id,omitempty"`

	connection *connection
}

func (m *Message) Highlight() error {
	return m.connection.post("/messages/%d/star", m.ID)
}

func (m *Message) Unhighlight() error {
	return m.connection.delete("/messages/%d/star", m.ID)
}
