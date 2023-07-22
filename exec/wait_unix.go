//go:build !windows
// +build !windows

package exec

import (
	"errors"
	"log"
	"os/exec"
	"syscall"
	"time"
)

var ErrTimeout = errors.New("command timed out")

// KillGrace is the amount of time we allow a process to shutdown before
// sending a SIGKILL.
const KillGrace = 5 * time.Second

// WaitTimeout waits for the given command to finish with a timeout.
// It assumes the command has already been started.
// If the command times out, it attempts to kill the process.
func WaitTimeout(c *exec.Cmd, timeout time.Duration) error {
	var kill *time.Timer
	term := time.AfterFunc(timeout, func() {
		err := c.Process.Signal(syscall.SIGTERM)
		if err != nil {
			log.Printf("E! [agent] Error terminating process: %s", err)
			return
		}

		kill = time.AfterFunc(KillGrace, func() {
			err := c.Process.Kill()
			if err != nil {
				log.Printf("E! [agent] Error killing process: %s", err)
				return
			}
		})
	})

	err := c.Wait()

	// Shutdown all timers (the kill timer and the term timer) before checking cmd err,
	// otherwise there is no chance to turn off these timers that have not expired.
	if kill != nil {
		kill.Stop()
	}
	termSent := !term.Stop()
	// For a timer created with AfterFunc(d, f), if t.Stop returns false, then
	// the timer has already expired and the function f has been started in its own goroutine.
	// So if termSent is true, it means the cmd does not finished before the term timer expired.

	// Now, we can check cmd err.
	// If the process exited without error treat it as success.
	// This allows a process to do a clean shutdown on signal.
	if err == nil {
		return nil
	}

	// If SIGTERM was sent then treat any process error as a timeout.
	if termSent {
		return ErrTimeout
	}

	// Otherwise there was an cmd error unrelated to termination.
	return err
}
