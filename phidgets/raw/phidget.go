package raw

// #cgo darwin CFLAGS: -I/Library/Frameworks/Phidget21.framework/Headers
// #cgo darwin LDFLAGS: -framework Phidget21
// #cgo linux CFLAGS: -lphidget21
// #cgo linux LDFLAGS: -lphidget21
// #include "phidget.h"
import "C"

import (
	"errors"
	"fmt"
	"time"
)

const (
	Any = -1
	False = C.PFALSE
	True = C.PTRUE
)

type Status int

const (
	Attached    = C.PHIDGET_ATTACHED
	NotAttached = C.PHIDGET_NOTATTACHED
)

type Class int

const (
	AccelerometerClass     = C.PHIDCLASS_ACCELEROMETER
	AdvancedServoClass     = C.PHIDCLASS_ADVANCEDSERVO
	AnalogClass            = C.PHIDCLASS_ANALOG
	BridgeClass            = C.PHIDCLASS_BRIDGE
	EncoderClass           = C.PHIDCLASS_ENCODER
	FrequencyCounterClass  = C.PHIDCLASS_FREQUENCYCOUNTER
	GPSClass               = C.PHIDCLASS_GPS
	InterfaceKitClass      = C.PHIDCLASS_INTERFACEKIT
	IRClass                = C.PHIDCLASS_IR
	LEDClass               = C.PHIDCLASS_LED
	MotorControlClass      = C.PHIDCLASS_MOTORCONTROL
	PHSensorClass          = C.PHIDCLASS_PHSENSOR
	RFIDClass              = C.PHIDCLASS_RFID
	ServoClass             = C.PHIDCLASS_SERVO
	SpatialClass           = C.PHIDCLASS_SPATIAL
	StepperClass           = C.PHIDCLASS_STEPPER
	TemperatureSensorClass = C.PHIDCLASS_TEMPERATURESENSOR
	TextLCDClass           = C.PHIDCLASS_TEXTLCD
	TextLEDClass           = C.PHIDCLASS_TEXTLED
	WeightSensorClass      = C.PHIDCLASS_WEIGHTSENSOR
)

type ID int

