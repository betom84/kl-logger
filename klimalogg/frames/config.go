package frames

// ConfigResponseData contains current configuration of klimalogg console
type ConfigResponseFrame struct {
	GetFrame
}

func NewConfigResponseFrame() ConfigResponseFrame {
	return ConfigResponseFrame{NewGetFrame()}
}

func (f ConfigResponseFrame) CfgChecksum() uint16 {
	return uint16(f.GetFrame[126])<<8 | uint16(f.GetFrame[127])
}

func (f ConfigResponseFrame) SignalQuality() int {
	return int(f.GetFrame[7])
}

func (f ConfigResponseFrame) TimeZone() int {
	return int(f.GetFrame[9])
}

func (f ConfigResponseFrame) HistoryIntervall() int {
	return int(f.GetFrame[10]) * 5
}

func (f ConfigResponseFrame) Contrast() int {
	return int(f.GetFrame[8] >> 4)
}

func (f ConfigResponseFrame) IsAlarmEnabled() bool {
	return (f.GetFrame[8] & 0x8) == 0
}

func (f ConfigResponseFrame) IsDCFEnabled() bool {
	return (f.GetFrame[8] & 0x4) != 0
}

func (f ConfigResponseFrame) TimeFormat() string {
	return []string{"24h", "12h"}[(f.GetFrame[8] & 0x2)]
}

func (f ConfigResponseFrame) TempFormat() string {
	return []string{"Celcius", "Fahrenheit"}[(f.GetFrame[8] & 0x1)]
}

func (f ConfigResponseFrame) SensorMin() int {
	return 0
}

func (f ConfigResponseFrame) SensorMax() int {
	return 8
}

func (f ConfigResponseFrame) TemperatureMax(sensor int) float32 {
	offset := 11 + (sensor * 3)
	return toTemperature(f.GetFrame[offset:offset+2], 1)
}

func (f ConfigResponseFrame) TemperatureMin(sensor int) float32 {
	offset := 12 + (sensor * 3)
	return toTemperature(f.GetFrame[offset:offset+2], 2)
}

func (f ConfigResponseFrame) HumidityMax(sensor int) uint {
	offset := 38 + (sensor * 2)
	return toHumidity(f.GetFrame[offset])
}

func (f ConfigResponseFrame) HumidityMin(sensor int) uint {
	offset := 39 + (sensor * 2)
	return toHumidity(f.GetFrame[offset])
}

func (f ConfigResponseFrame) Description(sensor int) string {
	if sensor == 0 {
		return "INDOOR"
	}

	d := [8]byte{}
	offset := 61 + ((sensor - 1) * 8)
	copy(d[:], f.GetFrame[offset:offset+8])

	return toString(d)
}

func (f ConfigResponseFrame) IsTemperatureMaxAlarmSet(sensor int) bool {
	/*
	   Humidity0Max: 00 00 00 00 01
	   Humidity0Min: 00 00 00 00 02
	   Temp0Max:     00 00 00 00 04
	   Temp0Min:     00 00 00 00 08
	   Humidity1Max: 00 00 00 00 10
	   Humidity1Min: 00 00 00 00 20
	   Temp1Max:     00 00 00 00 40
	   Temp1Min:     00 00 00 00 80
	   ...
	*/

	alarm := f.GetFrame[56:62]

	switch sensor {
	case 0:
		return (alarm[4] & 0x40) != 0
	case 1:
		return (alarm[4] & 0x04) != 0
	case 2:
		return (alarm[3] & 0x40) != 0
	case 3:
		return (alarm[3] & 0x04) != 0
	case 4:
		return (alarm[2] & 0x40) != 0
	case 5:
		return (alarm[2] & 0x04) != 0
	case 6:
		return (alarm[1] & 0x40) != 0
	case 7:
		return (alarm[1] & 0x04) != 0
	case 8:
		return (alarm[0] & 0x40) != 0
	}

	return false
}

func (f ConfigResponseFrame) IsTemperatureMinAlarmSet(sensor int) bool {
	// todo
	return false
}

func (f ConfigResponseFrame) IsHumidityMaxAlarmSet(sensor int) bool {
	// todo
	return false
}

func (f ConfigResponseFrame) IsHumidityMinAlarmSet(sensor int) bool {
	// todo
	return false
}

// FirstConfigRequestData to pair klimalogg console and ask for current configuration
type FirstConfigRequestFrame struct {
	SetFrame
}

func NewFirstConfigRequestFrame() FirstConfigRequestFrame {
	f := NewSetFrame()

	f.SetTypeID(ConfigRequest)
	f.SetDeviceID(0xf0f0)
	f.SetLoggerID(0xff)
	f.Write(make([]byte, 7))

	c := FirstConfigRequestFrame{f}
	c.SetCfgChecksum(0xffff)
	c.SetDeviceID(0xffff)
	c.SetComInterval(8)

	return c
}

func (f FirstConfigRequestFrame) SetCfgChecksum(checksum int) {
	f.SetFrame[7] = byte(checksum >> 8)
	f.SetFrame[8] = byte(checksum)
}

func (f FirstConfigRequestFrame) SetComInterval(interval int) {
	f.SetFrame[9] = 0x80
	f.SetFrame[10] = byte(interval)
}

func (f FirstConfigRequestFrame) SetDeviceID(deviceID uint16) {
	f.SetFrame[11] = byte(deviceID >> 8)
	f.SetFrame[12] = byte(deviceID)
}
func (f FirstConfigRequestFrame) DeviceID() uint16 {
	return uint16(f.SetFrame[11])<<8 | uint16(f.SetFrame[12])
}

func (f FirstConfigRequestFrame) SetLoggerID(loggerID uint8) {
	f.SetFrame[13] = byte(loggerID)
}

func (f FirstConfigRequestFrame) LoggerID() uint8 {
	return uint8(f.SetFrame[13])
}
