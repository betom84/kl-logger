package frames_test

import (
	"testing"

	"github.com/betom84/kl-logger/klimalogg/frames"
	"github.com/stretchr/testify/assert"
)

func TestGetFrame(t *testing.T) {
	frame := frames.NewGetFrame()
	assert.Len(t, frame, 273)

	sampleFrame := []byte{0x00, 0x00, 0xe5, 0x01, 0x07, 0x00, 0x30, 0x64, 0x1a, 0xb1, 0x13, 0x62, 0x10, 0x52, 0x14, 0x91, 0x85, 0xa3, 0x98, 0x32, 0x55, 0x01, 0x49, 0x17, 0x5d, 0x81, 0x41, 0x27, 0x43, 0x87, 0x36, 0x38, 0x56, 0x56, 0x14, 0xa1, 0x87, 0x29, 0x14, 0x91, 0x85, 0xa4, 0x89, 0x38, 0xaa, 0x01, 0x49, 0x17, 0x5d, 0x51, 0x49, 0x23, 0x75, 0x17, 0x44, 0x49, 0x4a, 0xaa, 0x14, 0xa1, 0x41, 0xc5, 0x14, 0x91, 0x85, 0xb2, 0x91, 0x40, 0x64, 0x01, 0x49, 0x17, 0x5e, 0x91, 0x4a, 0x22, 0x7b, 0x27, 0x32, 0x50, 0x26, 0x42, 0x14, 0xa2, 0x04, 0xc0, 0x14, 0x91, 0x85, 0xa4, 0x84, 0x38, 0x67, 0x01, 0x49, 0x17, 0x5d, 0x61, 0x4a, 0x22, 0x6c, 0x07, 0x44, 0x50, 0x06, 0x38, 0x14, 0xa2, 0x06, 0xc7, 0x14, 0x91, 0x85, 0xb2, 0x87, 0x41, 0xaa, 0x01, 0x49, 0x17, 0x5d, 0x31, 0x49, 0x19, 0x81, 0x57, 0x40, 0x52, 0x1a, 0xaa, 0xaa, 0x4a, 0xa4, 0xaa, 0xaa, 0x4a, 0xa4, 0xaa, 0xaa, 0xaa, 0xaa, 0x0a, 0xa4, 0xaa, 0x4a, 0xaa, 0xa4, 0xaa, 0x4a, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0x4a, 0xa4, 0xaa, 0xaa, 0x4a, 0xa4, 0xaa, 0xaa, 0xaa, 0xaa, 0x0a, 0xa4, 0xaa, 0x4a, 0xaa, 0xa4, 0xaa, 0x4a, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0x4a, 0xa4, 0xaa, 0xaa, 0x4a, 0xa4, 0xaa, 0xaa, 0xaa, 0xaa, 0x0a, 0xa4, 0xaa, 0x4a, 0xaa, 0xa4, 0xaa, 0x4a, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0x4a, 0xa4, 0xaa, 0xaa, 0x4a, 0xa4, 0xaa, 0xaa, 0xaa, 0xaa, 0x0a, 0xa4, 0xaa, 0x4a, 0xaa, 0xa4, 0xaa, 0x4a, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x39, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	err := frame.UnmarshalBinary(sampleFrame)

	assert.NoError(t, err)

	assert.Equal(t, uint16(0x0107), frame.DeviceID())
	assert.Equal(t, uint8(0), frame.LoggerID())
	assert.Equal(t, frames.CurrentWeatherResponse, frame.TypeID())
	assert.Equal(t, uint8(229), frame.Length())

	assert.Equal(t, "GetFrame{Length:229 DeviceID:0x0107 LoggerID:00 TypeID:CurrentWeatherResponse (0x30)}", frame.String())
}
