package campfire

import (
	"encoding/json"
	"fmt"
	"github.com/samuelkadolph/go/httpstream"
	"net/url"
)

type Stream struct {
	Error   <-chan error
	Message <-chan Message
	Room    *Room

	connection *connection
	error      chan error
	message    chan Message
	stream     *httpstream.Stream
}

func NewStream(room *Room) (*Stream, error) {
	var err error

	s := &Stream{}
	s.connection = room.connection
	s.error = make(chan error)
	s.message = make(chan Message)
	s.Error = s.error
	s.Message = s.message

	u := &url.URL{}
	u.Host = "streaming.campfirenow.com"
	u.Path = fmt.Sprintf("/room/%d/live.json", room.ID)

	if room.connection.Config.UseSSL {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}

	s.stream, err = httpstream.New(httpstream.Get{u}, httpstream.BasicAuth{s.connection.Username, s.connection.Password}, httpstream.CarriageReturn)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case d := <-s.stream.Data:
				var m Message
				err := json.Unmarshal(d, &m)
				if err != nil {
					s.error <- err
				}
				m.connection = s.connection
				s.message <- m
			case e := <-s.stream.Error:
				s.error <- e
			}
		}
	}()

	return s, nil
}
