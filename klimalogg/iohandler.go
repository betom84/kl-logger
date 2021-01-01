package klimalogg

import (
	"fmt"

	"github.com/betom84/kl-logger/klimalogg/models"
	"github.com/betom84/kl-logger/transceiver"
	"github.com/sirupsen/logrus"
)

// Repository to update/store received klimalogg data
type Repository interface {
	UpdateWeather(float32, uint)
	UpdateSettings(int, bool, bool, string, string)
}

// IOHandler processes frames received from klimalogg console
type IOHandler struct {
	LoggerID    uint8
	DeviceID    uint16
	Repostitory Repository
}

// HandleRequest from klimalogg console and create next frame to be send to klimalogg console
func (h IOHandler) HandleRequest(curr transceiver.Frame) (*transceiver.Frame, error) {
	var response *transceiver.Frame
	var err error

	switch curr.TypeID {
	case transceiver.CurrentWeatherResponse:
		response, err = h.handleCurrentWeather(curr)
	case transceiver.ConfigResponse:
		response, err = h.handleConfig(curr)
	case transceiver.RequestFirstConfigResponse:
		response, err = h.handleFirstConfig(curr)
	}

	if err != nil || response != nil {
		return response, err
	}

	next := &transceiver.Frame{}
	next.TypeID = transceiver.CurrentWeatherRequest
	next.DeviceID = curr.DeviceID
	next.LoggerID = curr.LoggerID

	data := models.CurrentWeatherRequestData{}
	next.SetData(data)

	return next, nil
}

func (h IOHandler) handleCurrentWeather(curr transceiver.Frame) (*transceiver.Frame, error) {
	data := models.CurrentWeatherResponseData{}
	err := data.UnmarshalBinary(curr.Data)
	if err != nil {
		return nil, err
	}

	if h.Repostitory != nil {
		h.Repostitory.UpdateWeather(data.Temperature[0], data.Humidity[0])
	} else {
		logrus.WithFields(logrus.Fields{
			"temperature": data.Temperature,
			"humidity":    data.Humidity,
			"cfgChecksum": data.CfgChecksum,
		}).Info("received current weather")
	}

	return nil, nil
}

func (h IOHandler) handleConfig(curr transceiver.Frame) (*transceiver.Frame, error) {
	data := models.ConfigResponseData{}
	err := data.UnmarshalBinary(curr.Data)
	if err != nil {
		return nil, err
	}

	if h.Repostitory != nil {
		s := data.Settings
		h.Repostitory.UpdateSettings(s.Contrast, s.Alert, s.DCF, s.TimeFormat, s.TempFormat)
	} else {
		logrus.WithFields(logrus.Fields{
			"cfgChecksum": data.CfgChecksum,
			"settings":    fmt.Sprintf("%+v", data.Settings),
		}).Info("received config")
	}

	return nil, nil
}

func (h IOHandler) handleFirstConfig(curr transceiver.Frame) (*transceiver.Frame, error) {
	next := &transceiver.Frame{}
	next.DeviceID = curr.DeviceID
	next.LoggerID = curr.LoggerID
	next.TypeID = transceiver.ConfigRequest

	data := models.FirstConfigRequestData{
		DeviceID: h.DeviceID,
		LoggerID: h.LoggerID,
	}

	next.SetData(data)

	return next, nil
}
