package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func run(proj string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("Project directory is required")
	}

	pipeline := make([]step, 1)
	pipeline[0] = newStep(
		"go build",
		"go",
		"Go Build: SUCCESS",
		proj,
		// Build target project using `go build` to validate program correctness.
		// Building multiple packages together ensures go build doesn't create an executable.
		[]string{"build", ".", "errors"},
	)

	for _, s := range pipeline {
		msg, err := s.execute()
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(out, msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	proj := flag.String("p", "", "Project directory")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
