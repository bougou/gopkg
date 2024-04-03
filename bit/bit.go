package bit

func SetBit7(b uint8) uint8 {
	return b | 0x80
}

func SetBit6(b uint8) uint8 {
	return b | 0x40
}

func SetBit5(b uint8) uint8 {
	return b | 0x20
}

func SetBit4(b uint8) uint8 {
	return b | 0x10
}

func SetBit3(b uint8) uint8 {
	return b | 0x08
}

func SetBit2(b uint8) uint8 {
	return b | 0x04
}

func SetBit1(b uint8) uint8 {
	return b | 0x02
}

func SetBit0(b uint8) uint8 {
	return b | 0x01
}

func ClearBit7(b uint8) uint8 {
	return b & 0x7f
}

func ClearBit6(b uint8) uint8 {
	return b & 0xbf
}

func ClearBit5(b uint8) uint8 {
	return b & 0xdf
}

func ClearBit4(b uint8) uint8 {
	return b & 0xef
}

func ClearBit3(b uint8) uint8 {
	return b & 0xf7
}

func ClearBit2(b uint8) uint8 {
	return b & 0xfb
}

func ClearBit1(b uint8) uint8 {
	return b & 0xfd
}

func ClearBit0(b uint8) uint8 {
	return b & 0xfe
}
func IsBit7Set(b uint8) bool {
	return b&0x80 == 0x80
}

func IsBit6Set(b uint8) bool {
	return b&0x40 == 0x40
}

func IsBit5Set(b uint8) bool {
	return b&0x20 == 0x20
}

func IsBit4Set(b uint8) bool {
	return b&0x10 == 0x10
}

func IsBit3Set(b uint8) bool {
	return b&0x08 == 0x08
}

func IsBit2Set(b uint8) bool {
	return b&0x04 == 0x04
}

func IsBit1Set(b uint8) bool {
	return b&0x02 == 0x02
}

func IsBit0Set(b uint8) bool {
	return b&0x01 == 0x01
}
