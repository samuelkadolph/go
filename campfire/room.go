package campfire

import (
	"fmt"
	"strconv"
	"time"
)

type Room struct {
	ActiveTokenValue string  `json:"active_token_value"`
	CreatedAt        string  `json:"created_at"`
	Full             bool    `json:"full"`
	ID               int64   `json:"id"`
	MembershipLimit  int     `json:"membership_limit"`
	Name             string  `json:"name"`
	OpenToGuests     bool    `json:"open_to_guests"`
	Topic            string  `json:"topic"`
	UpdatedAt        string  `json:"updated_at"`
	Users            []*User `json:"users"`

	client *Client
}

func (r *Room) Paste(body string) error {
	return r.message(body, "PasteMessage")
}

func (r *Room) Recent() ([]*Message, error) {
	return r.recent(nil)
}

func (r *Room) RecentSinceID(id int64) ([]*Message, error) {
	return r.recent(map[string]string{"since_message_id": strconv.FormatInt(id, 10)})
}

func (r *Room) RecentSinceIDWithLimit(id int64, limit int) ([]*Message, error) {
	return r.recent(map[string]string{"since_message_id": strconv.FormatInt(id, 10), "label": strconv.Itoa(limit)})
}

func (r *Room) RecentWithLimit(limit int) ([]*Message, error) {
	return r.recent(map[string]string{"label": strconv.Itoa(limit)})
}

func (r *Room) Say(message string) error {
	return r.message(message, "TextMessage")
}

func (r *Room) Sound(sound string) error {
	return r.message(sound, "SoundMessage")
}

func (r *Room) Transcripts() ([]*Message, error) {
	return r.getMessages(fmt.Sprintf("/room/%d/transcript", r.ID), nil)
}

func (r *Room) TranscriptsForDate(date time.Time) ([]*Message, error) {
	return r.getMessages(fmt.Sprintf("/room/%d/transcript/%s", r.ID, date.Format("2006/01/02")), nil)
}

func (r *Room) Tweet(tweet string) error {
	return r.message(tweet, "TweetMessage")
}

func (r *Room) Upload(body []byte) error {
	return nil
}

func (r *Room) UploadFile(path string) error {
	return nil
}

func (r *Room) Uploads() (uploads []*Upload, err error) {
	var wrapper struct {
		Uploads []*Upload `json:"uploads"`
	}
	if err = r.client.get(fmt.Sprintf("/room/%d/uploads", r.ID), nil, &wrapper); err != nil {
		return
	}
	uploads = wrapper.Uploads
	return
}

func (r *Room) message(body, t string) error {
	var wrapper struct {
		Message Message `json:"message"`
	}
	wrapper.Message.Body = &body
	wrapper.Message.Type = t
	return r.client.post(fmt.Sprintf("/room/%d/speak", r.ID), wrapper)
}

func (r *Room) recent(params map[string]string) ([]*Message, error) {
	return r.getMessages(fmt.Sprintf("/room/%d/recent", r.ID), params)
}

func (r *Room) getMessages(path string, params map[string]string) (msgs []*Message, err error) {
	var wrapper struct {
		Messages []*Message `json:"messages"`
	}
	if err = r.client.get(path, params, &wrapper); err != nil {
		return
	}
	for _, msg := range wrapper.Messages {
		msg.client = r.client
		msgs = append(msgs, msg)
	}
	return
}
