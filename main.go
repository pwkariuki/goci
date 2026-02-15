package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func run(proj string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("Project directory is required")
	}

	// Build target project using `go build` to validate program correctness.
	// Building multiple packages together ensures go build doesn't create an executable.
	args := []string{"build", ".", "errors"}
	cmd := exec.Command("go", args...)
	cmd.Dir = proj // set working directory to project directory

	if err := cmd.Run(); err != nil {
		return &stepErr{
			step:  "go build",
			msg:   "go build failed",
			cause: err,
		}
	}
	_, err := fmt.Fprintln(out, "Go Build: SUCCESS")
	return err
}

func main() {
	proj := flag.String("p", "", "Project directory")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
