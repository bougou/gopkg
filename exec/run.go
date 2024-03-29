package exec

import (
	"bytes"
	"os/exec"
	"time"
)

// CombinedOutputTimeout runs the given command with the given timeout and
// returns the combined output of stdout and stderr.
// If the command times out, it attempts to kill the process.
//
// 混合输出 stdout and stderr
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

// SeparatedOutputTimeout runs the given command with the given timeout and
// returns the output of stdout and stderr separately.
// If the command times out, it attempts to kill the process.
func SeparatedOutputTimeout(c *exec.Cmd, timeout time.Duration) (stdout []byte, stderr []byte, err error) {
	var o, e bytes.Buffer
	c.Stdout = &o
	c.Stderr = &e
	if err := c.Start(); err != nil {
		return nil, nil, err
	}
	err = WaitTimeout(c, timeout)
	return o.Bytes(), e.Bytes(), err
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
		// -n means non-interactive mode
		cmd = exec.Command("sudo", append([]string{"-n", command}, args...)...)
	}
	return CombinedOutputTimeout(cmd, time.Duration(timeout))
}
