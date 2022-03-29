package exec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

// Shell is a local command execution utility
type ShellCommand struct {
	command string

	silent     bool
	trimOutput bool
}

func NewShellCommand(command string) *ShellCommand {
	return &ShellCommand{
		command: command,
	}
}

func (sh *ShellCommand) SetSilent(silent bool) {
	sh.silent = silent
}

func (sh *ShellCommand) SetTrimOutput(trim bool) {
	sh.trimOutput = trim
}

// Exec runs a local shell command
func (sh *ShellCommand) Exec() (string, error) {
	command := sh.command

	command = strings.Replace(command, "\\\n", "", -1)
	commandParts := strings.Split(command, " ")
	reassembledCommandParts := []string{}
	commandPartBuffer := ""
	commandPartBufferQuoteChar := ""
	// the following reconstructs some command line pieces for quoted args with spaces
	// which is a case our simple command line part split doesn't handle above
	// this is very imperfect, but will do for us here until it doesn't
	partSearch, _ := regexp.Compile(`['"]`)
	for _, commandPart := range commandParts {
		matches := partSearch.FindAllStringIndex(commandPart, -1)
		if len(matches) == 1 {
			// we're either starting or ending a buffered command part
			if commandPartBuffer == "" {
				// starting a buffer
				if strings.Contains(commandPart, `"`) {
					commandPartBufferQuoteChar = `"`
				} else {
					commandPartBufferQuoteChar = `'`
				}
				commandPartBuffer = commandPart
			} else if strings.Contains(commandPart, commandPartBufferQuoteChar) {
				// finishing a buffer
				commandPartBuffer = fmt.Sprintf(`%s %s`, commandPartBuffer, commandPart)
				reassembledCommandParts = append(reassembledCommandParts,
					strings.Replace(commandPartBuffer, commandPartBufferQuoteChar, "", -1))
				commandPartBuffer = ""
				commandPartBufferQuoteChar = ""
			}
		} else if commandPartBuffer != "" {
			// in the middle of a re-assemble
			commandPartBuffer = fmt.Sprintf("%s %s", commandPartBuffer, commandPart)
		} else {
			// just a normal re-append
			reassembledCommandParts = append(reassembledCommandParts, commandPart)
		}
	}

	cmd := exec.Command(reassembledCommandParts[0], reassembledCommandParts[1:]...)
	cmd.Env = os.Environ()

	var output bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	cmd.Stdin = os.Stdin
	var stdout io.Writer
	if sh.silent {
		stdout = io.MultiWriter(&output)
	} else {
		stdout = io.MultiWriter(os.Stdout, &output)
	}
	var stderr io.Writer
	if sh.silent {
		stderr = io.MultiWriter(&output)
	} else {
		stderr = io.MultiWriter(os.Stderr, &output)
	}
	err := cmd.Start()
	if err != nil {
		return output.String(), err
	}
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
		wg.Done()
	}()

	_, errStderr = io.Copy(stderr, stderrIn)
	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		return output.String(), fmt.Errorf("shell error: %v", output.String())
	}
	if errStdout != nil {
		return output.String(), errStdout
	}
	if errStderr != nil {
		return output.String(), errStderr
	}
	if sh.trimOutput {
		return strings.TrimSpace(string(output.String())), nil
	}
	return output.String(), nil
}
