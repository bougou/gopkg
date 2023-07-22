package exec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Shell is a local command execution utility
type ShellCommand struct {
	// command string that you normally typed on shell
	command string
	// whether to print on the stdout and stderr of operating system.
	silent bool
	// whether to trim the output
	trim bool
	// timeout
	timeout time.Duration
}

func NewShellCommand(command string) *ShellCommand {
	return &ShellCommand{
		command: command,
	}
}

// RunShellCommand runs the specified command, and return the combined output content of
// stdout and stderr.
// If silent is true, no outputs are printed on os.stdout and os.stderr.
func RunShellCommand(command string, silent bool, trim bool) (output []byte, err error) {
	sc := NewShellCommand(command)
	sc.SetSilent(silent)
	sc.SetTrim(trim)
	return sc.Exec()
}

func RunShellCommandTimeout(command string, silent bool, trim bool, timeout time.Duration) (output []byte, err error) {
	sc := NewShellCommand(command)
	sc.SetSilent(silent)
	sc.SetTrim(trim)
	sc.SetTimeout(timeout)
	return sc.Exec()
}

func (sh *ShellCommand) SetSilent(silent bool) {
	sh.silent = silent
}

func (sh *ShellCommand) SetTrim(trim bool) {
	sh.trim = trim
}

func (sh *ShellCommand) SetTimeout(timeout time.Duration) {
	sh.timeout = timeout
}

func (sh *ShellCommand) Cmd() (cmd *exec.Cmd, err error) {
	shellstrList := ShellStr2List(sh.command)

	if len(shellstrList) == 0 {
		return nil, fmt.Errorf("no command str is specified")
	}

	if len(shellstrList) == 1 {
		cmd = exec.Command(shellstrList[0])
	}

	if len(shellstrList) > 1 {
		cmd = exec.Command(shellstrList[0], shellstrList[1:]...)
	}

	return
}

// SimpleExec just executes the shell command with setting stdin/stdout/stderr to
// standard input, standard output, and standard error file descriptors of the os.
// The behaviour of SimpleExec comes closest to that you type the command in shell and run it.
func (sh *ShellCommand) SimpleExec() error {
	cmd, err := sh.Cmd()
	if err != nil {
		return err
	}

	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if sh.timeout == 0 {
		return cmd.Run()
	}

	return RunTimeout(cmd, sh.timeout)
}

// Exec runs a local shell command, and return combined output of stdout and stderr.
func (sh *ShellCommand) Exec() ([]byte, error) {
	cmd, err := sh.Cmd()
	if err != nil {
		return []byte{}, err
	}

	var output bytes.Buffer

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

	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()

	if err := cmd.Start(); err != nil {
		return sh.mayTrim(output.Bytes()), fmt.Errorf("cmd start failed, err: %s", err)
	}

	if sh.timeout == 0 {
		err = cmd.Wait()
	} else {
		err = WaitTimeout(cmd, sh.timeout)
	}

	if err != nil {
		return sh.mayTrim(output.Bytes()), fmt.Errorf("cmd run failed, err: %s", err)
	}

	return sh.mayTrim(output.Bytes()), err
}

func (sh *ShellCommand) mayTrim(output []byte) []byte {
	if sh.trim {
		return []byte(strings.TrimSpace(string(output)))
	}
	return output
}

// Exec runs a local shell command, and return combined output of stdout and stderr.
func (sh *ShellCommand) deprecatedExec() ([]byte, error) {

	var output bytes.Buffer

	cmd, err := sh.Cmd()
	if err != nil {
		return sh.mayTrim(output.Bytes()), err
	}

	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

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

	if err := cmd.Start(); err != nil {
		return sh.mayTrim(output.Bytes()), fmt.Errorf("cmd start failed, err: %s", err)
	}

	// Copy the stdout/stderr of cmd during the cmd runs.

	var wg sync.WaitGroup
	wg.Add(2)
	var errStdout, errStderr error
	go func() {
		defer wg.Done()
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()
	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
		wg.Done()
	}()
	wg.Wait()

	// Wait cmd finished

	if err := cmd.Wait(); err != nil {
		return sh.mayTrim(output.Bytes()), fmt.Errorf("command wait failed: %s", err)
	}
	if errStdout != nil {
		return output.Bytes(), fmt.Errorf("write stdout failed, %s", errStdout)
	}
	if errStderr != nil {
		return output.Bytes(), fmt.Errorf("write stderr failed, %s", errStderr)
	}

	return sh.mayTrim(output.Bytes()), nil
}
