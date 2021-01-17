package frames

import (
	"fmt"
)

type SetFrameTypeID byte

// Types to use for SetFrame messages
const (
	HistoryRequest        SetFrameTypeID = 0x00
	SetTimeRequest        SetFrameTypeID = 0x01
	SetConfigRequest      SetFrameTypeID = 0x02
	ConfigRequest         SetFrameTypeID = 0x03
	CurrentWeatherRequest SetFrameTypeID = 0x04
	SendConfig            SetFrameTypeID = 0x20
	SendTime              SetFrameTypeID = 0x60
)

func (t SetFrameTypeID) String() string {
	var s string

	switch t {
	case HistoryRequest:
		s = "HistoryRequest"
	case SetTimeRequest:
		s = "SetTimeRequest"
	case SetConfigRequest:
		s = "SetConfigRequest"
	case ConfigRequest:
		s = "ConfigRequest"
	case CurrentWeatherRequest:
		s = "CurrentWeatherRequest"
	case SendConfig:
		s = "SendConfig"
	case SendTime:
		s = "SendTime"
	}

	return fmt.Sprintf("%s (0x%02x)", s, byte(t))
}

type SetFrame []byte

func NewSetFrame() SetFrame {
	f := make([]byte, 7, 273)

	f[0] = byte(0xd5) //klimalogg.SetFrame
	f[2] = 4

	return f
}

func (f SetFrame) Length() uint8 {
	return uint8(f[2])
}

func (f SetFrame) SetDeviceID(deviceID uint16) {
	f[3] = byte(deviceID >> 8)
	f[4] = byte(deviceID)
}

func (f SetFrame) DeviceID() uint16 {
	return uint16(f[3])<<8 | uint16(f[4])
}

func (f SetFrame) SetLoggerID(loggerID uint8) {
	f[5] = byte(loggerID)
}

func (f SetFrame) LoggerID() uint8 {
	return uint8(f[5])
}

func (f SetFrame) SetTypeID(typeID SetFrameTypeID) {
	f[6] = byte(typeID)
}

func (f SetFrame) TypeID() SetFrameTypeID {
	return SetFrameTypeID(f[6])
}

func (f *SetFrame) Write(data []byte) (int, error) {
	slice := *f
	if len(slice)+len(data) > cap(slice) {
		return 0, fmt.Errorf("data exceeds frame capacity")
	}

	slice = append(slice, data...)
	slice[2] = uint8(slice[2]) + uint8(len(data))

	*f = slice
	return len(data), nil
}

func (f SetFrame) String() string {
	return fmt.Sprintf("SetFrame{Length:%d DeviceID:0x%04x LoggerID:%02d TypeID:%s}", f.Length(), f.DeviceID(), f.LoggerID(), f.TypeID())
}

func (f SetFrame) MarshalBinary() ([]byte, error) {
	return f, nil
}
