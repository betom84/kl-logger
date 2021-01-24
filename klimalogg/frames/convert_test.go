package frames

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToInt(t *testing.T) {
	assert.Equal(t, uint(17), toInt(0x17))
	assert.Equal(t, uint(17), toInt(0x17, HighNibble, LowNibble))

	assert.Equal(t, uint(1), toInt(0x17, HighNibble))
	assert.Equal(t, uint(7), toInt(0x17, LowNibble))
}

func TestToHumditiy(t *testing.T) {
	assert.Equal(t, uint(50), toHumidity(0x50))
	assert.Equal(t, uint(99), toHumidity(0x99))
}

func TestToTemperature(t *testing.T) {
	assert.InDelta(t, float32(23.2), toTemperature([]byte{0x63, 0x20}, HighNibble), 0.01)
	assert.InDelta(t, float32(23.2), toTemperature([]byte{0x06, 0x32}, LowNibble), 0.01)
}

func TestToHexNumber(t *testing.T) {
	assert.Equal(t, uint8(0x01), toHexNumber(1))
	assert.Equal(t, uint8(0x50), toHexNumber(50))
	assert.Equal(t, uint8(0x99), toHexNumber(99))
	assert.Equal(t, uint8(0x00), toHexNumber(100))
}

func TestToDateTime(t *testing.T) {
	expected, err := time.Parse("2006-01-02 15:04", "2005-02-13 06:03")
	assert.NoError(t, err)
	assert.Equal(t, expected.Format(time.RFC3339), toDateTime([]byte{0x05, 0x21, 0x36, 0x21}, HighNibble).Format(time.RFC3339))

	expected, err = time.Parse("2006-01-02 15:04", "2021-01-14 18:39")
	assert.NoError(t, err)
	assert.Equal(t, expected.Format(time.RFC3339), toDateTime([]byte{0x21, 0x11, 0x48, 0xd9}, HighNibble).Format(time.RFC3339))

}
