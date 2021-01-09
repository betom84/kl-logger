package models

import (
	"time"

	"github.com/betom84/kl-logger/utils"
)

type SendTimeData struct {
	CfgChecksum uint16
	Now         time.Time
}

func (d SendTimeData) MarshalBinary() ([]byte, error) {
	weekday := d.Now.Weekday()
	if weekday == 0 {
		weekday = 7
	}

	hexDay := utils.ToHexNumber(uint(d.Now.Day()))
	hexMonth := utils.ToHexNumber(uint(d.Now.Month()))
	hexYear := utils.ToHexNumber(uint(d.Now.Year() - 2000))

	data := make([]byte, 9)

	data[0] = byte(d.CfgChecksum >> 8)
	data[1] = byte(d.CfgChecksum)
	data[2] = utils.ToHexNumber(uint(d.Now.Second()))
	data[3] = utils.ToHexNumber(uint(d.Now.Minute()))
	data[4] = utils.ToHexNumber(uint(d.Now.Hour()))
	data[5] = (hexDay << 4) | byte(weekday)
	data[6] = (hexMonth << 4) | (hexDay >> 4)
	data[7] = (hexYear << 4) | (hexMonth >> 4)
	data[8] = (hexYear >> 4)

	return data, nil
}
