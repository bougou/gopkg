package exec

import (
	"fmt"
	"testing"
)

func Test_Str2Cmd(t *testing.T) {
	tests := []struct {
		name   string
		cmdstr string
		list   []string
	}{
		{
			name:   "case 1",
			cmdstr: "ls -al",
			list:   []string{"ls", "-al"},
		},

		{
			name:   "case 2",
			cmdstr: `echo "hello world"`,
			list:   []string{"echo", "hello world"},
		},

		{
			name:   "case 3",
			cmdstr: `echo "Test" "$Hello world"`,
			list:   []string{"echo", "Test", `$Hello world`},
		},
	}

	compareStrList := func(a []string, b []string) error {
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
		if err := compareStrList(got, tt.list); err != nil {
			t.Errorf("test %s failed, got: %v, want: %v, reason: %v", tt.name, got, tt.list, err)
		}
	}
}
