package raw

// #include "ir.h"
import "C"

import (
	"runtime"
	"unsafe"
)

type IR struct {
	Phidget

	Code    <-chan IRCode
	Learn   <-chan IRLearn
	RawData <-chan IRRawData

	code             chan IRCode
	irHandle         C.CPhidgetIRHandle
	learn            chan IRLearn
	onCodeHandler    *C.handler
	onLearnHandler   *C.handler
	onRawDataHandler *C.handler
	rawData          chan IRRawData
}

type IRCode struct {
	Data     []byte
	BitCount int
	Repeat   int
}

type IRCodeInfo struct {
	BitCount         int
	Encoding         int
	Length           int
	Gap              int
	Trail            int
	Header           [2]int
	One              [2]int
	Zero             [2]int
	Repeat           [26]int
	MinRepeat        int
	ToggleMask       [16]byte
	CarrierFrequency int
	DutyCycle        int
}

type IRLearn struct {
	Data     []byte
	CodeInfo IRCodeInfo
}

type IRRawData struct {
	Data []int
}

func NewIR() (*IR, error) {
	ptr := new(C.CPhidgetIRHandle)
	if err := result(C.CPhidgetIR_create(ptr)); err != nil {
		return nil, err
	}

	ph := new(IR)
	if err := ph.initIR(*ptr); err != nil {
		return nil, err
	}

	return ph, nil
}

func (i *IR) GetLastCode() ([]byte, int, error) {
	d := make([]byte, 16)
	n := new(C.int)
	c := new(C.int)

	*n = C.int(len(d))

	if err := result(C.CPhidgetIR_getLastCode(i.irHandle, (*C.uchar)(&d[0]), n, c)); err != nil {
		return nil, 0, err
	}

	return d[:int(*n)], int(*c), nil
}

func (i *IR) GetLastLearnedCode() ([]byte, IRCodeInfo, error) {
	d := make([]byte, 16)
	n := new(C.int)
	ci := new(C.CPhidgetIR_CodeInfo)

	*n = C.int(len(d))

	if err := result(C.CPhidgetIR_getLastLearnedCode(i.irHandle, (*C.uchar)(&d[0]), n, ci)); err != nil {
		return nil, IRCodeInfo{}, err
	}

	return d[:int(*n)], irCodeInfoFromC(C.CPhidgetIR_CodeInfoHandle(unsafe.Pointer(ci))), nil
}

func (i *IR) GetRawData(max int) ([]int, error) {
	d := make([]C.int, max)
	n := new(C.int)

	*n = C.int(max)

	if err := result(C.CPhidgetIR_getRawData(i.irHandle, (*C.int)(&d[0]), n)); err != nil {
		return nil, err
	}

	r := make([]int, int(*n))
	for i := 0; i < int(*n); i++ {
		r[i] = int(d[i])
	}

	return r, nil
}

func (i *IR) Transmit(data []byte, info IRCodeInfo) error {
	return result(C.CPhidgetIR_Transmit(i.irHandle, (*C.uchar)(&data[0]), info.toC()))
}

func (i *IR) TransmitRaw(data []int, carrierFrequency, dutyCycle, gap int) error {
	d := make([]C.int, len(data))
	for i, v := range data {
		d[i] = C.int(v)
	}

	return result(C.CPhidgetIR_TransmitRaw(i.irHandle, (*C.int)(&d[0]), C.int(len(data)), C.int(carrierFrequency), C.int(dutyCycle), C.int(gap)))
}

func (i *IR) TransmitRepeat() error {
	return result(C.CPhidgetIR_TransmitRepeat(i.irHandle))
}

func irCodeInfoFromC(ci C.CPhidgetIR_CodeInfoHandle) IRCodeInfo {
	info := IRCodeInfo{}

	info.BitCount = int(ci.bitCount)
	info.Encoding = int(ci.encoding)
	info.Length = int(ci.length)
	info.Gap = int(ci.gap)
	info.Trail = int(ci.trail)
	info.Header = [...]int{int(ci.header[0]), int(ci.header[1])}
	info.One = [...]int{int(ci.one[0]), int(ci.one[1])}
	info.Zero = [...]int{int(ci.zero[0]), int(ci.zero[1])}
	info.Repeat = [...]int{int(ci.repeat[0]), int(ci.repeat[1]), int(ci.repeat[2]), int(ci.repeat[3]), int(ci.repeat[4]), int(ci.repeat[5]), int(ci.repeat[6]), int(ci.repeat[7]), int(ci.repeat[8]), int(ci.repeat[9]), int(ci.repeat[10]), int(ci.repeat[11]), int(ci.repeat[12]), int(ci.repeat[13]), int(ci.repeat[14]), int(ci.repeat[15]), int(ci.repeat[16]), int(ci.repeat[17]), int(ci.repeat[18]), int(ci.repeat[19]), int(ci.repeat[20]), int(ci.repeat[21]), int(ci.repeat[22]), int(ci.repeat[23]), int(ci.repeat[24]), int(ci.repeat[25])}
	info.MinRepeat = int(ci.min_repeat)
	info.ToggleMask = [...]byte{byte(ci.toggle_mask[0]), byte(ci.toggle_mask[1]), byte(ci.toggle_mask[2]), byte(ci.toggle_mask[3]), byte(ci.toggle_mask[4]), byte(ci.toggle_mask[5]), byte(ci.toggle_mask[6]), byte(ci.toggle_mask[7]), byte(ci.toggle_mask[8]), byte(ci.toggle_mask[9]), byte(ci.toggle_mask[10]), byte(ci.toggle_mask[11]), byte(ci.toggle_mask[12]), byte(ci.toggle_mask[13]), byte(ci.toggle_mask[14]), byte(ci.toggle_mask[15])}
	info.CarrierFrequency = int(ci.carrierFrequency)
	info.DutyCycle = int(ci.dutyCycle)

	return info
}

