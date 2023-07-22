package exec

import (
	"context"
	"io"
	osexec "os/exec"
	"syscall"
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

// Wraps exec.Cmd so we can capture errors.
type cmdWrapper osexec.Cmd

// assert the cmdWrapper implements the Cmd interface.
var _ Cmd = &cmdWrapper{}

func (cmd *cmdWrapper) SetDir(dir string) {
	cmd.Dir = dir
}

func (cmd *cmdWrapper) SetStdin(in io.Reader) {
	cmd.Stdin = in
}

func (cmd *cmdWrapper) SetStdout(out io.Writer) {
	cmd.Stdout = out
}

func (cmd *cmdWrapper) SetStderr(out io.Writer) {
	cmd.Stderr = out
}

func (cmd *cmdWrapper) SetEnv(env []string) {
	cmd.Env = env
}

func (cmd *cmdWrapper) StdoutPipe() (io.ReadCloser, error) {
	r, err := (*osexec.Cmd)(cmd).StdoutPipe()
	return r, handleError(err)
}

func (cmd *cmdWrapper) StderrPipe() (io.ReadCloser, error) {
	r, err := (*osexec.Cmd)(cmd).StderrPipe()
	return r, handleError(err)
}

func (cmd *cmdWrapper) Start() error {
	err := (*osexec.Cmd)(cmd).Start()
	return handleError(err)
}

func (cmd *cmdWrapper) Wait() error {
	err := (*osexec.Cmd)(cmd).Wait()
	return handleError(err)
}

// Run is part of the Cmd interface.
func (cmd *cmdWrapper) Run() error {
	err := (*osexec.Cmd)(cmd).Run()
	return handleError(err)
}

// CombinedOutput is part of the Cmd interface.
func (cmd *cmdWrapper) CombinedOutput() ([]byte, error) {
	out, err := (*osexec.Cmd)(cmd).CombinedOutput()
	return out, handleError(err)
}

func (cmd *cmdWrapper) Output() ([]byte, error) {
	out, err := (*osexec.Cmd)(cmd).Output()
	return out, handleError(err)
}

// Stop is part of the Cmd interface.
func (cmd *cmdWrapper) Stop() {
	c := (*osexec.Cmd)(cmd)

	if c.Process == nil {
		return
	}

	c.Process.Signal(syscall.SIGTERM)

	time.AfterFunc(10*time.Second, func() {
		if !c.ProcessState.Exited() {
			c.Process.Signal(syscall.SIGKILL)
		}
	})
}
