package frames

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	// HighNibble flag to address bits 2^4 to 2^7 of a byte
	HighNibble = 1

	// LowNibble flag to address bits 2^0 to 2^3 of a byte
	LowNibble = 2
)

func toInt(value byte, nibbles ...int) uint {
	high := uint8(value) >> 4
	low := uint8(value) & 0xf

	switch {
	case len(nibbles) == 0:
		fallthrough
	case len(nibbles) == 2 && nibbles[0] == HighNibble && nibbles[1] == LowNibble:
		return uint((high * 10) + low)
	case len(nibbles) == 1 && nibbles[0] == HighNibble:
		return uint(high)
	case len(nibbles) == 1 && nibbles[0] == LowNibble:
		return uint(low)
	}

	return 0
}

func toHumidity(value byte) uint {
	return toInt(value, 1, 2)
}

func toTemperature(value []byte, startNibble int) float32 {
	var t float32

	if startNibble == HighNibble {
		t = float32(value[0]>>4) * 10
		t += float32(value[0] & 0xf)
		t += float32(value[1]>>4) / 10
	} else {
		t = float32(value[0]&0xf) * 10
		t += float32(value[1] >> 4)
		t += float32(value[1]&0xf) / 10
	}

	return t - 40
}

// ToHexNumber creates a byte representation of a number without actually converting it.
// E.g. number 25 will be transformed to 0x25 instead of 0x19. That only works for numbers >0 and <100.
func toHexNumber(number uint) byte {
	if number >= 100 {
		return 0x00
	}

	return byte((number % 10) + uint(16*math.Floor(float64(number)/10)))
}

func toDateTime(dt []byte, startNibble uint) time.Time {
	var year, month, day uint
	var t1, t2, t3 uint

	switch {
	case startNibble == HighNibble && len(dt) >= 4:
		year = 2000 + toInt(dt[0])
		month = toInt(dt[1], 1)
		day = (toInt(dt[1], 2) * 10) + toInt(dt[2], 1)

		t1 = toInt(dt[2], 2)
		t2 = toInt(dt[3], 1)
		t3 = toInt(dt[3], 2)

	case startNibble == LowNibble && len(dt) >= 5:
		year = 2000 + (toInt(dt[0], 2)*10 + toInt(dt[1], 1))
		month = toInt(dt[1], 2)
		day = toInt(dt[2])

		t1 = toInt(dt[3], 1)
		t2 = toInt(dt[3], 2)
		t3 = toInt(dt[4], 1)
	}

	var hours, minutes uint
	if t1 < 10 {
		hours = t1
	} else if t1 >= 10 {
		hours = 10 + t1
	}

	if t2 < 10 {
		minutes = t2
	} else if t2 >= 10 {
		hours += 10
		minutes = (t2 - 10) * 10
	}

	minutes += t3

	t, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%04d-%02d-%02d %02d:%02d", year, month, day, hours, minutes))
	if err != nil {
		t = time.Time{}
	}

	return t
}

var charmap = []rune{
	' ', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'0', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I',
	'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S',
	'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '-', '+', '(',
	')', 'o', '*', ',', '/', '\\', ' ', '.', ' ', ' ',
	' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
	' ', ' ', ' ', '@',
}

func toString(b [8]byte) string {
	s := strings.Builder{}

	s.WriteString(toSubString(b[6:8], LowNibble))
	s.WriteString(toSubString(b[5:7], HighNibble))
	s.WriteString(toSubString(b[3:5], LowNibble))
	s.WriteString(toSubString(b[2:4], HighNibble))
	s.WriteString(toSubString(b[0:3], LowNibble))

	return strings.TrimSpace(s.String())
}

func toSubString(b []byte, startNibble int) string {
	var idx1, idx2 int

	if startNibble == HighNibble {
		idx1 = int((b[1]>>2)&0x3C) + int((b[0]>>2)&0x3)
		idx2 = int((b[0]<<4)&0x30) + int((b[0]>>4)&0xF)
	} else if startNibble == LowNibble {
		idx1 = int((b[1]<<2)&0x3C) + int((b[1]>>6)&0x3)
		idx2 = int(b[1]&0x30) + int(b[0]&0xF)
	}

	return fmt.Sprintf("%c%c", charmap[idx1], charmap[idx2])
}

func fromString(str string) [8]byte {
	chars := "!" + string(charmap[1:])

	s := strings.ReplaceAll(str, " ", "!")
	for i := len(s); i < 10; i++ {
		s += "!"
	}

	cIndices := make([]byte, len(s))
	for i, c := range s {
		cIndices[i] = byte(strings.IndexRune(chars, c))
	}

	result := [8]byte{}
	result[7] = ((cIndices[0] << 6) & 0xC0) + (cIndices[1] & 0x30) + ((cIndices[0] >> 2) & 0x0F)
	result[6] = ((cIndices[2] << 2) & 0xF0) + (cIndices[1] & 0x0F)
	result[5] = ((cIndices[3] << 4) & 0xF0) + ((cIndices[2] << 2) & 0x0C) + ((cIndices[3] >> 4) & 0x03)
	result[4] = ((cIndices[4] << 6) & 0xC0) + (cIndices[5] & 0x30) + ((cIndices[4] >> 2) & 0x0F)
	result[3] = ((cIndices[6] << 2) & 0xF0) + (cIndices[5] & 0x0F)
	result[2] = ((cIndices[7] << 4) & 0xF0) + ((cIndices[6] << 2) & 0x0C) + ((cIndices[7] >> 4) & 0x03)
	result[1] = ((cIndices[8] << 6) & 0xC0) + (cIndices[9] & 0x30) + ((cIndices[8] >> 2) & 0x0F)
	result[0] = (cIndices[9] & 0x0F)

	return result
}
