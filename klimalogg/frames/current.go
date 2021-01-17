package frames

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

func (f CurrentWeatherResponseFrame) HumidityMax(sensor int) uint {
	return toHumidity(f.GetFrame[18])
}

func (f CurrentWeatherResponseFrame) HumidityMin(sensor int) uint {
	return toHumidity(f.GetFrame[19])
}

func (f CurrentWeatherResponseFrame) Humidity(sensor int) uint {
	return toHumidity(f.GetFrame[20])
}

func (f CurrentWeatherResponseFrame) TemperatureMax(sensor int) float32 {
	return toTemperature(f.GetFrame[29:31], 2)
}

func (f CurrentWeatherResponseFrame) TemperatureMin(sensor int) float32 {
	return toTemperature(f.GetFrame[31:33], 1)
}

func (f CurrentWeatherResponseFrame) Temperature(sensor int) float32 {
	return toTemperature(f.GetFrame[32:34], 2)
}
