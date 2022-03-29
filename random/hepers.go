package random

import "math/rand"

// rand.Seed(time.Now().UnixNano())

func randIntRange(r *rand.Rand, min, max int) int {
	if min == max {
		return min
	}
	return r.Intn((max+1)-min) + min
}

// Generate random ASCII digit
func randDigit(r *rand.Rand) rune {
	return rune(byte(r.Intn(10)) + '0')
}

func number(r *rand.Rand, min int, max int) int { return randIntRange(r, min, max) }

func Bytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}
