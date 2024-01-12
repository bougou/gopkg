package regpat

import (
	"regexp"
)

// GetPatternCapturedMap returns a map holding capture group information for the input string
// by using sepecified regexp pattern.
// The capture group name is stored as map key, the matched group value is stored as map value of the key.
//
// The first parameter is normally created by  `reg := regexp.MustCompile(somePatternStr)`
func GetPatternCapturedMap(reg *regexp.Regexp, input string) map[string]string {
	match := reg.FindStringSubmatch(input)

	result := make(map[string]string)

	if len(match) == 0 {
		return result
	}

	for i, name := range reg.SubexpNames() {
		if i == 0 {
			continue
		}

		if i >= len(match) {
			break
		}

		result[name] = match[i]
	}

	return result
}

// GetPatternCaptured returns a matched capture group string content for the input string
// by using specified regexp pattern.
// Return empty string if not matched or if the capture group name does not exist.
func GetPatternCaptured(reg *regexp.Regexp, input string, captureGroup string) string {
	m := GetPatternCapturedMap(reg, input)
	s, ok := m[captureGroup]
	if ok {
		return s
	}
	return ""
}
