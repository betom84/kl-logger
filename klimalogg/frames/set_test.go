package frames_test

import (
	"testing"

	"github.com/betom84/kl-logger/klimalogg/frames"
	"github.com/stretchr/testify/assert"
)

func TestSetFrame(t *testing.T) {
	frame := frames.NewSetFrame()

	frame.SetDeviceID(4711)
	frame.SetLoggerID(1)
	frame.SetTypeID(frames.SendTime)

	assert.Equal(t, 4, frame.Length())
	assert.Equal(t, 7, len(frame))
	assert.Equal(t, 273, cap(frame))

	assert.Equal(t, uint16(4711), frame.DeviceID())
	assert.Equal(t, 1, frame.LoggerID())
	assert.Equal(t, frames.SendTime, frame.TypeID())

	assert.Equal(t, "SetFrame{Length:4 DeviceID:0x1267 LoggerID:01 TypeID:SendTime (0x60)}", frame.String())

	actual, err := frame.MarshalBinary()
	expected := []byte{0xd5, 0x0, 0x4, 0x12, 0x67, 0x1, 0x60}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	data := []byte("data")
	frame.Write(data)

	assert.Equal(t, 8, frame.Length())

	actual, err = frame.MarshalBinary()
	expected = []byte{0xd5, 0x0, 0x8, 0x12, 0x67, 0x1, 0x60, 0x64, 0x61, 0x74, 0x61}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
