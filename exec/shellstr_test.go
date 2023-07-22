package exec

import (
	"fmt"
	"testing"
)

func Test_Str2Cmd(t *testing.T) {
	tests := []struct {
		name     string
		cmdstr   string
		expected []string
	}{
		{
			name:     "case 1",
			cmdstr:   "ls -al",
			expected: []string{"ls", "-al"},
		},

		{
			name:     "case 2",
			cmdstr:   `echo "hello world"`,
			expected: []string{"echo", "hello world"},
		},

		{
			name:     "case 3",
			cmdstr:   `echo "Test" "$Hello world"`,
			expected: []string{"echo", `"Test"`, `$Hello world`},
		},
		{
			name:     "case 4",
			cmdstr:   `ls -al ~/output`,
			expected: []string{"ls", "-al", "~/output"},
		},
	}

	compareStrListFn := func(a []string, b []string) error {
		if len(a) != len(b) {
			return fmt.Errorf("length not matched, left: %d, right: %d", len(a), len(b))
		}

		for i := range a {
			if a[i] != b[i] {
				return fmt.Errorf("element at index (%d) not equal, left: %s, right: %s", i, a[i], b[i])
			}
		}
		return nil
	}

	for _, tt := range tests {

		got := ShellStr2List(tt.cmdstr)

		if err := compareStrListFn(got, tt.expected); err != nil {
			t.Errorf("test %s failed, got: %v, want: %v, reason: %v", tt.name, got, tt.expected, err)
		}
	}
}
