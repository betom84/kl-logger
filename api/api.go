package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/betom84/kl-logger/repository"
	"github.com/go-chi/chi"
)

type sensorWeather struct {
	ID                 int     `json:"ID"`
	Name               string  `json:"name"`
	Temperature        float32 `json:"temperature"`
	TemperatureMin     float32 `json:"temperatureMin"`
	TemperatureMinTime string  `json:"temperatureMinTime"`
	TemperatureMax     float32 `json:"temperatureMax"`
	TemperatureMaxTime string  `json:"temperatureMaxTime"`
	Humidity           uint    `json:"humidity"`
	HumidityMin        uint    `json:"humidityMin"`
	HumidityMinTime    string  `json:"humidityMinTim"`
	HumidityMax        uint    `json:"humidityMax"`
	HumidityMaxTime    string  `json:"humidityMaxTime"`
}

func newSensorWeather(sensorID int, name string, weather repository.WeatherSample) sensorWeather {
	return sensorWeather{
		ID:                 sensorID,
		Name:               name,
		Temperature:        weather.Temperature(sensorID),
		TemperatureMin:     weather.TemperatureMin(sensorID),
		TemperatureMinTime: weather.TemperatureMinTime(sensorID).Format(time.RFC3339),
		TemperatureMax:     weather.TemperatureMax(sensorID),
		TemperatureMaxTime: weather.TemperatureMaxTime(sensorID).Format(time.RFC3339),
		Humidity:           weather.Humidity(sensorID),
		HumidityMin:        weather.HumidityMin(sensorID),
		HumidityMinTime:    weather.HumidityMinTime(sensorID).Format(time.RFC3339),
		HumidityMax:        weather.HumidityMax(sensorID),
		HumidityMaxTime:    weather.HumidityMaxTime(sensorID).Format(time.RFC3339),
	}
}

func GetWeather(repo repository.Repository) http.HandlerFunc {

	type response struct {
		UpdatedAt     string          `json:"updatedAt"`
		SignalQuality int             `json:"signalQuality"`
		Sensors       []sensorWeather `json:"sensors"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		weather := repo.Weather()
		config := repo.Config()

		if weather == nil || config == nil {
			return
		}

		sensors := make([]sensorWeather, 0)
		for sensorID := weather.SensorMin(); sensorID <= weather.SensorMax(); sensorID++ {
			sensors = append(sensors, newSensorWeather(sensorID, config.Description(sensorID), weather))
		}

		response := response{
			UpdatedAt:     repo.LastWeatherUpdate().Format(time.RFC3339),
			SignalQuality: weather.SignalQuality(),
			Sensors:       sensors,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func GetWeatherBySensor(repo repository.Repository) http.HandlerFunc {

	type response struct {
		sensorWeather
		UpdatedAt     string `json:"updatedAt"`
		SignalQuality int    `json:"signalQuality"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		weather := repo.Weather()
		config := repo.Config()

		if weather == nil || config == nil {
			return
		}

		var sensorID int
		var err error
		if sensorID, err = strconv.Atoi(chi.URLParam(r, "sensor")); err != nil ||
			sensorID < weather.SensorMin() ||
			sensorID > weather.SensorMax() {
			panic("invalid sensor id")
		}

		response := response{
			UpdatedAt:     repo.LastWeatherUpdate().Format(time.RFC3339),
			SignalQuality: weather.SignalQuality(),
			sensorWeather: newSensorWeather(sensorID, config.Description(sensorID), weather),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

type sensorConfig struct {
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

func newSensorConfig(sensorID int, config repository.Configuration) sensorConfig {
	return sensorConfig{
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
	}
}

func GetConfig(repo repository.Repository) http.HandlerFunc {

	type response struct {
		UpdatedAt  string         `json:"updatedAt"`
		Alarm      bool           `json:"alarm"`
		DCF        bool           `json:"dcf"`
		TimeFormat string         `json:"time"`
		TempFormat string         `json:"temperature"`
		Sensors    []sensorConfig `json:"sensors"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		config := repo.Config()
		if config == nil {
			return
		}

		sensors := make([]sensorConfig, 0)
		for sensorID := 0; sensorID < 9; sensorID++ {
			sensors = append(sensors, newSensorConfig(sensorID, config))
		}

		response := response{
			UpdatedAt:  repo.LastConfigUpdate().Format(time.RFC3339),
			Alarm:      config.IsAlarmEnabled(),
			DCF:        config.IsDCFEnabled(),
			TimeFormat: config.TimeFormat(),
			TempFormat: config.TempFormat(),
			Sensors:    sensors,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func GetConfigBySensor(repo repository.Repository) http.HandlerFunc {

	type response struct {
		sensorConfig
		UpdatedAt string `json:"updatedAt"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		config := repo.Config()
		if config == nil {
			return
		}

		var sensorID int
		if sensorID, err := strconv.Atoi(chi.URLParam(r, "sensor")); err != nil ||
			sensorID < config.SensorMin() ||
			sensorID > config.SensorMax() {
			panic("invalid sensor id")
		}

		response := response{
			UpdatedAt:    repo.LastConfigUpdate().Format(time.RFC3339),
			sensorConfig: newSensorConfig(sensorID, config),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

type Traceable interface {
	StartTracing(io.Writer)
	StopTracing()
}

func GetTransceiverTrace(t Traceable) http.HandlerFunc {
	transceiver := t

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		s, err := strconv.Atoi(r.URL.Query().Get("seconds"))
		if err != nil || s == 0 {
			s = 5
		}

		b := strings.Builder{}
		b.WriteString("{ \"trace\": [")

		transceiver.StartTracing(&b)
		time.Sleep(time.Duration(s) * time.Second)
		transceiver.StopTracing()

		b.WriteString("{}]}")

		strings.NewReplacer("\n", ",").WriteString(w, b.String())

	}
}
