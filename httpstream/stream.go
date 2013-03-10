package httpstream

import (
	"bufio"
	"errors"
	"net/http"
	"time"
)

type Stream struct {
	Data  <-chan []byte
	Error <-chan error

	auth  Authenticator
	data  chan []byte
	end   Endpoint
	error chan error
	pro   Processor
}

func New(end Endpoint, auth Authenticator, pro Processor) (*Stream, error) {
	s := &Stream{}

	s.auth = auth
	s.data = make(chan []byte, 100)
	s.end = end
	s.error = make(chan error, 100)
	s.pro = pro

	s.Data = s.data
	s.Error = s.error

	resp, err := s.connect()
	if err != nil {
		return nil, err
	}

	go s.process(resp)

	return s, nil
}

func (s *Stream) Close() {
}

func (s *Stream) connect() (*http.Response, error) {
	c := &http.Client{}

	req := &http.Request{}
	req.Header = http.Header{}
	s.end.Target(req)
	s.auth.Authenticate(req)

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	return res, nil
}

func (s *Stream) process(res *http.Response) {
	var b []byte
	var err error
	var r *bufio.Reader

	r = bufio.NewReader(res.Body)

	for {
		if b, err = s.pro.Process(r); err != nil {
			res.Body.Close()

			time.Sleep(1 * time.Second)

			if res, err = s.connect(); err != nil {
				continue
			}

			r = bufio.NewReader(res.Body)
			continue
		}

		s.data <- b
	}
}
