package testdata

import (
	"time"

	"github.com/betom84/kl-logger/repository"
	"github.com/stretchr/testify/mock"
)

func MockRepository(w repository.WeatherSample, c repository.Configuration) *RepositoryMock {
	m := &RepositoryMock{mock.Mock{}}
	m.On("Weather").Return(w)
	m.On("Config").Return(c)
	m.On("LastWeatherUpdate").Return(time.Now())
	m.On("LastConfigUpdate").Return(time.Now())

	return m
}

type RepositoryMock struct{ mock.Mock }

func (m *RepositoryMock) NewListener() chan<- interface{} { return nil }
func (m *RepositoryMock) LastWeatherUpdate() time.Time    { return m.Called().Get(0).(time.Time) }
func (m *RepositoryMock) LastConfigUpdate() time.Time     { return m.Called().Get(0).(time.Time) }
func (m *RepositoryMock) Weather() repository.WeatherSample {
	return m.Called().Get(0).(repository.WeatherSample)
}
func (m *RepositoryMock) Config() repository.Configuration {
	return m.Called().Get(0).(repository.Configuration)
}

func MockWeatherSample(cb func(*WeatherSampleMock)) *WeatherSampleMock {
	m := &WeatherSampleMock{mock.Mock{}}

	if cb != nil {
		cb(m)
	}

	m.On("SignalQuality").Maybe().Return(100)
	m.On("CfgChecksum").Maybe().Return(0xabcd)
	m.On("SensorMin").Maybe().Return(0)
	m.On("SensorMax").Maybe().Return(8)
	m.On("IsSensorActive").Maybe().Return(true)
	m.On("Temperature", mock.Anything).Maybe().Return(float32(+20.0))
	m.On("TemperatureMax", mock.Anything).Maybe().Return(float32(+30.0))
	m.On("TemperatureMin", mock.Anything).Maybe().Return(float32(-10.0))
	m.On("TemperatureMaxTime", mock.Anything).Maybe().Return(time.Now())
	m.On("TemperatureMinTime", mock.Anything).Maybe().Return(time.Now())
	m.On("Humidity", mock.Anything).Maybe().Return(uint(50))
	m.On("HumidityMax", mock.Anything).Maybe().Return(uint(80))
	m.On("HumidityMin", mock.Anything).Maybe().Return(uint(20))
	m.On("HumidityMaxTime", mock.Anything).Maybe().Return(time.Now())
	m.On("HumidityMinTime", mock.Anything).Maybe().Return(time.Now())

	return m
}

type WeatherSampleMock struct{ mock.Mock }

func (m *WeatherSampleMock) SignalQuality() int           { return m.Called().Get(0).(int) }
func (m *WeatherSampleMock) CfgChecksum() uint16          { return m.Called().Get(0).(uint16) }
func (m *WeatherSampleMock) SensorMin() int               { return m.Called().Get(0).(int) }
func (m *WeatherSampleMock) SensorMax() int               { return m.Called().Get(0).(int) }
func (m *WeatherSampleMock) IsSensorActive(s int) bool    { return m.Called().Get(0).(bool) }
func (m *WeatherSampleMock) Temperature(s int) float32    { return m.Called(s).Get(0).(float32) }
func (m *WeatherSampleMock) TemperatureMin(s int) float32 { return m.Called(s).Get(0).(float32) }
func (m *WeatherSampleMock) TemperatureMinTime(s int) time.Time {
	return m.Called(s).Get(0).(time.Time)
}
func (m *WeatherSampleMock) TemperatureMax(s int) float32 { return m.Called(s).Get(0).(float32) }
func (m *WeatherSampleMock) TemperatureMaxTime(s int) time.Time {
	return m.Called(s).Get(0).(time.Time)
}
func (m *WeatherSampleMock) Humidity(s int) uint             { return m.Called(s).Get(0).(uint) }
func (m *WeatherSampleMock) HumidityMin(s int) uint          { return m.Called(s).Get(0).(uint) }
func (m *WeatherSampleMock) HumidityMinTime(s int) time.Time { return m.Called(s).Get(0).(time.Time) }
func (m *WeatherSampleMock) HumidityMax(s int) uint          { return m.Called(s).Get(0).(uint) }
func (m *WeatherSampleMock) HumidityMaxTime(s int) time.Time { return m.Called(s).Get(0).(time.Time) }

