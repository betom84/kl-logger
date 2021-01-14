package klimalogg

import (
	"encoding"
	"fmt"
	"time"

	"github.com/betom84/kl-logger/utils"
	"github.com/google/gousb"
	"github.com/sirupsen/logrus"
)

var (
	vendorID  gousb.ID = 0x6666
	productID gousb.ID = 0x5555
)

// MessageID ...
type MessageID byte

// MessageID to control transceiver
const (
	GetFrame           MessageID = 0x00
	SetRX              MessageID = 0xd0
	SetTX              MessageID = 0xd1
	SetFrame           MessageID = 0xd5
	SetState           MessageID = 0xd7
	SetPreamblePattern MessageID = 0xd8
	Execute            MessageID = 0xd9
	GetState           MessageID = 0xde
	ReadConfigFlashIn  MessageID = 0xdc
	ReadConfigFlashOut MessageID = 0xdd
	WriteRegister      MessageID = 0xf0
)

func (m MessageID) String() string {
	var name string

	switch m {
	case GetFrame:
		name = "GetFrame"
	case SetRX:
		name = "SetRX"
	case SetTX:
		name = "SetTX"
	case SetFrame:
		name = "SetFrame"
	case SetState:
		name = "SetState"
	case SetPreamblePattern:
		name = "SetPreamblePattern"
	case Execute:
		name = "Execute"
	case GetState:
		name = "GetState"
	case ReadConfigFlashIn:
		name = "ReadConfigFlashIn"
	case ReadConfigFlashOut:
		name = "ReadConfigFlashOut"
	case WriteRegister:
		name = "WriteRegister"
	}

	return fmt.Sprintf("%s (0x%02x)", name, byte(m))
}

type message struct {
	rtype   uint8
	request uint8
	value   uint16
	index   uint16
	length  uint
}

var messages = map[MessageID]message{
	GetFrame:           {gousb.ControlClass | gousb.ControlInterface | gousb.ControlIn, 0x01, 0x00003d6, 0x0000000, 0x111},
	SetRX:              {gousb.ControlClass | gousb.ControlInterface, 0x0000009, 0x00003d0, 0x0000000, 0x15},
	SetTX:              {gousb.ControlClass | gousb.ControlInterface, 0x0000009, 0x00003d1, 0x0000000, 0x15},
	SetFrame:           {gousb.ControlClass | gousb.ControlInterface, 0x0000009, 0x00003d5, 0x0000000, 0x111},
	SetState:           {gousb.ControlClass | gousb.ControlInterface, 0x0000009, 0x00003d7, 0x0000000, 0x15},
	SetPreamblePattern: {gousb.ControlClass | gousb.ControlInterface, 0x0000009, 0x00003d8, 0x0000000, 0x15},
	Execute:            {gousb.ControlClass | gousb.ControlInterface, 0x0000009, 0x00003d9, 0x0000000, 0x0f},
	GetState:           {gousb.ControlClass | gousb.ControlInterface | gousb.ControlIn, 0x01, 0x00003de, 0x0000000, 0x06},
	ReadConfigFlashIn:  {gousb.ControlClass | gousb.ControlInterface, 0x0000009, 0x00003dd, 0x0000000, 0x0f},
	ReadConfigFlashOut: {gousb.ControlClass | gousb.ControlInterface | gousb.ControlIn, 0x01, 0x00003dc, 0x0000000, 0x0f},
	WriteRegister:      {gousb.ControlClass | gousb.ControlInterface, 0x0000009, 0x00003f0, 0x0000000, 0x05},
}

// FlashProperty to identify flash memory register
type FlashProperty struct {
	address uint16
	length  uint8
}

// Supported configuration values stored in transceiver flash memory
var (
	FrequencyCorrection FlashProperty = FlashProperty{0x1F5, 4}
	TransceiverInfo     FlashProperty = FlashProperty{0x1F9, 7}
)

// Transceiver usb device
type Transceiver struct {
	VendorID  uint16
	ProductID uint16

	context              *gousb.Context
	device               *gousb.Device
	defaultInterface     *gousb.Interface
	defaultInterfaceDone func()
}

