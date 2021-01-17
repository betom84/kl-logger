package repository

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Default repository to store klimalogg data
var Default = &InMemory{}

// InMemory repository to store data in memory
type InMemory struct {
	temperature float32
	humidity    uint

	currentConfig    Configuration
	lastConfigUpdate time.Time

	currentWeather    WeatherSample
	lastWeatherUpdate time.Time
}

func (r *InMemory) Listen() chan<- interface{} {
	c := make(chan interface{})

	go func() {
		for {
			i := <-c

			switch t := i.(type) {
			case WeatherSample:
				r.updateWeather(t)
			case Configuration:
				r.updateConfiguration(t)
			default:
				logrus.WithField("interface", i).Warn("repository: unrecognized interface type")
			}
		}
	}()

	return c
}

func (r *InMemory) updateWeather(s WeatherSample) {
	r.currentWeather = s
	r.lastWeatherUpdate = time.Now()
}

func (r *InMemory) updateConfiguration(c Configuration) {
	r.currentConfig = c
	r.lastConfigUpdate = time.Now()
}

func (r InMemory) CurrentWeather() WeatherSample {
	return r.currentWeather
}

func (r InMemory) LastWeatherUpdate() time.Time {
	return r.lastWeatherUpdate
}

func (r InMemory) CurrentConfig() Configuration {
	return r.currentConfig
}

func (r InMemory) LastConfigUpdate() time.Time {
	return r.lastConfigUpdate
}