const (
	// Device        = C.PHIDID_INTERFACEKIT_4_8_8 // No idea
	// Device        = C.PHIDID_RFID               // No idea
	// Device        = C.PHIDID_TEXTLED_1x8        // No idea
	// Device        = C.PHIDID_TEXTLED_4x8        // No idea
	// Device1000    = C.PHIDID_SERVO_1MOTOR
	// Device1000Old = C.PHIDID_SERVO_1MOTOR_OLD
	// Device1001    = C.PHIDID_SERVO_4MOTOR
	// Device1001Old = C.PHIDID_SERVO_4MOTOR_OLD
	// Device1002    = C.PHIDID_ANALOG_4OUTPUT
	// Device1011    = C.PHIDID_INTERFACEKIT_2_2_2
	// Device1012    = C.PHIDID_INTERFACEKIT_0_16_16
	// Device1014    = C.PHIDID_INTERFACEKIT_0_0_4
	// Device1015    = C.PHIDID_LINEAR_TOUCH
	// Device1016    = C.PHIDID_ROTARY_TOUCH
	// Device1017    = C.PHIDID_INTERFACEKIT_0_0_8
	// Device1018    = C.PHIDID_INTERFACEKIT_8_8_8
	// Device1023    = C.PHIDID_RFID_2OUTPUT
	// Device1024    = C.PHIDID_RFID_2OUTPUT_READ_WRITE // Not Found
	// Device1030    = C.PHIDID_LED_64
	// Device1031    = C.PHIDID_LED_64_ADV
	// Device1040    = C.PHIDID_GPS
	// Device1043    = C.PHIDID_SPATIAL_ACCEL_3AXIS
	// Device1044    = C.PHIDID_SPATIAL_ACCEL_GYRO_COMPASS
	// Device1045    = C.PHIDID_TEMPERATURESENSOR_IR
	// Device1046    = C.PHIDID_BRIDGE_4INPUT
	// Device1047    = C.PHIDID_ENCODER_HS_4ENCODER_4INPUT
	// Device1048    = C.PHIDID_TEMPERATURESENSOR_4
	// Device1050    = C.PHIDID_WEIGHTSENSOR // Not Found
	// Device1051    = C.PHIDID_TEMPERATURESENSOR
	// Device1052    = C.PHIDID_ENCODER_1ENCODER_1INPUT
	// Device1054    = C.PHIDID_ACCELEROMETER_2AXIS // PhidgetFrequencyCounter
	// Device1054    = C.PHIDID_FREQUENCYCOUNTER_2INPUT
	IRID = C.PHIDID_IR
	// Device1057    = C.PHIDID_ENCODER_HS_1ENCODER
	// Device1058    = C.PHIDID_PHSENSOR
	// Device1059    = C.PHIDID_ACCELEROMETER_3AXIS
	// Device1060    = C.PHIDID_MOTORCONTROL_LV_2MOTOR_4INPUT
	// Device1061    = C.PHIDID_ADVANCEDSERVO_8MOTOR
	// Device1062    = C.PHIDID_UNIPOLAR_STEPPER_4MOTOR
	// Device1063    = C.PHIDID_BIPOLAR_STEPPER_1MOTOR
	// Device1064    = C.PHIDID_MOTORCONTROL_HC_2MOTOR
	// Device1065    = C.PHIDID_MOTORCONTROL_1MOTOR
	// Device1066    = C.PHIDID_ADVANCEDSERVO_1MOTOR
	// Device1203    = C.PHIDID_INTERFACEKIT_8_8_8_w_LCD
	// Device1203    = C.PHIDID_TEXTLCD_2x20_w_8_8_8
	// Device1204    = C.PHIDID_TEXTLCD_ADAPTER
	// Device1210    = C.PHIDID_TEXTLCD_2x20             // Not Found
	// Device1221    = C.PHIDID_INTERFACEKIT_0_8_8_w_LCD // Not Found
	// Device1221    = C.PHIDID_TEXTLCD_2x20_w_0_8_8     // Not Found
)

type eventType int

const (
	phidgetAttach     = C.phidgetAttach
	phidgetConnect    = C.phidgetConnect
	phidgetDetach     = C.phidgetDetach
	phidgetDisconnect = C.phidgetDisconnect
)

type Phidget struct {
	Attached     <-chan bool
	Connected    <-chan bool
	Detached     <-chan bool
	Disconnected <-chan bool
	Error        <-chan error

	attached            chan bool
	connected           chan bool
	detached            chan bool
	disconnected        chan bool
	error               chan error
	handle              C.CPhidgetHandle
	onAttachHandler     *C.handler
	onConnectHandler    *C.handler
	onDetachHandler     *C.handler
	onDisconnectHandler *C.handler
	onErrorHandler      *C.handler
}

func ErrorDescription(code int) string {
	str := new(*C.char)
	C.CPhidget_getErrorDescription(C.int(code), str)
	return C.GoString(*str)
}

func LibraryVersion() (string, error) {
	str := new(*C.char)
	err := result(C.CPhidget_getLibraryVersion(str))
	if err != nil {
		return "", err
	}
	return C.GoString(*str), nil
}

func (p *Phidget) Close() error {
	return result(C.CPhidget_close(p.handle))
}

func (p *Phidget) GetDeviceClass() (Class, error) {
	ptr := new(C.CPhidget_DeviceClass)
	err := result(C.CPhidget_getDeviceClass(p.handle, ptr))
	if err != nil {
		return 0, err
	}
	return Class(*ptr), nil
}

func (p *Phidget) GetDeviceID() (ID, error) {
	ptr := new(C.CPhidget_DeviceID)
	err := result(C.CPhidget_getDeviceID(p.handle, ptr))
	if err != nil {
		return 0, err
	}
	return ID(*ptr), nil
}

