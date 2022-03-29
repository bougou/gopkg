package exec

import (
	"strings"
	"testing"
)

func TestExec(t *testing.T) {
	sh := ShellCommand{
		command: "echo testing-output",
	}

	output, err := sh.Exec()
	if err != nil {
		t.Errorf("Command failed: %v", err)
	}
	if output != "testing-output" {
		t.Errorf("want: 'testing-output', got: '%s'", output)
	}
}

func TestExecTrimmedOutput(t *testing.T) {
	sh := ShellCommand{
		command:    "echo testing-output",
		trimOutput: true,
	}

	output, err := sh.Exec()
	if err != nil {
		t.Errorf("Command failed: %v", err)
	}
	if output != "testing-output" {
		t.Errorf("Expected output 'testing-output', but received '%s'", output)
	}
}

func TestExecFailure(t *testing.T) {
	sh := ShellCommand{
		command: "command-that-does-not-exist",
	}
	_, err := sh.Exec()
	if err == nil {
		t.Error("Didn't receive expected failure from command")
	}
}

func TestOutputLineSpliting(t *testing.T) {
	sh := ShellCommand{
		command:    `printf one\ntwo\nthree`,
		silent:     true,
		trimOutput: true,
	}

	output, err := sh.Exec()
	if err != nil {
		t.Errorf("Command failed: %v", err)
	}
	if len(strings.Fields(output)) != 3 {
		t.Errorf("Expected split length of output, expected 3, but received %v", len(strings.Fields(output)))
	}
}

func TestReassembleCommandParts(t *testing.T) {
	shs := []ShellCommand{
		{
			command:    `grep -r 'some text with spaces' .`,
			silent:     false,
			trimOutput: false,
		},
		{
			command:    `grep -r "some text with spaces" .`,
			silent:     false,
			trimOutput: false,
		},
	}

	for _, sh := range shs {
		_, err := sh.Exec()
		if err != nil {
			t.Errorf("Command failed: %v", err)
		}
	}
}
