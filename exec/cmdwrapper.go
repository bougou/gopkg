package exec

import (
	"io"
	osexec "os/exec"
	"syscall"
	"time"
)

// Wraps exec.Cmd so we can capture errors.
type cmdWrapper osexec.Cmd

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
