package repository

import "time"

// Default repository to store klimalogg data
var Default = InMemory{}

// InMemory repository to store data in memory
type InMemory struct {
	temperature float32
	humidity    uint
	lastUpdate  time.Time
}

// UpdateWeather
func (r *InMemory) UpdateWeather(temperature float32, humidity uint) {
	r.temperature = temperature
	r.humidity = humidity

	r.lastUpdate = time.Now()
}

// UpdateSettings
func (r *InMemory) UpdateSettings(contrast int, alert bool, dcf bool, timeFormat string, tempFormat string) {
	r.lastUpdate = time.Now()
}

// Temperature
func (r InMemory) Temperature() float32 {
	return r.temperature
}

// Humidity
func (r InMemory) Humidity() uint {
	return r.humidity
}

// LastUpdate
func (r InMemory) LastUpdate() time.Time {
	return r.lastUpdate
}
