package repository

type Configuration interface {
	SignalQuality() int
	TimeZone() int
	HistoryIntervall() int
	CfgChecksum() uint16

	Contrast() int
	IsAlarmEnabled() bool
	IsDCFEnabled() bool
	TempFormat() string
	TimeFormat() string

	Description(int) string

	TemperatureMin(int) float32
	TemperatureMax(int) float32
	IsTemperatureMinAlarmSet(int) bool
	IsTemperatureMaxAlarmSet(int) bool

	HumidityMin(int) uint
	HumidityMax(int) uint
	IsHumidityMinAlarmSet(int) bool
	IsHumidityMaxAlarmSet(int) bool
}

type WeatherSample interface {
	SignalQuality() int
	CfgChecksum() uint16

	Temperature(int) float32
	TemperatureMin(int) float32
	TemperatureMax(int) float32

	Humidity(int) uint
	HumidityMin(int) uint
	HumidityMax(int) uint
}
