package repository

import "time"

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

	SensorMin() int
	SensorMax() int

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
	SensorMin() int
	SensorMax() int

	IsSensorActive(int) bool

	Temperature(int) float32
	TemperatureMin(int) float32
	TemperatureMinTime(int) time.Time
	TemperatureMax(int) float32
	TemperatureMaxTime(int) time.Time

	Humidity(int) uint
	HumidityMin(int) uint
	HumidityMinTime(int) time.Time
	HumidityMax(int) uint
	HumidityMaxTime(int) time.Time
}
