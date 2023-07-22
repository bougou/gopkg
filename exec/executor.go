package exec

import (
	"context"
	osexec "os/exec"
)

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
