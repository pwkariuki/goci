package main

import "os/exec"

// Type for a step in the CI pipeline
type step struct {
	name    string   // the step name
	exe     string   // executable name of external tool to execute
	args    []string // arguments for the executable
	message string   // output message in case of success
	proj    string   // project to execute the step on
}

// Constructor to instantiate a new step
func newStep(name, exe, message, proj string, args []string) step {
	return step{
		name:    name,
		exe:     exe,
		message: message,
		args:    args,
		proj:    proj,
	}
}

// Execute external program
func (s step) execute() (string, error) {
	cmd := exec.Command(s.exe, s.args...)
	cmd.Dir = s.proj

	if err := cmd.Run(); err != nil {
		return "", &stepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	return s.message, nil
}
