package phidgets

import (
	"github.com/samuelkadolph/go/phidgets/raw"
)

type Connector interface {
	Open(*Phidget) error
}

type Label struct {
	Label string
}

type RemoteIPLabel struct {
	Label    string
	Address  string
	Port     int
	Password string
}

type RemoteIPSerial struct {
	Serial   int
	Address  string
	Port     int
	Password string
}

type RemoteLabel struct {
	Label    string
	Server   string
	Password string
}

type RemoteSerial struct {
	Serial   int
	Server   string
	Password string
}

type Serial struct {
	Serial int
}

var (
	Any = Serial{raw.Any}
	Remote = RemoteSerial{raw.Any, "", ""}
)

func (c Label) Open(p *Phidget) error {
	return p.rawPhidget.OpenLabel(c.Label)
}

func (c RemoteIPLabel) Open(p *Phidget) error {
	return p.rawPhidget.OpenLabelRemoteIP(c.Label, c.Address, c.Port, c.Password)
}

func (c RemoteIPSerial) Open(p *Phidget) error {
	return p.rawPhidget.OpenRemoteIP(c.Serial, c.Address, c.Port, c.Password)
}

func (c RemoteLabel) Open(p *Phidget) error {
	return p.rawPhidget.OpenLabelRemote(c.Label, c.Server, c.Password)
}

func (c RemoteSerial) Open(p *Phidget) error {
	return p.rawPhidget.OpenRemote(c.Serial, c.Server, c.Password)
}

func (c Serial) Open(p *Phidget) error {
	return p.rawPhidget.Open(c.Serial)
}
