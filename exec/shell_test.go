package exec

import (
	"testing"
)

func TestShellCommandExec(t *testing.T) {
	cases := []struct {
		name string

		command     string
		silent      bool
		trim        bool
		expect      string
		expectError bool
		noexpect    bool
	}{
		{
			name:    "test1",
			command: "echo testing-output",
			trim:    false,
			expect:  "testing-output\n",
		},
		{
			name:    "test2",
			command: "echo testing-output",
			silent:  true,
			trim:    true,
			expect:  "testing-output",
		},
		{
			name:        "test3",
			command:     "command-that-does-not-exist",
			expectError: true,
		},
		{
			name:    "test4",
			command: `printf one\ntwo\nthree`,
			silent:  true,
			trim:    true,
			expect:  "one\ntwo\nthree",
		},
		{
			name:     "test5",
			command:  `grep -r 'some text with spaces' .`,
			silent:   true,
			trim:     false,
			noexpect: true,
		},
		{
			name:     "test6",
			command:  `grep -r "some text with spaces" .`,
			trim:     false,
			noexpect: true,
		},
	}

	for _, tt := range cases {
		sc := ShellCommand{
			command: tt.command,
			silent:  tt.silent,
			trim:    tt.trim,
		}

		output, err := sc.Exec()
		if err != nil {
			if !tt.expectError {
				t.Errorf("command error, not expect error")
			}
		}

		if tt.noexpect {
			continue
		}

		if !isByteSliceEqual(output, []byte(tt.expect)) {
			t.Errorf("test case: %s failed, expect: '%s', got: '%s'", tt.name, tt.expect, output)
		}
	}

}

func isByteSliceEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
