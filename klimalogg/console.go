package klimalogg

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/betom84/kl-logger/klimalogg/ax5051"
	"github.com/betom84/kl-logger/klimalogg/frames"
	"github.com/sirupsen/logrus"
)

// Console represents a klimalogg console
type Console struct {
	transceiver          *Transceiver
	stopCommunication    chan bool
	communicationRunning bool
	listeners            []chan<- interface{}
	cfgChecksum          uint16
	deviceID             uint16
	loggerID             uint8
	serial               string
}

// NewConsole using given transceiver (default frequency 868MHz)
func NewConsole(t *Transceiver) (*Console, error) {
	var err error

	c := &Console{}
	c.transceiver = t
	c.loggerID = 0
	c.listeners = make([]chan<- interface{}, 0)

	c.deviceID, c.serial, err = c.getTransceiverInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to read transceiver info; %v", err)
	}

	logrus.WithFields(logrus.Fields{
		"deviceID": fmt.Sprintf("%d (%04x)", c.deviceID, c.deviceID),
		"loggerID": 0,
		"serial":   c.serial,
	}).Debug("initialise klimalogg console")

	err = c.CorrectFrequency(float64(868300000))
	return c, err
}

// AddListener to receive weather- and config updates from klimalogg console
func (c *Console) AddListener(l chan<- interface{}) {
	c.listeners = append(c.listeners, l)
}

func (c *Console) notifyListeners(f interface{}) {
	for _, l := range c.listeners {
		l <- f
	}
}

// Close klimalogg console, stop running communication
func (c *Console) Close() {
	if c.communicationRunning {
		c.stopCommunication <- true
		close(c.stopCommunication)

		c.communicationRunning = false
	}
}

// StartCommunication with klimalogg console via transceiver (non-blocking)
func (c *Console) StartCommunication() error {
	if c.communicationRunning == true {
		return fmt.Errorf("communication already running")
	}

	c.prepareCommunication()

	c.stopCommunication = make(chan bool)

	go func() {
		nextReadinessProbe := time.NewTicker(1 * time.Second)

		for {
			select {
			case <-c.stopCommunication:
				logrus.Info("stop communication")
				return

			case <-nextReadinessProbe.C:
				if !c.transceiver.IsReady() {
					//nextReadinessProbe = time.NewTicker(5 * time.Millisecond)
					continue
				}
			}

			c.processCurrentFrame()

			//nextReadinessProbe = time.NewTicker(75 * time.Millisecond)
		}
	}()

	c.communicationRunning = true
	return nil
}

func (c Console) processCurrentFrame() {
	current := frames.NewGetFrame()
	err := c.transceiver.GetFrame(&current)
	if err != nil {
		logrus.WithError(err).Error("unable to get current frame")

		return
	}

	logrus.WithFields(logrus.Fields{
		"typeID":   current.TypeID(),
		"deviceID": fmt.Sprintf("%04x", current.DeviceID()),
		"loggerID": current.LoggerID(),
	}).Debug("handle frame")

	next, err := c.handleFrame(current)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"current": current,
		}).Error("failed to handle current frame")

		return
	}

	logrus.WithFields(logrus.Fields{
		"typeID":   next.TypeID(),
		"deviceID": fmt.Sprintf("%04x", next.DeviceID()),
		"loggerID": next.LoggerID(),
	}).Debug("set next frame")

	err = c.transceiver.SetFrame(next)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"current": current,
			"next":    next,
		}).Error("failed to set next frame")

		return
	}

	err = c.transceiver.Set(SetTX, []byte{})
	if err != nil {
		logrus.WithError(err).Error("failed to set TX")

		return
	}
}

func (c *Console) handleFrame(curr frames.GetFrame) (*frames.SetFrame, error) {
	var next *frames.SetFrame

	switch curr.TypeID() {
	case frames.CurrentWeatherResponse:
		c.notifyListeners(frames.CurrentWeatherResponseFrame{GetFrame: curr})

	case frames.ConfigResponse:
		f := frames.ConfigResponseFrame{GetFrame: curr}
		c.cfgChecksum = f.CfgChecksum()
		c.notifyListeners(f)

	case frames.RequestFirstConfigResponse:
		f := frames.NewFirstConfigRequestFrame()
		f.SetDeviceID(c.deviceID)
		f.SetLoggerID(c.loggerID)
		next = &f.SetFrame

	case frames.RequestTimeResponse:
		f := frames.NewSendTimeFrame()
		f.SetDeviceID(c.deviceID)
		f.SetLoggerID(c.loggerID)
		f.SetCfgChecksum(c.cfgChecksum)
		f.SetTime(time.Now())
		next = &f.SetFrame

	default:
		logrus.WithField("frame", curr).Warn("handle unsupported frame type")
	}

	if next != nil {
		return next, nil
	}

	currWeather := frames.NewCurrentWeatherRequestFrame()
	currWeather.SetCfgChecksum(c.cfgChecksum)
	currWeather.SetDeviceID(c.deviceID)
	currWeather.SetLoggerID(c.loggerID)

	return &currWeather.SetFrame, nil
}

func (c Console) prepareCommunication() {
	c.transceiver.Set(Execute, []byte{0x05})
	c.transceiver.Set(SetPreamblePattern, []byte{0xaa})
	c.transceiver.Set(SetState, []byte{0x00})

	time.Sleep(1 * time.Second)
	c.transceiver.Set(SetRX, []byte{})

	c.transceiver.Set(SetPreamblePattern, []byte{0xaa})
	c.transceiver.Set(SetState, []byte{0x1e})

	time.Sleep(1 * time.Second)
	c.transceiver.Set(SetRX, []byte{})
}

// CorrectFrequency of transceiver to communicate with klimalogg console based on radio frequency band (EU/US)
func (c Console) CorrectFrequency(baseFrequency float64) error {
	var initFreq uint32 = uint32(math.Floor(baseFrequency / float64(16000000) * float64(16777216)))

	corr, err := c.transceiver.ReadConfigFlash(FrequencyCorrection)
	if err != nil {
		return err
	}

	corrFreq := uint32(corr[0])<<24 | uint32(corr[1])<<16 | uint32(corr[2])<<8 | uint32(corr[3])

	frequency := initFreq + corrFreq
	if (frequency % 2) == 0 {
		frequency++
	}

	var register = ax5051.RegisterDefaultValues()
	register[ax5051.FREQ3] = byte((frequency >> 24) & 0xff)
	register[ax5051.FREQ2] = byte((frequency >> 16) & 0xff)
	register[ax5051.FREQ1] = byte((frequency >> 8) & 0xff)
	register[ax5051.FREQ0] = byte((frequency >> 0) & 0xff)

	logrus.WithFields(logrus.Fields{
		"initFreq":  initFreq,
		"corrFreq":  corrFreq,
		"frequency": frequency,
	}).Debug("correct frequency")

	for k, v := range register {
		err := c.transceiver.Set(WriteRegister, []byte{(k & 0x7f), 0x01, v, 0x00})
		if err != nil {
			break
		}
	}

	return err
}

func (c Console) getTransceiverInfo() (uint16, string, error) {
	info, err := c.transceiver.ReadConfigFlash(TransceiverInfo)
	if err != nil {
		return 0, "", err
	}

	deviceID := uint16(info[5])<<8 | uint16(info[6])

	serial := make([]string, len(info))
	for i := range info {
		serial[i] = fmt.Sprintf("%02d", info[i])
	}

	return deviceID, strings.Join(serial, ""), nil
}

// Transceiver used to communicate with klimalogg console
func (c Console) Transceiver() *Transceiver {
	return c.transceiver
}
