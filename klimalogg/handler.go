package klimalogg

import (
	"fmt"
	"time"

	"github.com/betom84/kl-logger/klimalogg/frames"
	"github.com/sirupsen/logrus"
)

// Repository to update/store received klimalogg data
type Repository interface {
	UpdateWeather(float32, uint)
	UpdateSettings(int, bool, bool, string, string)
}

// Handler processes frames received from klimalogg console
type Handler struct {
	LoggerID    uint8
	DeviceID    uint16
	Repostitory Repository

	cfgChecksum int
}

// HandleRequest from klimalogg console and create next frame to be send to klimalogg console
func (h Handler) HandleRequest(curr frames.GetFrame) (*frames.SetFrame, error) {
	var response *frames.SetFrame

	switch curr.TypeID() {
	case frames.CurrentWeatherResponse:
		h.handleCurrentWeather(frames.CurrentWeatherResponseFrame{GetFrame: curr})
	case frames.ConfigResponse:
		h.handleConfig(frames.ConfigResponseFrame{GetFrame: curr})
	case frames.RequestFirstConfigResponse:
		response = h.handleFirstConfig()
	case frames.RequestTimeResponse:
		response = h.handleTime()
	default:
		logrus.WithField("frame", curr).Warn("handle unsupported frame type")
	}

	if response != nil {
		return response, nil
	}

	next := frames.NewCurrentWeatherRequestFrame()
	next.SetCfgChecksum(h.cfgChecksum)
	next.SetDeviceID(h.DeviceID)
	next.SetLoggerID(int(h.LoggerID))

	return &next.SetFrame, nil
}

func (h Handler) handleCurrentWeather(f frames.CurrentWeatherResponseFrame) {
	if h.Repostitory != nil {
		h.Repostitory.UpdateWeather(f.Temperature(0), f.Humidity(0))
	} else {
		logrus.WithFields(logrus.Fields{
			"temperature": f.Temperature(0),
			"humidity":    f.Humidity(0),
			"cfgChecksum": f.CfgChecksum(),
		}).Info("received current weather")
	}
}

func (h Handler) handleConfig(f frames.ConfigResponseFrame) {
	h.cfgChecksum = f.CfgChecksum()

	if h.Repostitory != nil {
		s := f.Settings()
		h.Repostitory.UpdateSettings(s.Contrast, s.Alert, s.DCF, s.TimeFormat, s.TempFormat)
	} else {
		logrus.WithFields(logrus.Fields{
			"cfgChecksum": f.CfgChecksum(),
			"settings":    fmt.Sprintf("%+v", f.Settings()),
		}).Info("received config")
	}
}

func (h Handler) handleFirstConfig() *frames.SetFrame {
	next := frames.NewFirstConfigRequestFrame()
	next.SetDeviceID(int(h.DeviceID))
	next.SetLoggerID(int(h.LoggerID))

	return &next.SetFrame
}

func (h Handler) handleTime() *frames.SetFrame {
	next := frames.NewSendTimeFrame()
	next.SetDeviceID(h.DeviceID)
	next.SetLoggerID(int(h.LoggerID))
	next.SetCfgChecksum(h.cfgChecksum)
	next.SetTime(time.Now())

	return &next.SetFrame
}
