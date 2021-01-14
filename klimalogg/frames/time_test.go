package frames_test

import (
	"testing"
	"time"

	"github.com/betom84/kl-logger/klimalogg/frames"
	"github.com/stretchr/testify/assert"
)

func TestSendTimeFrame(t *testing.T) {
	now, err := time.Parse("Mon 2006-01-02 15:04:05", "Thu 2014-10-30 21:58:25")
	assert.NoError(t, err)

	frame := frames.NewSendTimeFrame()
	frame.SetDeviceID(263)
	frame.SetTime(now)

	actual, err := frame.MarshalBinary()
	expected := []byte{0xd5, 0x00, 0x0d, 0x01, 0x07, 0x00, 0x60, 0x0, 0x0, 0x25, 0x58, 0x21, 0x04, 0x03, 0x41, 0x01}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
