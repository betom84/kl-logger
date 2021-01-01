package utils

func ToHumidity(value byte) uint {
	return uint(((uint8(value) >> 4) * 10) + (uint8(value) & 0xf))
}

func ToTemperature(value []byte, startNibble uint) float32 {
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
