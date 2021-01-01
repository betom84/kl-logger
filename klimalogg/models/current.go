package models

import (
	"fmt"

	"github.com/betom84/kl-logger/utils"
)

// CurrentWeatherRequestData to ask klimalogg console for current weather data
type CurrentWeatherRequestData struct {
	CfgChecksum uint16
}

// MarshalBinary request for klimalogg console
func (f CurrentWeatherRequestData) MarshalBinary() ([]byte, error) {
	var comInt uint8 = 8
	var lastHistoryAddr uint32 = 0xffffff

	return []byte{byte(f.CfgChecksum >> 8), byte(f.CfgChecksum), 0x80, comInt, uint8(lastHistoryAddr >> 16), uint8(lastHistoryAddr >> 8), uint8(lastHistoryAddr)}, nil
}

// CurrentWeatherResponseData contains current weather data received from klimalogg console
type CurrentWeatherResponseData struct {
	SignalQuality  int
	CfgChecksum    int
	Humidity       [1]uint
	HumidityMin    [1]uint
	HumidityMax    [1]uint
	Temperature    [1]float32
	TemperatureMin [1]float32
	TemperatureMax [1]float32
}

// UnmarshalBinary response from klimalogg console
func (f *CurrentWeatherResponseData) UnmarshalBinary(data []byte) error {
	if len(data) != 226 {
		return fmt.Errorf("unexpected length of data; expected 226B, got %dB", len(data))
	}

	f.SignalQuality = int(data[0])
	f.CfgChecksum = int(uint16(data[1])<<8 | uint16(data[2]))

	//f.HumidityMaxDT[0] = data[3:7]
	//f.HumidityMinDT[0] = data[7:11]
	f.HumidityMax[0] = utils.ToHumidity(data[11])
	f.HumidityMin[0] = utils.ToHumidity(data[12])
	f.Humidity[0] = utils.ToHumidity(data[13])

	//f.TemperatureMaxDT[0] = data[14l:18h]
	//f.TemperatureMinDT[0] = data[18l:22h]
	f.TemperatureMax[0] = utils.ToTemperature(data[22:24], 2)
	f.TemperatureMin[0] = utils.ToTemperature(data[24:26], 1)
	f.Temperature[0] = utils.ToTemperature(data[25:27], 2)

	return nil
}
