package frames_test

import (
	"testing"

	"github.com/betom84/kl-logger/klimalogg/frames"
	"github.com/betom84/kl-logger/repository"
	"github.com/stretchr/testify/assert"
)

func TestConigResponseData(t *testing.T) {
	sampleFrame := []byte{
		0x00, 0x00, 0x7d, 0x01, 0x07, 0x00, 0x20, 0x64, 0x54, 0x00, 0x00, 0x80, 0x04, 0x00, 0x80, 0x04,
		0x00, 0x80, 0x04, 0x00, 0x80, 0x04, 0x00, 0x80, 0x04, 0x00, 0x80, 0x04, 0x00, 0x80, 0x04, 0x00,
		0x80, 0x04, 0x00, 0x80, 0x04, 0x00, 0x70, 0x20, 0x70, 0x20, 0x70, 0x20, 0x70, 0x20, 0x70, 0x20,
		0x70, 0x20, 0x70, 0x20, 0x70, 0x20, 0x70, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x08, 0xd3, 0xd5, 0x7f, 0xd2, 0x00, 0x00, 0x00, 0x00, 0x07, 0xb8, 0x76, 0xd2, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x71, 0x7f, 0x97, 0x00, 0x00, 0x00, 0x00, 0x85, 0xf4, 0x4c, 0x56, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0x1a, 0xb1,
		0x6c,
	}

	frame := frames.NewConfigResponseFrame()
	err := frame.UnmarshalBinary(sampleFrame)
	assert.NoError(t, err)

	assert.Implements(t, (*repository.Configuration)(nil), frames.NewConfigResponseFrame())

	assert.Equal(t, uint16(263), frame.DeviceID())
	assert.Equal(t, uint8(0), frame.LoggerID())
	assert.Equal(t, frames.ConfigResponse, frame.TypeID())
	assert.Equal(t, uint8(125), frame.Length())

	assert.Equal(t, uint16(6833), frame.CfgChecksum())
	assert.Equal(t, 100, frame.SignalQuality())
	assert.Equal(t, 0, frame.TimeZone())
	assert.Equal(t, 0, frame.HistoryIntervall())

	assert.Equal(t, 5, frame.Contrast())
	assert.Equal(t, true, frame.IsAlarmEnabled())
	assert.Equal(t, true, frame.IsDCFEnabled())
	assert.Equal(t, "24h", frame.TimeFormat())
	assert.Equal(t, "Celcius", frame.TempFormat())

	for sensor := 0; sensor <= 8; sensor++ {
		assert.Equal(t, float32(0), frame.TemperatureMin(sensor))
		assert.Equal(t, float32(40), frame.TemperatureMax(sensor))
		assert.Equal(t, uint(20), frame.HumidityMin(sensor))
		assert.Equal(t, uint(70), frame.HumidityMax(sensor))

		assert.Equal(t, false, frame.IsTemperatureMaxAlarmSet(sensor))
		assert.Equal(t, false, frame.IsTemperatureMinAlarmSet(sensor))
		assert.Equal(t, false, frame.IsHumidityMinAlarmSet(sensor))
		assert.Equal(t, false, frame.IsHumidityMaxAlarmSet(sensor))
	}

	assert.Equal(t, "INDOOR", frame.Description(0))
	assert.Equal(t, "AUSSEN", frame.Description(1))
	assert.Equal(t, "ALTAR", frame.Description(2))
	assert.Equal(t, "TURM", frame.Description(3))
	assert.Equal(t, "ORGEL", frame.Description(4))
}

func TestFirstConfigResponse(t *testing.T) {
	frame := frames.NewFirstConfigRequestFrame()

	assert.Equal(t, uint16(0xffff), frame.DeviceID())
	assert.Equal(t, uint8(0), frame.LoggerID())

	actual, err := frame.MarshalBinary()
	expected := []byte{0xd5, 0x0, 0x0b, 0xf0, 0xf0, 0xff, 0x3, 0xff, 0xff, 0x80, 0x8, 0xff, 0xff, 0x0}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
