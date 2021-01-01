package transceiver

import (
	"encoding"
	"fmt"
)

// TypeID of data in frame
type TypeID byte

// Types to use for SetFrame messages
const (
	HistoryRequest        TypeID = 0x00
	SetTimeRequest        TypeID = 0x01
	SetConfigRequest      TypeID = 0x02
	ConfigRequest         TypeID = 0x03
	CurrentWeatherRequest TypeID = 0x04
	SendConfig            TypeID = 0x20
	SendTime              TypeID = 0x60
)

// Types used by GetFrame messages
const (
	DataWrittenResponse        TypeID = 0x10
	ConfigResponse             TypeID = 0x20
	CurrentWeatherResponse     TypeID = 0x30
	HistoryResponse            TypeID = 0x40
	RequestHistoryResponse     TypeID = 0x50
	RequestFirstConfigResponse TypeID = 0x51
	RequestConfigResponse      TypeID = 0x52
	RequestTimeResponse        TypeID = 0x53
)

const frameByteLength = 273

// Frame is send to or recieved from klimalogg console
type Frame struct {
	MessageID MessageID
	Length    uint8
	DeviceID  uint16
	LoggerID  uint8
	TypeID    TypeID
	Data      []byte
}

func (f Frame) String() string {
	return fmt.Sprintf("{MessageID:%02x Length:%d DeviceID:%04x LoggerID:%d TypeID:%02x}", f.MessageID, f.Length, f.DeviceID, f.LoggerID, f.TypeID)
}

// UnmarshalBinary response from klimalogg console
func (f *Frame) UnmarshalBinary(data []byte) error {
	if len(data) != frameByteLength {
		return fmt.Errorf("unexpected length of data; expected %dB, got %dB", frameByteLength, len(data))
	}

	f.MessageID = MessageID(data[0])
	f.Length = data[2]
	f.DeviceID = uint16(data[3])<<8 | uint16(data[4])
	f.LoggerID = data[5]
	f.TypeID = TypeID(data[6])

	dataLength := int(f.Length) - 4
	if dataLength+7 > frameByteLength {
		return fmt.Errorf("invalid frame message length %dB, exceeding frame length %dB", dataLength, frameByteLength)
	}

	f.Data = data[7 : dataLength+8]

	return nil
}

// MarshalBinary request for klimalogg console
func (f Frame) MarshalBinary() ([]byte, error) {
	data := make([]byte, 273)

	data[0] = byte(f.MessageID)
	data[2] = uint8(len(f.Data) + 4)
	data[3] = byte(f.DeviceID >> 8)
	data[4] = byte(f.DeviceID)
	data[5] = byte(f.LoggerID)
	data[6] = byte(f.TypeID)

	for i, b := range f.Data {
		data[i+7] = b
	}

	return data, nil
}

// SetData of frame
func (f *Frame) SetData(data encoding.BinaryMarshaler) error {
	b, err := data.MarshalBinary()
	if err != nil {
		return err
	}

	f.Data = b
	return nil
}
