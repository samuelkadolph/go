package raw

// #include "interfacekit.h"
import "C"

import (
	"runtime"
)

type onChangeType int

const (
	inputChanged  = C.inputChanged
	outputChanged = C.outputChanged
	sensorChanged = C.sensorChanged
)

type InterfaceKit struct {
	Phidget

	InputChanged  <-chan InterfaceKitChange
	OutputChanged <-chan InterfaceKitChange
	SensorChanged <-chan InterfaceKitChange

	ifkHandle              C.CPhidgetInterfaceKitHandle
	inputChanged           chan InterfaceKitChange
	onInputChangedHandler  *C.handler
	onOutputChangedHandler *C.handler
	onSensorChangedHandler *C.handler
	outputChanged          chan InterfaceKitChange
	sensorChanged          chan InterfaceKitChange
}

type InterfaceKitChange struct {
	Index int
	Value int
}

func NewInterfaceKit() (*InterfaceKit, error) {
	h := new(C.CPhidgetInterfaceKitHandle)
	if err := result(C.CPhidgetInterfaceKit_create(h)); err != nil {
		return nil, err
	}

	ph := new(InterfaceKit)
	if err := ph.initInterfaceKit(*h); err != nil {
		return nil, err
	}

	return ph, nil
}

func (i *InterfaceKit) GetDataRate(index int) (int, error) {
	return resultWithInt(func(p *C.int) C.int { return C.CPhidgetInterfaceKit_getDataRate(i.ifkHandle, C.int(index), p) })
}

func (i *InterfaceKit) GetDataRateMax(index int) (int, error) {
	return resultWithInt(func(p *C.int) C.int { return C.CPhidgetInterfaceKit_getDataRateMax(i.ifkHandle, C.int(index), p) })
}

func (i *InterfaceKit) GetDataRateMin(index int) (int, error) {
	return resultWithInt(func(p *C.int) C.int { return C.CPhidgetInterfaceKit_getDataRateMin(i.ifkHandle, C.int(index), p) })
}

func (i *InterfaceKit) GetInputCount() (int, error) {
	return resultWithInt(func(p *C.int) C.int { return C.CPhidgetInterfaceKit_getInputCount(i.ifkHandle, p) })
}

func (i *InterfaceKit) GetInputState(index int) (int, error) {
	return resultWithInt(func(p *C.int) C.int { return C.CPhidgetInterfaceKit_getInputState(i.ifkHandle, C.int(index), p) })
}

func (i *InterfaceKit) GetOutputCount() (int, error) {
	return resultWithInt(func(p *C.int) C.int { return C.CPhidgetInterfaceKit_getOutputCount(i.ifkHandle, p) })
}

func (i *InterfaceKit) GetOutputState(index int) (int, error) {
	return resultWithInt(func(p *C.int) C.int { return C.CPhidgetInterfaceKit_getOutputState(i.ifkHandle, C.int(index), p) })
}

func (i *InterfaceKit) GetRatiometric() (int, error) {
	return resultWithInt(func(p *C.int) C.int { return C.CPhidgetInterfaceKit_getRatiometric(i.ifkHandle, p) })
}

func (i *InterfaceKit) GetSensorChangeTrigger(index int) (int, error) {
	return resultWithInt(func(p *C.int) C.int {
		return C.CPhidgetInterfaceKit_getSensorChangeTrigger(i.ifkHandle, C.int(index), p)
	})
}

func (i *InterfaceKit) GetSensorCount() (int, error) {
	return resultWithInt(func(p *C.int) C.int { return C.CPhidgetInterfaceKit_getSensorCount(i.ifkHandle, p) })
}

func (i *InterfaceKit) GetSensorRawValue(index int) (int, error) {
	return resultWithInt(func(p *C.int) C.int { return C.CPhidgetInterfaceKit_getSensorRawValue(i.ifkHandle, C.int(index), p) })
}

func (i *InterfaceKit) GetSensorValue(index int) (int, error) {
	return resultWithInt(func(p *C.int) C.int { return C.CPhidgetInterfaceKit_getSensorValue(i.ifkHandle, C.int(index), p) })
}

func (i *InterfaceKit) SetDataRate(index, rate int) error {
	return result(C.CPhidgetInterfaceKit_setDataRate(i.ifkHandle, C.int(index), C.int(rate)))
}

func (i *InterfaceKit) SetOutputState(index, state int) error {
	return result(C.CPhidgetInterfaceKit_setOutputState(i.ifkHandle, C.int(index), C.int(state)))
}

func (i *InterfaceKit) SetRatiometric(ratiometric int) error {
	return result(C.CPhidgetInterfaceKit_setRatiometric(i.ifkHandle, C.int(ratiometric)))
}

func (i *InterfaceKit) SetSensorChangeTrigger(index, trigger int) error {
	return result(C.CPhidgetInterfaceKit_setSensorChangeTrigger(i.ifkHandle, C.int(index), C.int(trigger)))
}

func (i *InterfaceKit) cleanupInterfaceKit() {
	i.unsetOnChangeHandler(sensorChanged, &i.onSensorChangedHandler)
	i.unsetOnChangeHandler(outputChanged, &i.onOutputChangedHandler)
	i.unsetOnChangeHandler(inputChanged, &i.onInputChangedHandler)
	i.cleanup()
}

func (i *InterfaceKit) initInterfaceKit(h C.CPhidgetInterfaceKitHandle) error {
	var err error

	runtime.SetFinalizer(i, func(i *InterfaceKit) { i.cleanupInterfaceKit() })

	i.ifkHandle = h

	if err = i.initPhidget(C.CPhidgetHandle(h)); err != nil {
		return err
	}

	i.inputChanged = make(chan InterfaceKitChange, channelSize)
	i.outputChanged = make(chan InterfaceKitChange, channelSize)
	i.sensorChanged = make(chan InterfaceKitChange, channelSize)

	i.InputChanged = i.inputChanged
	i.OutputChanged = i.outputChanged
	i.SensorChanged = i.sensorChanged

	if i.onInputChangedHandler, err = i.setOnChangeHandler(i.inputChanged, inputChanged); err != nil {
		return err
	}

	if i.onOutputChangedHandler, err = i.setOnChangeHandler(i.outputChanged, outputChanged); err != nil {
		return err
	}

	if i.onSensorChangedHandler, err = i.setOnChangeHandler(i.sensorChanged, sensorChanged); err != nil {
		return err
	}

	return nil
}

func (i *InterfaceKit) setOnChangeHandler(c chan InterfaceKitChange, t onChangeType) (*C.handler, error) {
	h, err := createHandler(func(h *C.handler) C.int {
		return C.setOnChangeHandler(i.ifkHandle, h, C.onChangeType(t))
	})
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			r := C.onChangeAwait(h)

			change := InterfaceKitChange{}
			change.Index = int(r.index)
			change.Value = int(r.value)

			C.onChangeResultFree(r)

			select {
			case c <- change:
			default:
			}
		}
	}()

	return h, nil
}

func (i *InterfaceKit) unsetOnChangeHandler(t onChangeType, h **C.handler) {
	C.unsetOnChangeHandler(i.ifkHandle, C.onChangeType(t))
	C.handlerFree(*h)
	*h = nil
}
