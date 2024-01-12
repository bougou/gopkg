package regpat

import (
	"regexp"
	"testing"
)

func Test_NumberInteger(t *testing.T) {

	tests := []struct {
		input  string
		expect bool
	}{

		{
			input:  "123",
			expect: true,
		},
		{
			input:  "123.456",
			expect: false,
		},
		{
			input:  ".456",
			expect: false,
		},
		{
			input:  "123.", // note
			expect: false,
		},
	}

	for _, tt := range tests {
		pattern := regexp.MustCompile("^" + NumberInteger + "$")
		matched := pattern.Match([]byte(tt.input))
		if matched != tt.expect {
			t.Errorf("match for input (%s) error: result (%v), expected (%v)", tt.input, matched, tt.expect)
		}
	}
}

func Test_NumberFloat(t *testing.T) {

	tests := []struct {
		input  string
		expect bool
	}{

		{
			input:  "123",
			expect: false,
		},
		{
			input:  "123.456",
			expect: true,
		},
		{
			input:  ".456",
			expect: false,
		},
		{
			input:  "123.", // note
			expect: false,
		},
	}

	for _, tt := range tests {
		pattern := regexp.MustCompile("^" + NumberFloat + "$")
		matched := pattern.Match([]byte(tt.input))
		if matched != tt.expect {
			t.Errorf("match for input (%s) error: result (%v), expected (%v)", tt.input, matched, tt.expect)
		}
	}
}

func Test_Number(t *testing.T) {

	tests := []struct {
		input  string
		expect bool
	}{

		{
			input:  "123",
			expect: true,
		},
		{
			input:  "123.456",
			expect: true,
		},
		{
			input:  ".456",
			expect: false,
		},
		{
			input:  "123.", // note
			expect: false,
		},
	}

	for _, tt := range tests {
		pattern := regexp.MustCompile("^" + Number + "$")
		matched := pattern.Match([]byte(tt.input))
		if matched != tt.expect {
			t.Errorf("match for input (%s) error: result (%v), expected (%v)", tt.input, matched, tt.expect)
		}
	}
}

func Test_NumberAllowNoInteger(t *testing.T) {

	tests := []struct {
		input  string
		expect bool
	}{

		{
			input:  "123",
			expect: true,
		},
		{
			input:  "123.456",
			expect: true,
		},
		{
			input:  ".456",
			expect: true,
		},
		{
			input:  "123.", // note
			expect: false,
		},
	}

	for _, tt := range tests {
		pattern := regexp.MustCompile("^" + NumberAllowNoInteger + "$")
		matched := pattern.Match([]byte(tt.input))
		if matched != tt.expect {
			t.Errorf("match for input (%s) error: result (%v), expected (%v)", tt.input, matched, tt.expect)
		}
	}
}

func Test_NumberAllowNoDecimal(t *testing.T) {

	tests := []struct {
		input  string
		expect bool
	}{

		{
			input:  "123",
			expect: true,
		},
		{
			input:  "123.456",
			expect: true,
		},
		{
			input:  ".456",
			expect: true,
		},
		{
			input:  "123.", // note
			expect: true,
		},
	}

	for _, tt := range tests {
		pattern := regexp.MustCompile("^" + NumberAllowNoDecimal + "$")
		matched := pattern.Match([]byte(tt.input))
		if matched != tt.expect {
			t.Errorf("match for input (%s) error: result (%v), expected (%v)", tt.input, matched, tt.expect)
		}
	}
}