func (i *IR) cleanupIR() {
	i.unsetOnRawDataHandler()
	i.unsetOnLearnHandler()
	i.unsetOnCodeHandler()
	i.cleanup()
}

func (i *IR) initIR(h C.CPhidgetIRHandle) error {
	runtime.SetFinalizer(i, func(i *IR) { i.cleanupIR() })

	i.irHandle = h

	if err := i.initPhidget(C.CPhidgetHandle(h)); err != nil {
		return nil
	}

	i.code = make(chan IRCode, channelSize)
	i.learn = make(chan IRLearn, channelSize)
	i.rawData = make(chan IRRawData, channelSize)

	i.Code = i.code
	i.Learn = i.learn
	i.RawData = i.rawData

	if err := i.setOnCodeHandler(); err != nil {
		return err
	}

	if err := i.setOnLearnHandler(); err != nil {
		return err
	}

	if err := i.setOnRawDataHandler(); err != nil {
		return err
	}

	return nil
}

func (i *IR) setOnCodeHandler() error {
	var err error

	i.onLearnHandler, err = createHandler(func(h *C.handler) C.int {
		return C.setOnCodeHandler(i.irHandle, h)
	})
	if err != nil {
		return nil
	}

	go func() {
		for {
			r := C.onCodeAwait(i.onLearnHandler)

			code := IRCode{}
			code.Data = C.GoBytes(unsafe.Pointer(r.data), r.dataLength)
			code.BitCount = int(r.bitCount)
			code.Repeat = int(r.repeat)

			C.onCodeResultFree(r)

			select {
			case i.code <- code:
			default:
			}
		}
	}()

	return nil
}

func (i *IR) setOnLearnHandler() error {
	var err error

	i.onCodeHandler, err = createHandler(func(h *C.handler) C.int {
		return C.setOnLearnHandler(i.irHandle, h)
	})
	if err != nil {
		return nil
	}

	go func() {
		for {
			r := C.onLearnAwait(i.onCodeHandler)

			learn := IRLearn{}
			learn.Data = C.GoBytes(unsafe.Pointer(r.data), r.dataLength)
			learn.CodeInfo = irCodeInfoFromC(r.codeInfo)

			C.onLearnResultFree(r)

			select {
			case i.learn <- learn:
			default:
			}
		}
	}()

	return nil
}

func (i *IR) setOnRawDataHandler() error {
	var err error

	i.onRawDataHandler, err = createHandler(func(h *C.handler) C.int {
		return C.setOnRawDataHandler(i.irHandle, h)
	})
	if err != nil {
		return nil
	}

	go func() {
		for {
			r := C.onRawDataAwait(i.onRawDataHandler)

			raw := IRRawData{}
			raw.Data = []int{}

			d := (*[2048]int)(unsafe.Pointer(r.data))[:r.dataLength]
			for _, i := range d {
				raw.Data = append(raw.Data, i)
			}

			C.onRawDataResultFree(r)

			select {
			case i.rawData <- raw:
			default:
			}
		}
	}()

	return nil
}

func (i *IRCodeInfo) toC() C.CPhidgetIR_CodeInfoHandle {
	ci := new(C.struct__CPhidgetIR_CodeInfo)

	ci.bitCount = C.int(i.BitCount)

	return ci
}

func (i *IR) unsetOnCodeHandler() {
	C.unsetOnCodeHandler(i.handle)
	C.handlerFree(i.onCodeHandler)
	i.onCodeHandler = nil
}

func (i *IR) unsetOnLearnHandler() {
	C.unsetOnCodeHandler(i.handle)
	C.handlerFree(i.onLearnHandler)
	i.onLearnHandler = nil
}

func (i *IR) unsetOnRawDataHandler() {
	C.unsetOnCodeHandler(i.handle)
	C.handlerFree(i.onRawDataHandler)
	i.onRawDataHandler = nil
}
