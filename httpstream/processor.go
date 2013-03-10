package httpstream

import (
	"bufio"
	"bytes"
)

type Processor interface {
	Process(*bufio.Reader) ([]byte, error)
}

type Delimiter struct {
	Delim byte
}

type Raw struct {
	b []byte
}

var CarriageReturn = Delimiter{'\r'}
var Newline = Delimiter{'\n'}

func (p Delimiter) Process(r *bufio.Reader) ([]byte, error) {
	b, err := r.ReadBytes(p.Delim)
	if err != nil {
		return nil, err
	}
	return bytes.TrimSpace(b), nil
}

func (p Raw) Process(r *bufio.Reader) ([]byte, error) {
	if p.b == nil {
		p.b = make([]byte, 10240)
	}

	n, err := r.Read(p.b)
	if err != nil {
		return nil, err
	}

	return p.b[:n], nil
}