// Get message via usb control
func (t Transceiver) Get(id MessageID) ([]byte, error) {
	m := messages[id]
	b := make([]byte, int(m.length))

	len, err := t.device.Control(m.rtype, m.request, m.value, m.index, b)

	if id != GetState {
		logrus.WithFields(logrus.Fields{
			"id":   fmt.Sprintf("%02x", byte(id)),
			"len":  len,
			"data": utils.Prettify(b),
		}).Tracef("<- %s", id)
	}

	return b, err
}

// Set message via usb control
func (t Transceiver) Set(id MessageID, data []byte) error {
	m := messages[id]
	b := make([]byte, int(m.length))
	b[0] = byte(id)

	for idx, d := range data {
		b[idx+1] = d
	}

	len, err := t.device.Control(m.rtype, m.request, m.value, m.index, b)

	logrus.WithFields(logrus.Fields{
		"id":   fmt.Sprintf("%02x", byte(id)),
		"len":  len,
		"data": utils.Prettify(b),
	}).Tracef("-> %s", id)

	return err
}

// GetFrame message via usb control
func (t Transceiver) GetFrame(unmarshaler encoding.BinaryUnmarshaler) error {
	b, err := t.Get(GetFrame)
	if err != nil {
		return err
	}

	return unmarshaler.UnmarshalBinary(b)
}

// SetFrame message via usb control
func (t Transceiver) SetFrame(marshaler encoding.BinaryMarshaler) error {
	b, err := marshaler.MarshalBinary()
	if err != nil {
		return err
	}

	return t.Set(SetFrame, b[1:])
}

// IsReady returns true previous message was processed (and response is available)
func (t Transceiver) IsReady() bool {
	b, err := t.Get(GetState)
	if err != nil {
		logrus.WithError(err).Error("transceiver readiness probe failed")

		if err.Error() == "libusb: no device [code -4]" {
			panic(err)
		}

		return false
	}

	return b[1] == 0x16
}

func (t Transceiver) String() string {
	return t.device.String()
}

// ReadConfigFlash value from transceiver flash memory
func (t Transceiver) ReadConfigFlash(value FlashProperty) ([]byte, error) {

	// todo, init buffer with 0xcc (?)

	err := t.Set(ReadConfigFlashIn, []byte{0x0a, byte((value.address >> 8) & 0xff), byte((value.address) & 0xff)})
	if err != nil {
		return []byte{}, err
	}

	out, err := t.Get(ReadConfigFlashOut)
	if err != nil {
		return []byte{}, err
	}

	b := make([]byte, int(value.length))
	for i := range b {
		b[i] = out[i+4]
	}

	return b, err
}

// Open transceiver usb device and claim default interface
func (t *Transceiver) Open() error {
	var err error

	if t.context != nil || t.device != nil || t.defaultInterface != nil {
		return fmt.Errorf("transceiver device already opened; close first to reopen")
	}

	t.context = gousb.NewContext()
	t.device, err = t.context.OpenDeviceWithVIDPID(gousb.ID(t.VendorID), gousb.ID(t.ProductID))
	if err != nil {
		return fmt.Errorf("could not open transceiver; %v", err)
	}

	err = t.device.SetAutoDetach(true)
	if err != nil {
		return fmt.Errorf("failed to enable auto kernel driver detachment; %v", err)
	}

	t.device.ControlTimeout = 1 * time.Second

	t.defaultInterface, t.defaultInterfaceDone, err = t.device.DefaultInterface()

	return err
}

// Close transceiver usb device
func (t *Transceiver) Close() error {
	if t.defaultInterfaceDone != nil {
		t.defaultInterfaceDone()

		t.defaultInterface = nil
		t.defaultInterfaceDone = nil
	}

	if t.device != nil {
		err := t.device.Close()
		if err != nil {
			return err
		}

		t.device = nil
	}

	if t.context != nil {
		err := t.context.Close()
		if err != nil {
			return err
		}

		t.context = nil
	}

	return nil
}

// NewTransceiver usb device identified by product and vendor ID
func NewTransceiver(vendorID uint16, productID uint16) *Transceiver {
	return &Transceiver{
		VendorID:  vendorID,
		ProductID: productID,
	}
}
