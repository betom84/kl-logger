package klimalogg

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/betom84/kl-logger/transceiver"
	"github.com/betom84/kl-logger/transceiver/ax5051"
	"github.com/betom84/kl-logger/utils"
	"github.com/sirupsen/logrus"
)

// Console represents a klimalogg console
type Console struct {
	LoggerID             uint8
	transceiver          *transceiver.Transceiver
	handler              IOHandler
	stopCommunication    chan bool
	communicationRunning bool
}

// NewConsole using given transceiver
func NewConsole(t *transceiver.Transceiver) *Console {
	return &Console{transceiver: t, handler: IOHandler{}, LoggerID: 0}
}

// Initialise default klimalogg console (868MHz)
func (c *Console) Initialise(repository Repository) error {
	deviceID, serial, err := c.getTransceiverInfo()
	if err != nil {
		return fmt.Errorf("failed to read transceiver info; %v", err)
	}

	c.handler.DeviceID = deviceID
	c.handler.LoggerID = c.LoggerID
	c.handler.Repostitory = repository

	logrus.WithFields(logrus.Fields{
		"deviceID": fmt.Sprintf("%d (%04x)", c.handler.DeviceID, c.handler.DeviceID),
		"loggerID": c.handler.LoggerID,
		"serial":   serial,
	}).Debug("initialise klimalogg console")

	err = c.CorrectFrequency(float64(868300000))
	if err != nil {
		return fmt.Errorf("failed to correct frequency; %v", err)
	}

	return nil
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

			c.processRequestFrame()

			//nextReadinessProbe = time.NewTicker(75 * time.Millisecond)
		}
	}()

	c.communicationRunning = true
	return nil
}

func (c Console) processRequestFrame() {
	request := transceiver.Frame{}
	err := c.transceiver.GetFrame(&request)
	if err != nil {
		logrus.WithError(err).Error("unable to get request frame")

		return
	}

	logrus.WithFields(logrus.Fields{
		"typeID":   fmt.Sprintf("%02x", request.TypeID),
		"data":     utils.Prettify(request.Data),
		"deviceID": fmt.Sprintf("%04x", request.DeviceID),
		"loggerID": request.LoggerID,
	}).Debug("handle frame")

	response, err := c.handler.HandleRequest(request)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"request": request,
		}).Error("failed to handle request frame")

		return
	}

	logrus.WithFields(logrus.Fields{
		"typeID":   fmt.Sprintf("%02x", response.TypeID),
		"data":     utils.Prettify(response.Data),
		"deviceID": fmt.Sprintf("%04x", response.DeviceID),
		"loggerID": response.LoggerID,
	}).Debug("set response frame")

	err = c.transceiver.SetFrame(response)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err,
			"request":  request,
			"response": response,
		}).Error("unable to set response frame")

		return
	}

	err = c.transceiver.Set(transceiver.SetTX, []byte{})
	if err != nil {
		logrus.WithError(err).Error("unable to set TX")

		return
	}
}

func (c Console) prepareCommunication() {
	c.transceiver.Set(transceiver.Execute, []byte{0x05})
	c.transceiver.Set(transceiver.SetPreamblePattern, []byte{0xaa})
	c.transceiver.Set(transceiver.SetState, []byte{0x00})

	time.Sleep(1 * time.Second)
	c.transceiver.Set(transceiver.SetRX, []byte{})

	c.transceiver.Set(transceiver.SetPreamblePattern, []byte{0xaa})
	c.transceiver.Set(transceiver.SetState, []byte{0x1e})

	time.Sleep(1 * time.Second)
	c.transceiver.Set(transceiver.SetRX, []byte{})
}

// CorrectFrequency of transceiver to communicate with klimalogg console based on radio frequency band (EU/US)
func (c Console) CorrectFrequency(baseFrequency float64) error {
	var initFreq uint32 = uint32(math.Floor(baseFrequency / float64(16000000) * float64(16777216)))

	corr, err := c.transceiver.ReadConfigFlash(transceiver.FrequencyCorrection)
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
		err := c.transceiver.Set(transceiver.WriteRegister, []byte{(k & 0x7f), 0x01, v, 0x00})
		if err != nil {
			break
		}
	}

	return err
}

func (c Console) getTransceiverInfo() (uint16, string, error) {
	info, err := c.transceiver.ReadConfigFlash(transceiver.TransceiverInfo)
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
