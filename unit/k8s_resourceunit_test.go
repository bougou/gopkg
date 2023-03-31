package unit

import (
	"testing"
)

func Test_ParseK8SResourceStrToFloat64(t *testing.T) {
	tests := []struct {
		input  string
		expect float64
	}{
		{
			input:  "100m",
			expect: 0.1,
		},
		{
			input:  "1024Mi",
			expect: 1073741824,
		},
		{
			input:  "1.2K",
			expect: 1200,
		},
	}

	for _, tt := range tests {
		got := ParseK8SResourceStrToFloat64(tt.input)
		if got != tt.expect {
			t.Errorf("not match, expect: %f， got: %f", tt.expect, got)
		}
	}
}

func TestParseK8SResourceFloat64ToStr(t *testing.T) {
	tests := []struct {
		input  float64
		expect string
	}{
		{
			input:  0.1,
			expect: "100m",
		},
		{
			input:  2048,
			expect: "2Ki",
		},
		{
			input:  568,
			expect: "568",
		},
		{
			input:  0.568,
			expect: "568m",
		},
		{
			input:  1025,
			expect: "1Ki",
		},
		{
			input:  1024 * 1024 * 1024, // 1073741824
			expect: "1Gi",
		},
		{
			input:  1024*1024*1024 + 1024*1024*1024*0.049, // 1G + 0.049G
			expect: "1Gi",
		},
		{
			input:  1024*1024*1024 + 1024*1024*1024*0.05, // 1G + 0.05G, 0.05 是阈值
			expect: "2Gi",
		},
		{
			input:  1024*1024*1024 + 1024*1024*1024*0.051, // 1G + 0.051G
			expect: "2Gi",
		},
	}

	for i, tt := range tests {
		got := ParseK8SResourceFloat64ToStr(tt.input)
		if got != tt.expect {
			t.Errorf("test %d failed, expect: %s， got: %s", i, tt.expect, got)
		}
	}
}
