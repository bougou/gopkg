package exec

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	osexec "os/exec"
	"time"
)

// Interface is an interface that presents a subset of the os/exec API. Use this
// when you want to inject fakeable/mockable exec behavior.
type Interface interface {
	// Command returns a Cmd instance which can be used to run a single command.
	// This follows the pattern of package os/exec.
	Command(cmd string, args ...string) Cmd

	// CommandContext returns a Cmd instance which can be used to run a single command.
	//
	// The provided context is used to kill the process if the context becomes done
	// before the command completes on its own. For example, a timeout can be set in
	// the context.
	CommandContext(ctx context.Context, cmd string, args ...string) Cmd

	// LookPath wraps os/exec.LookPath
	LookPath(file string) (string, error)
}

// Cmd is an interface that presents an API that is very similar to Cmd from os/exec.
// As more functionality is needed, this can grow. Since Cmd is a struct, we will have
// to replace fields with get/set method pairs.
type Cmd interface {
	// Run runs the command to the completion.
	Run() error
	// CombinedOutput runs the command and returns its combined standard output
	// and standard error. This follows the pattern of package os/exec.
	CombinedOutput() ([]byte, error)
	// Output runs the command and returns standard output, but not standard err
	Output() ([]byte, error)
	SetDir(dir string)
	SetStdin(in io.Reader)
	SetStdout(out io.Writer)
	SetStderr(out io.Writer)
	SetEnv(env []string)

	// StdoutPipe and StderrPipe for getting the process' Stdout and Stderr as
	// Readers
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)

	// Start and Wait are for running a process non-blocking
	Start() error
	Wait() error

	// Stops the command by sending SIGTERM. It is not guaranteed the
	// process will stop before this function returns. If the process is not
	// responding, an internal timer function will send a SIGKILL to force
	// terminate after 10 seconds.
	Stop()
}

// Implements Interface in terms of really exec()ing.
type executor struct{}

// New returns a new Interface which will os/exec to run commands.
func New() Interface {
	return &executor{}
}

// Command is part of the Interface interface.
func (executor *executor) Command(cmd string, args ...string) Cmd {
	return (*cmdWrapper)(osexec.Command(cmd, args...))
}

// CommandContext is part of the Interface interface.
func (executor *executor) CommandContext(ctx context.Context, cmd string, args ...string) Cmd {
	return (*cmdWrapper)(osexec.CommandContext(ctx, cmd, args...))
}

// LookPath is part of the Interface interface
func (executor *executor) LookPath(file string) (string, error) {
	return osexec.LookPath(file)
}

// CombinedOutputTimeout runs the given command with the given timeout and
// returns the combined output of stdout and stderr.
// If the command times out, it attempts to kill the process.
func CombinedOutputTimeout(c *exec.Cmd, timeout time.Duration) ([]byte, error) {
	var b bytes.Buffer
	c.Stdout = &b
	c.Stderr = &b
	if err := c.Start(); err != nil {
		return nil, err
	}
	err := WaitTimeout(c, timeout)
	return b.Bytes(), err
}

// StdOutputTimeout runs the given command with the given timeout and
// returns the output of stdout.
// If the command times out, it attempts to kill the process.
func StdOutputTimeout(c *exec.Cmd, timeout time.Duration) ([]byte, error) {
	var b bytes.Buffer
	c.Stdout = &b
	c.Stderr = nil
	if err := c.Start(); err != nil {
		return nil, err
	}
	err := WaitTimeout(c, timeout)
	return b.Bytes(), err
}

// RunTimeout runs the given command with the given timeout.
// If the command times out, it attempts to kill the process.
func RunTimeout(c *exec.Cmd, timeout time.Duration) error {
	if err := c.Start(); err != nil {
		return err
	}
	return WaitTimeout(c, timeout)
}

func RunCmd(timeout time.Duration, sudo bool, command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	if sudo {
		cmd = exec.Command("sudo", append([]string{"-n", command}, args...)...)
	}
	return CombinedOutputTimeout(cmd, time.Duration(timeout))
}