func (p *Phidget) GetDeviceLabel() (string, error) {
	return resultWithString(func(ptr **C.char) C.int { return C.CPhidget_getDeviceLabel(p.handle, ptr) })
}

func (p *Phidget) GetDeviceName() (string, error) {
	return resultWithString(func(ptr **C.char) C.int { return C.CPhidget_getDeviceName(p.handle, ptr) })
}

func (p *Phidget) GetDeviceStatus() (int, error) {
	return resultWithInt(func(ptr *C.int) C.int { return C.CPhidget_getDeviceStatus(p.handle, ptr) })
}

func (p *Phidget) GetDeviceType() (string, error) {
	return resultWithString(func(ptr **C.char) C.int { return C.CPhidget_getDeviceType(p.handle, ptr) })
}

func (p *Phidget) GetDeviceVersion() (int, error) {
	return resultWithInt(func(ptr *C.int) C.int { return C.CPhidget_getDeviceVersion(p.handle, ptr) })
}

func (p *Phidget) GetSerialNumber() (int, error) {
	return resultWithInt(func(ptr *C.int) C.int { return C.CPhidget_getSerialNumber(p.handle, ptr) })
}

func (p *Phidget) GetServerAddress() (string, int, error) {
	addr := new(*C.char)
	port := new(C.int)
	err := result(C.CPhidget_getServerAddress(p.handle, addr, port))
	if err != nil {
		return "", 0, err
	}
	return C.GoString(*addr), int(*port), nil
}

func (p *Phidget) GetServerID() (string, error) {
	return resultWithString(func(ptr **C.char) C.int { return C.CPhidget_getServerID(p.handle, ptr) })
}

func (p *Phidget) GetServerStatus() (int, error) {
	return resultWithInt(func(ptr *C.int) C.int { return C.CPhidget_getServerStatus(p.handle, ptr) })
}

func (p *Phidget) Open(serial int) error {
	return result(C.CPhidget_open(p.handle, C.int(serial)))
}

func (p *Phidget) OpenLabel(label string) error {
	return result(C.CPhidget_openLabel(p.handle, convertString(label)))
}

func (p *Phidget) OpenLabelRemote(label, server, password string) error {
	return result(C.CPhidget_openLabelRemote(p.handle, convertString(label), convertString(server), convertString(password)))
}

func (p *Phidget) OpenLabelRemoteIP(label string, address string, port int, password string) error {
	return result(C.CPhidget_openLabelRemoteIP(p.handle, convertString(label), convertString(address), C.int(port), convertString(password)))
}

func (p *Phidget) OpenRemote(serial int, server, password string) error {
	return result(C.CPhidget_openRemote(p.handle, C.int(serial), convertString(server), convertString(password)))
}

func (p *Phidget) OpenRemoteIP(serial int, address string, port int, password string) error {
	return result(C.CPhidget_openRemoteIP(p.handle, C.int(serial), convertString(address), C.int(port), convertString(password)))
}

func (p *Phidget) SetDeviceLabel(label string) error {
	return result(C.CPhidget_setDeviceLabel(p.handle, convertString(label)))
}

func (p *Phidget) WaitForAttachment(timeout time.Duration) error {
	return result(C.CPhidget_waitForAttachment(p.handle, C.int(timeout/time.Millisecond)))
}

func convertString(str string) *C.char {
	if str == "" {
		return nil
	}
	return C.CString(str)
}

func result(result C.int) error {
	code := int(result)
	if code != 0 {
		return errors.New(fmt.Sprintf("%s (%d)", ErrorDescription(code), code))
	}
	return nil
}

func resultWithInt(f func(*C.int) C.int) (int, error) {
	ptr := new(C.int)
	err := result(f(ptr))
	if err != nil {
		return 0, err
	}
	return int(*ptr), nil
}

func resultWithString(f func(**C.char) C.int) (string, error) {
	ptr := new(*C.char)
	err := result(f(ptr))
	if err != nil {
		return "", err
	}
	return C.GoString(*ptr), nil
}

