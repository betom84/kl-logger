package frames

import "math"

func toHumidity(value byte) uint {
	return uint(((uint8(value) >> 4) * 10) + (uint8(value) & 0xf))
}

func toTemperature(value []byte, startNibble uint) float32 {
	var t float32

	if startNibble == 1 {
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
