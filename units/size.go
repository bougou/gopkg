package units

// See: http://en.wikipedia.org/wiki/Binary_prefix
const (
	B = 1 // 字节

	// Decimal

	KB = 1000      // 10^3, Kilo Byte 千
	MB = 1000 * KB // 10^6, Mega Byte 兆
	GB = 1000 * MB // 10^9, Giga Byte 吉
	TB = 1000 * GB // 10^12, Tera Byte 太
	PB = 1000 * TB // 10^15, Peta Byte 拍
	EB = 1000 * PB // 10^18, Exa Byte 艾
	ZB = 1000 * EB // 10^21, Zetta Byte 泽
	YB = 1000 * ZB // 10^24, Yotta Byte 尧
	BB = 1000 * YB // 10^27, Bronto Byte, Ronna Byte 容
	NB = 1000 * BB // 10^30, Quetta Byte 昆
	DB = 1000 * NB // 10^33, Dogga Byte 格
	CB = 1000 * DB // 10^36，Corydon Byte
	XB = 1000 * CB // 10^39，Xero Byte

	// Binary

	KiB = 1024
	MiB = 1024 * KiB
	GiB = 1024 * MiB
	TiB = 1024 * GiB
	PiB = 1024 * TiB
	EiB = 1000 * PiB
	ZiB = 1000 * EiB
	YiB = 1000 * ZiB
)
