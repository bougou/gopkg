package strutil

func IsUpperLetter(c int32) bool {
	return c >= 'A' && c <= 'Z'
}

func IsLowerLetter(c int32) bool {
	return c >= 'a' && c <= 'z'
}
