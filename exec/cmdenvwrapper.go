package exec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type CmdEnvWrapper struct {
	*exec.Cmd
	env []string

	debug bool
}

func NewCmdEnvWrapper(cmd *exec.Cmd, env ...string) *CmdEnvWrapper {
	c := &CmdEnvWrapper{cmd, env, false}

	// Must, if os.Environ is not assigned to cmd.Env, some important os native environment variables may be lost,
	// and thus influence the cmd's execute behaviour
	c.Cmd.Env = os.Environ()
	c.Cmd.Env = append(c.Cmd.Env, env...)

	return c
}

func (c *CmdEnvWrapper) Run() error {
	if c.debug {
		fmt.Println("[Run]:", c.String())
	}

	return c.Cmd.Run()
}

func (c *CmdEnvWrapper) SetDebug(debug bool) {
	c.debug = debug
}

// String returns a human-readable description of c.
// It have environment variables prefixed to the string of exec.Cmd
func (c *CmdEnvWrapper) String() string {
	b := new(strings.Builder)

	// don't loop c.Cmd.Env, it contains os native environment variables
	for _, e := range c.env {
		b.WriteString(e)
		b.WriteByte(' ')
	}

	b.WriteString(c.Cmd.String())

	return b.String()
}