func MockConfiguration(cb func(*ConfigurationMock)) *ConfigurationMock {
	m := &ConfigurationMock{mock.Mock{}}

	if cb != nil {
		cb(m)
	}

	m.On("SignalQuality").Maybe().Return(100)
	m.On("TimeZone").Maybe().Return(1)
	m.On("HistoryIntervall").Maybe().Return(15)
	m.On("CfgChecksum").Maybe().Return(0xabcd)
	m.On("Contract").Maybe().Return(5)
	m.On("IsAlarmEnabled").Maybe().Return(true)
	m.On("IsDCFEnabled").Maybe().Return(true)
	m.On("TempFormat").Maybe().Return("Fahrenheit")
	m.On("TimeFormat").Maybe().Return("24h")
	m.On("SensorMin").Maybe().Return(0)
	m.On("SensorMax").Maybe().Return(8)
	m.On("Description", mock.Anything).Maybe().Return("INDOOR")
	m.On("TemperatureMax", mock.Anything).Maybe().Return(float32(+30.0))
	m.On("TemperatureMin", mock.Anything).Maybe().Return(float32(-10.0))
	m.On("IsTemperatureMinAlarmSet", mock.Anything).Maybe().Return(false)
	m.On("IsTemperatureMaxAlarmSet", mock.Anything).Maybe().Return(false)
	m.On("HumidityMax", mock.Anything).Maybe().Return(uint(80))
	m.On("HumidityMin", mock.Anything).Maybe().Return(uint(20))
	m.On("IsHumidityMinAlarmSet", mock.Anything).Maybe().Return(false)
	m.On("IsHumidityMaxAlarmSet", mock.Anything).Maybe().Return(false)

	return m
}

type ConfigurationMock struct{ mock.Mock }

func (m *ConfigurationMock) SignalQuality() int                  { return m.Called().Get(0).(int) }
func (m *ConfigurationMock) TimeZone() int                       { return m.Called().Get(0).(int) }
func (m *ConfigurationMock) HistoryIntervall() int               { return m.Called().Get(0).(int) }
func (m *ConfigurationMock) CfgChecksum() uint16                 { return m.Called().Get(0).(uint16) }
func (m *ConfigurationMock) Contrast() int                       { return m.Called().Get(0).(int) }
func (m *ConfigurationMock) IsAlarmEnabled() bool                { return m.Called().Get(0).(bool) }
func (m *ConfigurationMock) IsDCFEnabled() bool                  { return m.Called().Get(0).(bool) }
func (m *ConfigurationMock) TempFormat() string                  { return m.Called().Get(0).(string) }
func (m *ConfigurationMock) TimeFormat() string                  { return m.Called().Get(0).(string) }
func (m *ConfigurationMock) SensorMin() int                      { return m.Called().Get(0).(int) }
func (m *ConfigurationMock) SensorMax() int                      { return m.Called().Get(0).(int) }
func (m *ConfigurationMock) Description(s int) string            { return m.Called(s).Get(0).(string) }
func (m *ConfigurationMock) TemperatureMin(s int) float32        { return m.Called(s).Get(0).(float32) }
func (m *ConfigurationMock) TemperatureMax(s int) float32        { return m.Called(s).Get(0).(float32) }
func (m *ConfigurationMock) IsTemperatureMinAlarmSet(s int) bool { return m.Called(s).Get(0).(bool) }
func (m *ConfigurationMock) IsTemperatureMaxAlarmSet(s int) bool { return m.Called(s).Get(0).(bool) }
func (m *ConfigurationMock) HumidityMin(s int) uint              { return m.Called(s).Get(0).(uint) }
func (m *ConfigurationMock) HumidityMax(s int) uint              { return m.Called(s).Get(0).(uint) }
func (m *ConfigurationMock) IsHumidityMinAlarmSet(s int) bool    { return m.Called(s).Get(0).(bool) }
func (m *ConfigurationMock) IsHumidityMaxAlarmSet(s int) bool    { return m.Called(s).Get(0).(bool) }
