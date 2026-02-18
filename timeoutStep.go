package main

import (
	"context"
	"errors"
	"os/exec"
	"time"
)

// Type for an external command step that times out
type timeoutStep struct {
	step
	timeout time.Duration // timeout value of this step
}

func newTimeoutStep(name, exe, message, proj string,
	args []string, timeout time.Duration) timeoutStep {

	s := timeoutStep{
		step:    newStep(name, exe, message, proj, args),
		timeout: timeout,
	}

	// If no timeout provided to constructor, default to 30 seconds
	if s.timeout == 0 {
		s.timeout = 30 * time.Second
	}

	return s
}

// Execute external program with timeout context
func (s timeoutStep) execute() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, s.exe, s.args...)
	cmd.Dir = s.proj

	if err := cmd.Run(); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return "", &stepErr{
				step:  s.name,
				msg:   "failed time out",
				cause: context.DeadlineExceeded,
			}
		}

		return "", &stepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	return s.message, nil
}
