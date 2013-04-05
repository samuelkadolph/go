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

func (i *InterfaceKit) Ratiometric() (bool, error) {
	s, err := i.rawIFK.GetRatiometric()
	return stateToBool(s), err
}

func (i *InterfaceKit) SetRatiometric(state bool) error {
	return i.rawIFK.SetRatiometric(boolToState(state))
}

func (i *InterfaceKitInput) State() (bool, error) {
	s, err := i.ifk.rawIFK.GetInputState(i.Index)
	return stateToBool(s), err
}

func (o *InterfaceKitOutput) SetState(state bool) error {
	return o.ifk.rawIFK.SetOutputState(o.Index, boolToState(state))
}

func (o *InterfaceKitOutput) State() (bool, error) {
	s, err := o.ifk.rawIFK.GetOutputState(o.Index)
	return stateToBool(s), err
}

func (s *InterfaceKitSensor) ChangeTrigger() (int, error) {
	return s.ifk.rawIFK.GetSensorChangeTrigger(s.Index)
}

func (s *InterfaceKitSensor) DataRate() (int, error) {
	return s.ifk.rawIFK.GetDataRate(s.Index)
}

func (s *InterfaceKitSensor) DataRateMax() (int, error) {
	return s.ifk.rawIFK.GetDataRateMax(s.Index)
}

func (s *InterfaceKitSensor) DataRateMin() (int, error) {
	return s.ifk.rawIFK.GetDataRateMin(s.Index)
}

func (s *InterfaceKitSensor) RawValue() (int, error) {
	return s.ifk.rawIFK.GetSensorRawValue(s.Index)
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
		return raw.True
	}
	return raw.False
}

func newInterfaceKitInput(n int, i *InterfaceKit) *InterfaceKitInput {
	o := new(InterfaceKitInput)

	o.changed = make(chan int, channelSize)
	o.Changed = o.changed
	o.Index = n
	o.ifk = i

	return o
}

func newInterfaceKitOutput(n int, i *InterfaceKit) *InterfaceKitOutput {
	o := new(InterfaceKitOutput)

	o.changed = make(chan int, channelSize)
	o.Changed = o.changed
	o.Index = n
	o.ifk = i

	return o
}

func newInterfaceKitSensor(n int, i *InterfaceKit) *InterfaceKitSensor {
	s := new(InterfaceKitSensor)

	s.changed = make(chan int, channelSize)
	s.Changed = s.changed
	s.Index = n
	s.ifk = i

	return s
}

func stateToBool(s int) bool {
	if s == raw.True {
		return true
	}
	return false
}

func (i *InterfaceKit) createInputs() {
	if i.Inputs != nil {
		return
	}

	c, err := i.rawIFK.GetInputCount()
	if err != nil {
		return
	}

	i.Inputs = make([]*InterfaceKitInput, c)
	for n := 0; n < c; n++ {
		i.Inputs[n] = newInterfaceKitInput(n, i)
	}

	go func() {
		for c := range i.rawIFK.InputChanged {
			select {
			case i.Inputs[c.Index].changed <- c.Value:
			default:
			}
		}
	}()
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
			select {
			case i.Outputs[c.Index].changed <- c.Value:
			default:
			}
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
			select {
			case i.Sensors[c.Index].changed <- c.Value:
			default:
			}
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
