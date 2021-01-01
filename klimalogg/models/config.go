package models

import (
	"fmt"

	"github.com/betom84/kl-logger/utils"
)

// Settings of klimalogg console
type Settings struct {
	Contrast   int
	Alert      bool
	DCF        bool
	TimeFormat string
	TempFormat string
}

// ConfigResponseData contains current configuration of klimalogg console
type ConfigResponseData struct {
	CfgChecksum     int
	SignalQuality   int
	Settings        Settings
	TimeZone        byte
	HistoryInterval byte
	HumidityMin     [1]uint
	HumidityMax     [1]uint
	TemperatureMin  [1]float32
	TemperatureMax  [1]float32
	AlarmSet        [5]byte
	Description     [1]string
}

// UnmarshalBinary response from klimalogg console
func (d *ConfigResponseData) UnmarshalBinary(data []byte) error {
	if len(data) != 122 {
		return fmt.Errorf("unexpected length of data; expected 122B, got %dB", len(data))
	}

	d.SignalQuality = int(data[0])
	d.TimeZone = data[2]
	d.HistoryInterval = data[3]

	d.Settings = Settings{
		Contrast:   int(data[1] >> 4),
		Alert:      (data[1] & 0x8) == 0,
		DCF:        (data[1] & 0x4) != 0,
		TimeFormat: []string{"24h", "12h"}[(data[1] & 0x2)],
		TempFormat: []string{"Celcius", "Fahrenheit"}[(data[1] & 0x1)],
	}

	d.TemperatureMax[0] = utils.ToTemperature(data[3:5], 1)
	d.TemperatureMin[0] = utils.ToTemperature(data[4:6], 2)
	d.HumidityMax[0] = utils.ToHumidity(data[31])
	d.HumidityMin[0] = utils.ToHumidity(data[32])

	for i, b := range data[49:54] {
		d.AlarmSet[i] = b
	}

	d.Description[0] = "INDOOR"

	d.CfgChecksum = int(uint16(data[119])<<8 | uint16(data[120]))

	return nil
}

// FirstConfigRequestData to pair klimalogg console and ask for current configuration
type FirstConfigRequestData struct {
	LoggerID uint8
	DeviceID uint16
}

// MarshalBinary request for klimalogg console
func (f FirstConfigRequestData) MarshalBinary() ([]byte, error) {
	var cfgChecksum uint16 = 0xffff
	var comInt uint8 = 8

	return []byte{byte(cfgChecksum >> 8), byte(cfgChecksum), 0x80, comInt, uint8(f.DeviceID >> 8), uint8(f.DeviceID), f.LoggerID}, nil
}
