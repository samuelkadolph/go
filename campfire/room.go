package campfire

import (
	n "github.com/samuelkadolph/go/nullable"
	"strconv"
	"strings"
	"time"
)

type Room struct {
	ActiveTokenValue n.String `json:"active_token_value,omitempty"`
	CreatedAt        n.String `json:"created_at,omitempty"`
	Full             n.Bool   `json:"full,omitempty"`
	ID               n.Int64  `json:"id,omitempty"`
	MembershipLimit  n.Int    `json:"membership_limit,omitempty"`
	Name             n.String `json:"name,omitempty"`
	OpenToGuests     n.Bool   `json:"open_to_guests,omitempty"`
	Topic            n.String `json:"topic,omitempty"`
	UpdatedAt        n.String `json:"updated_at,omitempty"`
	Users            []*User  `json:"users,omitempty"`

	connection *connection
}

const (
	_56K          = "56k"
	BUELLER       = "bueller"
	CRICKETS      = "crickets"
	DANGERZONE    = "dangerzone"
	DEEPER        = "deeper"
	DRAMA         = "drama"
	GREATJOB      = "greatjob"
	HORN          = "horn"
	HORROR        = "horror"
	INCONCEIVABLE = "inconceivable"
	LIVE          = "live"
	LOGGINS       = "loggins"
	NOOOO         = "noooo"
	NYAN          = "nyan"
	OHMY          = "ohmy"
	OHYEAH        = "ohyeah"
	PUSHIT        = "pushit"
	RIMSHOT       = "rimshot"
	SAX           = "sax"
	SECRET        = "secret"
	TADA          = "tada"
	TMYK          = "tmyk"
	TROMBONE      = "trombone"
	VUVUZELA      = "vuvuzela"
	YEAH          = "yeah"
	YODEL         = "yodel"
)

func (r *Room) Join() error {
	return r.connection.post("/room/%d/join", r.ID)
}

func (r *Room) Leave() error {
	return r.connection.post("/room/%d/leave", r.ID)
}

func (r *Room) Lock() error {
	return r.connection.post("/room/%d/lock", r.ID)
}

func (r *Room) Message(body string) error {
	return r.message(body, "")
}

func (r *Room) Paste(paste string) error {
	return r.message(paste, "PasteMessage")
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

func (r *Room) Rename(name string) error {
	var wrapper struct {
		Room Room `json:"room"`
	}

	wrapper.Room.Name = (n.String)(name)

	return r.connection.put("/room/%d", r.ID, wrapper)
}

func (r *Room) SetTopic(topic string) error {
	var wrapper struct {
		Room Room `json:"room"`
	}

	wrapper.Room.Topic = (n.String)(topic)

	return r.connection.put("/room/%d", r.ID, wrapper)
}

func (r *Room) Unlock() error {
	return r.connection.post("/room/%d/unlock", r.ID)
}

func (r *Room) Sound(sound string) error {
	return r.message(sound, "SoundMessage")
}

func (r *Room) Speak(message string) error {
	return r.message(message, "TextMessage")
}

func (r *Room) Transcripts() ([]*Message, error) {
	return r.messages("/room/%d/transcript", r.ID)
}

func (r *Room) TranscriptsForDate(date time.Time) ([]*Message, error) {
	return r.messages("/room/%d/transcript/%s", r.ID, date.Format("2006/01/02"))
}

func (r *Room) Tweet(tweet string) error {
	return r.message(tweet, "TweetMessage")
}

func (r *Room) Uploads() ([]*Upload, error) {
	var wrapper struct {
		Uploads []*Upload `json:"uploads"`
	}

	if err := r.connection.get("/room/%d/uploads", r.ID, &wrapper); err != nil {
		return nil, err
	}

	return wrapper.Uploads, nil
}

func (r *Room) message(b, t string) error {
	var wrapper struct {
		Message Message `json:"message"`
	}

	wrapper.Message.Body = (n.String)(strings.Replace(b, "\n", "&#xA;", -1))

	if t != "" {
		wrapper.Message.Type = (n.String)(t)
	}

	return r.connection.post("/room/%d/speak", r.ID, wrapper)
}

func (r *Room) messages(path string, a ...interface{}) ([]*Message, error) {
	var wrapper struct {
		Messages []*Message `json:"messages"`
	}

	if err := r.connection.get(path, append(a, &wrapper)...); err != nil {
		return nil, err
	}

	msgs := make([]*Message, len(wrapper.Messages))
	for i, m := range wrapper.Messages {
		msgs[i] = m
		msgs[i].connection = r.connection
	}

	return msgs, nil
}

func (r *Room) recent(params map[string]string) ([]*Message, error) {
	return r.messages("/room/%d/recent", r.ID, params)
}
