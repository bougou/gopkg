package bit

import (
	"fmt"
	"testing"
)

func TestPrint(t *testing.T) {
	tests := []struct {
		format string
		data   []byte
		expect string
	}{
		{"%20X", []byte("xyz"), "              78797A"},
		{"% 20X", []byte("xyz"), "            78 79 7A"},
		{"%#20X", []byte("xyz"), "            0X78797A"},
		{"%# 20X", []byte("xyz"), "      0X78 0X79 0X7A"},
		{"%-20X", []byte("xyz"), "78797A              "},
		{"% -20X", []byte("xyz"), "78 79 7A            "},
		{"%-#20X", []byte("xyz"), "0X78797A            "},
		{"%-# 20X", []byte("xyz"), "0X78 0X79 0X7A      "},

		{"%# 02X", []byte("xyz"), "0X78 0X79 0X7A"},
		{"%# 02x", []byte("xyz"), "0x78 0x79 0x7a"},
	}

	for _, tt := range tests {
		got := fmt.Sprintf(tt.format, tt.data)
		if got != tt.expect {
			t.Errorf("test not matched, expect: '%s', got: '%s'", tt.expect, got)
		}
	}

}
