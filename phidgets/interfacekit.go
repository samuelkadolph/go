package phidgets

import (
	"github.com/samuelkadolph/go/phidgets/raw"
)

type InterfaceKit struct {
	Phidget

	Inputs  []*InterfaceKitInput
	Outputs []*InterfaceKitOutput
	Sensors []*InterfaceKitSensor

	rawIFK *raw.InterfaceKit
}

type InterfaceKitInput struct {
	Changed <-chan int
	Index   int

	changed chan int
	ifk     *InterfaceKit
}

type InterfaceKitOutput struct {
	Changed <-chan int
	Index   int

	changed chan int
	ifk     *InterfaceKit
}

type InterfaceKitSensor struct {
	Changed <-chan int
	Index   int

	changed chan int
	ifk     *InterfaceKit
}

func NewInterfaceKit() (*InterfaceKit, error) {
	ifk := new(InterfaceKit)

	r, err := raw.NewInterfaceKit()
	if err != nil {
		return nil, err
	}

	if err := ifk.initInterfaceKit(r); err != nil {
		return nil, err
	}

	return ifk, nil
}

func (o *InterfaceKitOutput) SetState(state bool) error {
	return o.ifk.rawIFK.SetOutputState(o.Index, boolToState(state))
}

func (s *InterfaceKitSensor) SetChangeTrigger(trigger int) error {
	return s.ifk.rawIFK.SetSensorChangeTrigger(s.Index, trigger)
}

func (s *InterfaceKitSensor) SetDataRate(rate int) error {
	return s.ifk.rawIFK.SetDataRate(s.Index, rate)
}

func (s *InterfaceKitSensor) Value() (int, error) {
	return s.ifk.rawIFK.GetSensorValue(s.Index)
}

func boolToState(b bool) int {
	if b {
		return 1
	}
	return 0
}

func newInterfaceKitOutput(n int, i *InterfaceKit) *InterfaceKitOutput {
	o := new(InterfaceKitOutput)

	o.changed = make(chan int)
	o.Changed = o.changed
	o.Index = n
	o.ifk = i

	return o
}

func newInterfaceKitSensor(n int, i *InterfaceKit) *InterfaceKitSensor {
	s := new(InterfaceKitSensor)

	s.changed = make(chan int)
	s.Changed = s.changed
	s.Index = n
	s.ifk = i

	return s
}

func (i *InterfaceKit) createInputs() {
}

func (i *InterfaceKit) createOutputs() {
	if i.Outputs != nil {
		return
	}

	c, err := i.rawIFK.GetOutputCount()
	if err != nil {
		return
	}

	i.Outputs = make([]*InterfaceKitOutput, c)
	for n := 0; n < c; n++ {
		i.Outputs[n] = newInterfaceKitOutput(n, i)
	}

	go func() {
		for c := range i.rawIFK.OutputChanged {
			i.Outputs[c.Index].changed <- c.Value
		}
	}()
}

func (i *InterfaceKit) createSensors() {
	if i.Sensors != nil {
		return
	}

	c, err := i.rawIFK.GetSensorCount()
	if err != nil {
		return
	}

	i.Sensors = make([]*InterfaceKitSensor, c)
	for n := 0; n < c; n++ {
		i.Sensors[n] = newInterfaceKitSensor(n, i)
	}

	go func() {
		for c := range i.rawIFK.SensorChanged {
			i.Sensors[c.Index].changed <- c.Value
		}
	}()
}

func (i *InterfaceKit) initInterfaceKit(r *raw.InterfaceKit) error {
	i.rawIFK = r

	if err := i.initPhidget(&r.Phidget); err != nil {
		return err
	}

	i.attached = func() {
		i.createInputs()
		i.createOutputs()
		i.createSensors()
	}

	return nil
}