func (p *Phidget) cleanup() {
	p.unsetOnErrorHandler()
	p.unsetOnDisconnectHandler()
	p.unsetOnDetachHandler()
	p.unsetOnConnectHandler()
	p.unsetOnAttachHandler()
	C.CPhidget_delete(p.handle)
}

func (p *Phidget) initPhidget(h C.CPhidgetHandle) error {
	p.handle = h

	p.attached = make(chan bool)
	p.connected = make(chan bool)
	p.detached = make(chan bool)
	p.disconnected = make(chan bool)
	p.error = make(chan error)

	p.Attached = p.attached
	p.Connected = p.connected
	p.Detached = p.detached
	p.Disconnected = p.disconnected
	p.Error = p.error

	if err := p.setOnAttachHandler(p.attached); err != nil {
		return err
	}

	if err := p.setOnConnectHandler(p.connected); err != nil {
		return err
	}

	if err := p.setOnDetachHandler(p.detached); err != nil {
		return err
	}

	if err := p.setOnDisconnectHandler(p.detached); err != nil {
		return err
	}

	if err := p.setOnErrorHandler(p.error); err != nil {
		return err
	}

	return nil
}

func (p *Phidget) setOnAttachHandler(c chan bool) error {
	var err error

	p.onAttachHandler, err = p.setOnEventHandler(c, phidgetAttach)
	if err != nil {
		return err
	}

	return nil
}

func (p *Phidget) setOnConnectHandler(c chan bool) error {
	var err error

	p.onAttachHandler, err = p.setOnEventHandler(c, phidgetConnect)
	if err != nil {
		return err
	}

	return nil
}

func (p *Phidget) setOnDetachHandler(c chan bool) error {
	var err error

	p.onAttachHandler, err = p.setOnEventHandler(c, phidgetDetach)
	if err != nil {
		return err
	}

	return nil
}

func (p *Phidget) setOnDisconnectHandler(c chan bool) error {
	var err error

	p.onAttachHandler, err = p.setOnEventHandler(c, phidgetDisconnect)
	if err != nil {
		return err
	}

	return nil
}

func (p *Phidget) setOnErrorHandler(c chan error) error {
	var err error

	p.onErrorHandler, err = createHandler(func(h *C.handler) C.int {
		return C.setOnErrorHandler(p.handle, h)
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			r := C.onErrorAwait(p.onErrorHandler)
			e := errors.New(C.GoString(r.string))
			C.onErrorResultFree(r)
			c <- e
		}
	}()

	return nil
}

func (p *Phidget) setOnEventHandler(c chan bool, t eventType) (*C.handler, error) {
	h, err := createHandler(func(h *C.handler) C.int {
		return C.setOnEventHandler(p.handle, h, C.eventType(t))
	})
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			C.onEventAwait(h)
			c <- true
		}
	}()

	return h, nil
}

func (p *Phidget) unsetOnAttachHandler() {
	C.unsetOnEventHandler(p.handle, phidgetAttach)
	C.handlerFree(p.onAttachHandler)
	p.onAttachHandler = nil
}
func (p *Phidget) unsetOnConnectHandler() {
	C.unsetOnEventHandler(p.handle, phidgetConnect)
	C.handlerFree(p.onConnectHandler)
	p.onConnectHandler = nil
}
func (p *Phidget) unsetOnDetachHandler() {
	C.unsetOnEventHandler(p.handle, phidgetDetach)
	C.handlerFree(p.onDetachHandler)
	p.onDetachHandler = nil
}
func (p *Phidget) unsetOnDisconnectHandler() {
	C.unsetOnEventHandler(p.handle, phidgetDisconnect)
	C.handlerFree(p.onDisconnectHandler)
	p.onDisconnectHandler = nil
}

func (p *Phidget) unsetOnErrorHandler() {
	C.unsetOnErrorHandler(p.handle)
	C.handlerFree(p.onErrorHandler)
	p.onErrorHandler = nil
}

func (p *Phidget) unsetOnEventHandler(t eventType) {
	C.unsetOnEventHandler(p.handle, C.eventType(t))
}
