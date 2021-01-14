package ax5051

// Registernames
const (
	REVISION     = 0x0
	SCRATCH      = 0x1
	POWERMODE    = 0x2
	XTALOSC      = 0x3
	FIFOCTRL     = 0x4
	FIFODATA     = 0x5
	IRQMASK      = 0x6
	IFMODE       = 0x8
	PINCFG1      = 0x0C
	PINCFG2      = 0x0D
	MODULATION   = 0x10
	ENCODING     = 0x11
	FRAMING      = 0x12
	CRCINIT3     = 0x14
	CRCINIT2     = 0x15
	CRCINIT1     = 0x16
	CRCINIT0     = 0x17
	FREQ3        = 0x20
	FREQ2        = 0x21
	FREQ1        = 0x22
	FREQ0        = 0x23
	FSKDEV2      = 0x25
	FSKDEV1      = 0x26
	FSKDEV0      = 0x27
	IFFREQHI     = 0x28
	IFFREQLO     = 0x29
	PLLLOOP      = 0x2C
	PLLRANGING   = 0x2D
	PLLRNGCLK    = 0x2E
	TXPWR        = 0x30
	TXRATEHI     = 0x31
	TXRATEMID    = 0x32
	TXRATELO     = 0x33
	MODMISC      = 0x34
	FIFOCONTROL2 = 0x37
	ADCMISC      = 0x38
	AGCTARGET    = 0x39
	AGCATTACK    = 0x3A
	AGCDECAY     = 0x3B
	AGCCOUNTER   = 0x3C
	CICDEC       = 0x3F
	DATARATEHI   = 0x40
	DATARATELO   = 0x41
	TMGGAINHI    = 0x42
	TMGGAINLO    = 0x43
	PHASEGAIN    = 0x44
	FREQGAIN     = 0x45
	FREQGAIN2    = 0x46
	AMPLGAIN     = 0x47
	TRKFREQHI    = 0x4C
	TRKFREQLO    = 0x4D
	XTALCAP      = 0x4F
	SPAREOUT     = 0x60
	TESTOBS      = 0x68
	APEOVER      = 0x70
	TMMUX        = 0x71
	PLLVCOI      = 0x72
	PLLCPEN      = 0x73
	PLLRNGMISC   = 0x74
	AGCMANUAL    = 0x78
	ADCDCLEVEL   = 0x79
	RFMISC       = 0x7A
	TXDRIVER     = 0x7B
	REF          = 0x7C
	RXMISC       = 0x7D
)

var registerDefaultValues = map[byte]byte{
	IFMODE:     0x00,
	MODULATION: 0x41,
	ENCODING:   0x07,
	FRAMING:    0x84,
	CRCINIT3:   0xff,
	CRCINIT2:   0xff,
	CRCINIT1:   0xff,
	CRCINIT0:   0xff,
	FREQ3:      0x38,
	FREQ2:      0x90,
	FREQ1:      0x00,
	FREQ0:      0x01,
	PLLLOOP:    0x1d,
	PLLRANGING: 0x08,
	PLLRNGCLK:  0x03,
	MODMISC:    0x03,
	SPAREOUT:   0x00,
	TESTOBS:    0x00,
	APEOVER:    0x00,
	TMMUX:      0x00,
	PLLVCOI:    0x01,
	PLLCPEN:    0x01,
	RFMISC:     0xb0,
	REF:        0x23,
	IFFREQHI:   0x20,
	IFFREQLO:   0x00,
	ADCMISC:    0x01,
	AGCTARGET:  0x0e,
	AGCATTACK:  0x11,
	AGCDECAY:   0x0e,
	CICDEC:     0x3f,
	DATARATEHI: 0x19,
	DATARATELO: 0x66,
	TMGGAINHI:  0x01,
	TMGGAINLO:  0x96,
	PHASEGAIN:  0x03,
	FREQGAIN:   0x04,
	FREQGAIN2:  0x0a,
	AMPLGAIN:   0x06,
	AGCMANUAL:  0x00,
	ADCDCLEVEL: 0x10,
	RXMISC:     0x35,
	FSKDEV2:    0x00,
	FSKDEV1:    0x31,
	FSKDEV0:    0x27,
	TXPWR:      0x03,
	TXRATEHI:   0x00,
	TXRATEMID:  0x51,
	TXRATELO:   0xec,
	TXDRIVER:   0x88,
}

// RegisterDefaultValues to initialize transceiver with
func RegisterDefaultValues() map[byte]byte {
	copy := make(map[byte]byte, len(registerDefaultValues))
	for k, v := range registerDefaultValues {
		copy[k] = v
	}

	return copy
}
