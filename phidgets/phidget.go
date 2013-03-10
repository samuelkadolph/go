package phidgets

import (
	"github.com/samuelkadolph/go/phidgets/raw"
	"time"
)

type Phidget struct {
	Attached     <-chan bool
	Connected    <-chan bool
	Detached     <-chan bool
	Disconnected <-chan bool
	Error        <-chan error

	attached   func()
	detached   func()
	rawPhidget *raw.Phidget
}

func (p *Phidget) Class() (raw.Class, error) {
	return p.rawPhidget.GetDeviceClass()
}

func (p *Phidget) Close() error {
	return p.rawPhidget.Close()
}

func (p *Phidget) ID() (raw.ID, error) {
	return p.rawPhidget.GetDeviceID()
}

func (p *Phidget) Label() (string, error) {
	return p.rawPhidget.GetDeviceLabel()
}

func (p *Phidget) Name() (string, error) {
	return p.rawPhidget.GetDeviceName()
}

func (p *Phidget) Open(c Connector) error {
	return c.Open(p)
}

func (p *Phidget) Serial() (int, error) {
	return p.rawPhidget.GetSerialNumber()
}

func (p *Phidget) ServerAddress() (string, int, error) {
	return p.rawPhidget.GetServerAddress()
}

func (p *Phidget) ServerID() (string, error) {
	return p.rawPhidget.GetServerID()
}

func (p *Phidget) ServerStatus() (int, error) {
	return p.rawPhidget.GetServerStatus()
}

func (p *Phidget) SetLabel(label string) error {
	return p.rawPhidget.SetDeviceLabel(label)
}

func (p *Phidget) Status() (int, error) {
	return p.rawPhidget.GetDeviceStatus()
}

func (p *Phidget) Type() (string, error) {
	return p.rawPhidget.GetDeviceType()
}

func (p *Phidget) Version() (int, error) {
	return p.rawPhidget.GetDeviceVersion()
}

func (p *Phidget) WaitForAttachment(timeout time.Duration) error {
	return p.rawPhidget.WaitForAttachment(timeout)
}

func (p *Phidget) initPhidget(r *raw.Phidget) error {
	p.rawPhidget = r

	attached := make(chan bool)
	connected := make(chan bool)
	detached := make(chan bool)
	disconnected := make(chan bool)

	p.Attached = attached
	p.Connected = connected
	p.Detached = detached
	p.Disconnected = disconnected
	p.Error = p.rawPhidget.Error

	go func() {
		for v := range p.rawPhidget.Attached {
			if p.attached != nil {
				p.attached()
			}
			attached <- v
		}
	}()

	go func() {
		for v := range p.rawPhidget.Connected {
			if p.detached != nil {
				p.detached()
			}
			connected <- v
		}
	}()

	go func() {
		for v := range p.rawPhidget.Detached {
			detached <- v
		}
	}()

	go func() {
		for v := range p.rawPhidget.Disconnected {
			disconnected <- v
		}
	}()

	return nil
}
