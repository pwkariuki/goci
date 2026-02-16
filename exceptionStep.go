package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Type for a step that handles program output
type exceptionStep struct {
	step
}

// Constructor to instantiate a new exceptionStep
func newExceptionStep(name, exe, message, proj string,
	args []string) exceptionStep {

	s := exceptionStep{}
	s.step = newStep(name, exe, message, proj, args)
	return s
}

// Execute external program and capture its output
func (s exceptionStep) execute() (string, error) {
	cmd := exec.Command(s.exe, s.args...)
	var out bytes.Buffer
	cmd.Stdout = &out // program output copied into buffer
	cmd.Dir = s.proj

	if err := cmd.Run(); err != nil {
		return "", &stepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	// If buffer contains output, at least one file doesn't match format
	if out.Len() > 0 {
		return "", &stepErr{
			step:  s.name,
			msg:   fmt.Sprintf("invalid format: %s", out.String()),
			cause: nil, // formatting errors have no underlying cause
		}
	}

	return s.message, nil
}
