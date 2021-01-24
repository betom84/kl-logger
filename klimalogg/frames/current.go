package frames

import "time"

// CurrentWeatherRequestFrame to ask klimalogg console for current weather data
type CurrentWeatherRequestFrame struct {
	SetFrame
}

func NewCurrentWeatherRequestFrame() CurrentWeatherRequestFrame {
	s := NewSetFrame()
	s.SetTypeID(CurrentWeatherRequest)
	s.Write(make([]byte, 7))

	f := CurrentWeatherRequestFrame{s}
	f.SetLastHistory(0xffffff)
	f.SetComInterval(8)

	return f
}

func (f CurrentWeatherRequestFrame) SetCfgChecksum(checksum uint16) {
	f.SetFrame[7] = byte(checksum >> 8)
	f.SetFrame[8] = byte(checksum)
}

func (f CurrentWeatherRequestFrame) SetComInterval(interval int) {
	f.SetFrame[9] = 0x80
	f.SetFrame[10] = byte(interval)
}

func (f CurrentWeatherRequestFrame) ComInterval() int {
	return int(f.SetFrame[10])
}

func (f CurrentWeatherRequestFrame) SetLastHistory(addr uint32) {
	f.SetFrame[11] = uint8(addr >> 16)
	f.SetFrame[12] = uint8(addr >> 8)
	f.SetFrame[13] = uint8(addr)
}

func (f CurrentWeatherRequestFrame) LastHistory() uint32 {
	return uint32(f.SetFrame[11])<<16 | uint32(f.SetFrame[12])<<8 | uint32(f.SetFrame[13])
}

// CurrentWeatherResponseData contains current weather data received from klimalogg console
type CurrentWeatherResponseFrame struct {
	GetFrame
}

func NewCurrentWeatherResponseFrame() CurrentWeatherResponseFrame {
	return CurrentWeatherResponseFrame{NewGetFrame()}
}

func (f CurrentWeatherResponseFrame) CfgChecksum() uint16 {
	return uint16(f.GetFrame[8])<<8 | uint16(f.GetFrame[9])
}

func (f CurrentWeatherResponseFrame) SignalQuality() int {
	return int(f.GetFrame[7])
}

func (f CurrentWeatherResponseFrame) getSensorOffset(sensor int) uint {
	if sensor < f.SensorMin() || sensor > f.SensorMax() {
		return 0
	}

	return uint(24 * sensor)
}

func (f CurrentWeatherResponseFrame) SensorMin() int {
	return 0
}

func (f CurrentWeatherResponseFrame) SensorMax() int {
	return 8
}

func (f CurrentWeatherResponseFrame) HumidityMaxTime(sensor int) time.Time {
	o := f.getSensorOffset(sensor)
	return toDateTime(f.GetFrame[10+o:14+o], HighNibble)
}

func (f CurrentWeatherResponseFrame) HumidityMinTime(sensor int) time.Time {
	o := f.getSensorOffset(sensor)
	return toDateTime(f.GetFrame[14+o:18+o], HighNibble)
}

func (f CurrentWeatherResponseFrame) HumidityMax(sensor int) uint {
	o := f.getSensorOffset(sensor)
	return toHumidity(f.GetFrame[18+o])
}

func (f CurrentWeatherResponseFrame) HumidityMin(sensor int) uint {
	o := f.getSensorOffset(sensor)
	return toHumidity(f.GetFrame[19+o])
}

func (f CurrentWeatherResponseFrame) Humidity(sensor int) uint {
	o := f.getSensorOffset(sensor)
	return toHumidity(f.GetFrame[20+o])
}

func (f CurrentWeatherResponseFrame) TemperatureMaxTime(sensor int) time.Time {
	o := f.getSensorOffset(sensor)
	return toDateTime(f.GetFrame[21+o:26+o], LowNibble)
}

func (f CurrentWeatherResponseFrame) TemperatureMinTime(sensor int) time.Time {
	o := f.getSensorOffset(sensor)
	return toDateTime(f.GetFrame[25+o:30+o], LowNibble)
}

func (f CurrentWeatherResponseFrame) TemperatureMax(sensor int) float32 {
	o := f.getSensorOffset(sensor)
	return toTemperature(f.GetFrame[29+o:31+o], LowNibble)
}

func (f CurrentWeatherResponseFrame) TemperatureMin(sensor int) float32 {
	o := f.getSensorOffset(sensor)
	return toTemperature(f.GetFrame[31+o:33+o], HighNibble)
}

func (f CurrentWeatherResponseFrame) Temperature(sensor int) float32 {
	o := f.getSensorOffset(sensor)
	return toTemperature(f.GetFrame[32+o:34+o], LowNibble)
}
