package regpat

import (
	"regexp"
	"testing"
)

func Test_GetPatternCaptured(t *testing.T) {
	tests := []struct {
		pattern     string // regexp pattern
		input       string // input string to match
		expectField string // capture group name
		expectValue string
	}{
		{
			pattern:     `^(?P<Number>[0-9.]+)(?P<Unit>m||[KMGTP]i?)$`,
			input:       "100m",
			expectField: "Number",
			expectValue: "100",
		},
		{
			pattern:     `^(?P<Number>[0-9.]+)(?P<Unit>m||[KMGTP]i?)$`,
			input:       "100m",
			expectField: "Unit",
			expectValue: "m",
		},
		{
			pattern:     `^(?P<Number>[0-9.]+)(?P<Unit>m||[KMGTP]i?)$`,
			input:       "4096Mi",
			expectField: "Number",
			expectValue: "4096",
		},
		{
			pattern:     `^(?P<Number>[0-9.]+)(?P<Unit>m||[KMGTP]i?)$`,
			input:       "4096Mi",
			expectField: "Unit",
			expectValue: "Mi",
		},
		{
			pattern:     `^(?P<Number>[0-9.]+)(?P<Unit>m||[KMGTP]i?)$`,
			input:       "4096M",
			expectField: "Unit",
			expectValue: "M",
		},
	}

	for i, tt := range tests {

		r := regexp.MustCompile(tt.pattern)

		got := GetPatternCaptured(r, tt.input, tt.expectField)
		if got != tt.expectValue {
			t.Errorf("test %d failed, expect: %s, got: %s", i, tt.expectValue, got)
		}
	}

}
