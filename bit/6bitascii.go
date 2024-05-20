package bit

func Decode6bitASCII(raw []byte) []byte {
	// 6-bit ASCII definition
	var ascci6bit = [64]byte{
		' ', '!', '"', '#', '$', '%', '&', '\'',
		'(', ')', '*', '+', ',', '-', '.', '/',
		'0', '1', '2', '3', '4', '5', '6', '7',
		'8', '9', ':', ';', '<', '=', '>', '?',
		'@', 'A', 'B', 'C', 'D', 'E', 'F', 'G',
		'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
		'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W',
		'X', 'Y', 'Z', '[', '\\', ']', '^', '_',
	}

	var leftover byte
	var s []byte

	for i := 0; i < len(raw); i++ {
		// every 3 bytes pack 4 chars, so we can calculate
		// character positions in a byte based on the remainder of division by 3.
		switch i % 3 {
		case 0:
			idx := raw[i] & 0x3f            // 6 right bits are an index of one char
			leftover = (raw[i] & 0xc0) >> 6 // 2 left bits are leftovers

			s = append(s, ascci6bit[idx])
		case 1:
			idx := leftover | (raw[i]&0x0f)<<2 // index of one char is 2-bit leftover as prefix plus 4 right bits
			leftover = (raw[i] & 0xf0) >> 4    // 4 left bits are leftovers

			s = append(s, ascci6bit[idx])
		case 2:
			idx := (raw[i]&0x03)<<4 | leftover // index of one char is 2 right bits plus 4-bit leftover as suffix
			leftover = 0                       // cleanup leftover calculation

			s = append(s, ascci6bit[idx])

			idx = (raw[i] & 0xfc) >> 2 // 6 left bits are an index of one char
			s = append(s, ascci6bit[idx])
		}
	}

	return s
}
