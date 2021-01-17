package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/betom84/kl-logger/repository"
)

func getCurrentWeather() http.HandlerFunc {
	type sensor struct {
		ID             int     `json:"ID"`
		Name           string  `json:"name"`
		Temperature    float32 `json:"temperature"`
		TemperatureMin float32 `json:"temperatureMin"`
		TemperatureMax float32 `json:"temperatureMax"`
		Humidity       uint    `json:"humidity"`
		HumidityMin    uint    `json:"humidityMin"`
		HumidityMax    uint    `json:"humidityMax"`
	}

	type response struct {
		UpdatedAt     string   `json:"updatedAt"`
		SignalQuality int      `json:"signalQuality"`
		Sensors       []sensor `json:"sensors"`
	}

	repo := repository.Default

	return func(w http.ResponseWriter, r *http.Request) {
		weather := repo.CurrentWeather()
		config := repo.CurrentConfig()

		sensors := make([]sensor, 0)
		for sensorID := 0; sensorID < 9; sensorID++ {
			sensors = append(sensors, sensor{
				ID:             sensorID,
				Name:           config.Description(sensorID),
				Temperature:    weather.Temperature(sensorID),
				TemperatureMin: weather.TemperatureMin(sensorID),
				TemperatureMax: weather.TemperatureMax(sensorID),
				Humidity:       weather.Humidity(sensorID),
				HumidityMin:    weather.HumidityMin(sensorID),
				HumidityMax:    weather.HumidityMax(sensorID),
			})
		}

		response := response{
			UpdatedAt:     repo.LastWeatherUpdate().Format(time.RFC3339),
			SignalQuality: weather.SignalQuality(),
			Sensors:       sensors,
		}

		json.NewEncoder(w).Encode(response)
		w.Header().Set("Content-Type", "application/json")
	}
}

func getCurrentConfig() http.HandlerFunc {
	type sensor struct {
		ID                       int     `json:"ID"`
		Name                     string  `json:"name"`
		TemperatureMin           float32 `json:"temperatureMin"`
		TemperatureMax           float32 `json:"temperatureMax"`
		IsTemperatureMinAlarmSet bool    `json:"temperatureMinAlarmSet"`
		IsTemperatureMaxAlarmSet bool    `json:"temperatureMaxAlarmSet"`
		HumidityMin              uint    `json:"humidityMin"`
		HumidityMax              uint    `json:"humidityMax"`
		IsHumidityMinAlarmSet    bool    `json:"humidityMinAlarmSet"`
		IsHumidityMaxAlarmSet    bool    `json:"humidityMaxAlarmSet"`
	}

	type response struct {
		UpdatedAt  string   `json:"updatedAt"`
		Alarm      bool     `json:"alarm"`
		DCF        bool     `json:"dcf"`
		TimeFormat string   `json:"time"`
		TempFormat string   `json:"temperature"`
		Sensors    []sensor `json:"sensors"`
	}

	repo := repository.Default

	return func(w http.ResponseWriter, r *http.Request) {
		config := repo.CurrentConfig()

		sensors := make([]sensor, 0)
		for sensorID := 0; sensorID < 9; sensorID++ {
			sensors = append(sensors, sensor{
				ID:                       sensorID,
				Name:                     config.Description(sensorID),
				TemperatureMin:           config.TemperatureMin(sensorID),
				TemperatureMax:           config.TemperatureMax(sensorID),
				IsTemperatureMinAlarmSet: config.IsTemperatureMinAlarmSet(sensorID),
				IsTemperatureMaxAlarmSet: config.IsTemperatureMaxAlarmSet(sensorID),
				HumidityMin:              config.HumidityMin(sensorID),
				HumidityMax:              config.HumidityMax(sensorID),
				IsHumidityMinAlarmSet:    config.IsHumidityMinAlarmSet(sensorID),
				IsHumidityMaxAlarmSet:    config.IsHumidityMaxAlarmSet(sensorID),
			})
		}

		response := response{
			UpdatedAt:  repo.LastConfigUpdate().Format(time.RFC3339),
			Alarm:      config.IsAlarmEnabled(),
			DCF:        config.IsDCFEnabled(),
			TimeFormat: config.TimeFormat(),
			TempFormat: config.TempFormat(),
			Sensors:    sensors,
		}

		json.NewEncoder(w).Encode(response)
		w.Header().Set("Content-Type", "application/json")
	}
}
