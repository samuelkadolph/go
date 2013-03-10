package phidgets

import (
	"github.com/samuelkadolph/go/phidgets/raw"
)

type IR struct {
	Phidget

	Code    <-chan IRCode
	Learn   <-chan IRLearn
	RawData <-chan IRRawData
	Error   <-chan error

	rawIR *raw.IR
}

type IRCode raw.IRCode

type IRCodeInfo raw.IRCodeInfo

type IRLearn struct {
	Data     []byte
	CodeInfo IRCodeInfo
}

type IRRawData raw.IRRawData

func NewIR() (*IR, error) {
	ir := new(IR)

	r, err := raw.NewIR()
	if err != nil {
		return nil, err
	}

	if err := ir.initIR(r); err != nil {
		return nil, err
	}

	return ir, nil
}

func (i *IR) GetLastCode() ([]byte, int, error) {
	return i.rawIR.GetLastCode()
}

func (i *IR) GetLastLearnedCode() ([]byte, IRCodeInfo, error) {
	d, ci, err := i.rawIR.GetLastLearnedCode()
	return d, IRCodeInfo(ci), err
}

func (i *IR) GetRawData(max int) ([]int, error) {
	return i.rawIR.GetRawData(max)
}

func (i *IR) Transmit(data []byte, info IRCodeInfo) error {
	return i.rawIR.Transmit(data, raw.IRCodeInfo(info))
}

func (i *IR) TransmitRaw(data []int, carrierFrequency, dutyCycle, gap int) error {
	return i.rawIR.TransmitRaw(data, carrierFrequency, dutyCycle, gap)
}

func (i *IR) TransmitRepeat() error {
	return i.rawIR.TransmitRepeat()
}

func (i *IR) initIR(r *raw.IR) error {
	i.rawIR = r

	if err := i.initPhidget(&r.Phidget); err != nil {
		return err
	}

	code := make(chan IRCode)
	learn := make(chan IRLearn)
	raw := make(chan IRRawData)

	i.Code = code
	i.Learn = learn
	i.RawData = raw
	i.Error = i.rawIR.Error

	go func() {
		for c := range i.rawIR.Code {
			code <- IRCode(c)
		}
	}()
	go func() {
		for l := range i.rawIR.Learn {
			learn <- IRLearn{l.Data, IRCodeInfo(l.CodeInfo)}
		}
	}()
	go func() {
		for r := range i.rawIR.RawData {
			raw <- IRRawData(r)
		}
	}()

	return nil
}
