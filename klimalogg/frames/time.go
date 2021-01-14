package frames

import (
	"time"

	"github.com/betom84/kl-logger/utils"
)

// SendTimeFrame to send time to klimalogg console
type SendTimeFrame struct {
	SetFrame
}

// NewSendTimeFrame to set time now
func NewSendTimeFrame() SendTimeFrame {
	f := NewSetFrame()

	f.SetTypeID(SendTime)
	f.Write(make([]byte, 9))

	return SendTimeFrame{f}
}

func (f SendTimeFrame) SetCfgChecksum(checksum int) {
	f.SetFrame[7] = byte(checksum >> 8)
	f.SetFrame[8] = byte(checksum)
}

func (f SendTimeFrame) SetTime(t time.Time) {
	weekday := t.Weekday()
	if weekday == 0 {
		weekday = 7
	}

	hexDay := utils.ToHexNumber(uint(t.Day()))
	hexMonth := utils.ToHexNumber(uint(t.Month()))
	hexYear := utils.ToHexNumber(uint(t.Year() - 2000))

	f.SetFrame[9] = utils.ToHexNumber(uint(t.Second()))
	f.SetFrame[10] = utils.ToHexNumber(uint(t.Minute()))
	f.SetFrame[11] = utils.ToHexNumber(uint(t.Hour()))
	f.SetFrame[12] = (hexDay << 4) | byte(weekday)
	f.SetFrame[13] = (hexMonth << 4) | (hexDay >> 4)
	f.SetFrame[14] = (hexYear << 4) | (hexMonth >> 4)
	f.SetFrame[15] = (hexYear >> 4)
}
