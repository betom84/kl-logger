package frames

import (
	"fmt"
)

type GetFrameTypeID byte

// Types used by GetFrame messages
const (
	DataWrittenResponse        GetFrameTypeID = 0x10
	ConfigResponse             GetFrameTypeID = 0x20
	CurrentWeatherResponse     GetFrameTypeID = 0x30
	HistoryResponse            GetFrameTypeID = 0x40
	RequestHistoryResponse     GetFrameTypeID = 0x50
	RequestFirstConfigResponse GetFrameTypeID = 0x51
	RequestConfigResponse      GetFrameTypeID = 0x52
	RequestTimeResponse        GetFrameTypeID = 0x53
)

func (t GetFrameTypeID) String() string {
	var s string

	switch t {
	case ConfigResponse:
		s = "ConfigResponse"
	case DataWrittenResponse:
		s = "DataWrittenResponse"
	case CurrentWeatherResponse:
		s = "CurrentWeatherResponse"
	case HistoryResponse:
		s = "HistoryResponse"
	case RequestHistoryResponse:
		s = "RequestHistoryResponse"
	case RequestFirstConfigResponse:
		s = "RequestFirstConfigResponse"
	case RequestConfigResponse:
		s = "RequestConfigResponse"
	case RequestTimeResponse:
		s = "RequestTimeResponse"
	}

	return fmt.Sprintf("%s (0x%02x)", s, byte(t))
}

type GetFrame []byte

func NewGetFrame() GetFrame {
	f := make([]byte, 273)
	return GetFrame(f)
}

func (f GetFrame) Length() uint8 {
	return uint8(f[2])
}

func (f GetFrame) DeviceID() uint16 {
	return uint16(f[3])<<8 | uint16(f[4])
}

func (f GetFrame) LoggerID() uint8 {
	return uint8(f[5])
}

func (f GetFrame) TypeID() GetFrameTypeID {
	return GetFrameTypeID(f[6])
}

func (f GetFrame) String() string {
	return fmt.Sprintf("GetFrame{Length:%d DeviceID:0x%04x LoggerID:%02d TypeID:%s}", f.Length(), f.DeviceID(), f.LoggerID(), f.TypeID())
}

func (f *GetFrame) UnmarshalBinary(data []byte) error {
	if len(data) > len(*f) {
		return fmt.Errorf("unexpected length of data; expected %dB, got %dB", len(*f), len(data))
	}

	if data[0] != byte(0x0) { // klimalogg.GetFrame
		return fmt.Errorf("data can not be unmarshaled; found unexpected message type 0x%02x", data[0])
	}

	copy(*f, data)
	//*f = data
	return nil
}
