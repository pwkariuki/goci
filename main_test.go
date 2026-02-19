package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		proj     string
		out      string
		expErr   error
		setupGit bool
		mockCmd  func(ctx context.Context, name string, arg ...string) *exec.Cmd
	}{
		{name: "success", proj: "./testdata/tool/",
			out: "Go Build: SUCCESS\n" +
				"Go Test: SUCCESS\n" +
				"Gofmt: SUCCESS\n" +
				"Git Push: SUCCESS\n",
			expErr:   nil,
			setupGit: true,
			mockCmd:  nil},
		{name: "successMock", proj: "./testdata/tool/",
			out: "Go Build: SUCCESS\n" +
				"Go Test: SUCCESS\n" +
				"Gofmt: SUCCESS\n" +
				"Git Push: SUCCESS\n",
			expErr:   nil,
			setupGit: false,
			mockCmd:  mockCmdContext},
		{name: "fail", proj: "./testdata/toolErr",
			out:      "",
			expErr:   &stepErr{step: "go build"},
			setupGit: false,
			mockCmd:  nil},
		{name: "failFormat", proj: "./testdata/toolFmtErr",
			out:      "",
			expErr:   &stepErr{step: "go fmt"},
			setupGit: false,
			mockCmd:  nil},
		{name: "failTimeout", proj: "./testdata/tool",
			out:      "",
			expErr:   context.DeadlineExceeded,
			setupGit: false,
			mockCmd:  mockCmdTimeout},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupGit {
				_, err := exec.LookPath("git")
				if err != nil {
					t.Skip("Git not installed. Skipping test.")
				}

				cleanup := setupGit(t, tc.proj)
				defer cleanup()
			}

			if tc.mockCmd != nil {
				command = tc.mockCmd
			}

			var out bytes.Buffer
			err := run(tc.proj, &out)

			// test expects an error
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Expected error: %q. Got 'nil' instead", tc.expErr)
					return
				}

				if !errors.Is(err, tc.expErr) {
					t.Errorf("Expected error: %q. Got %q.", tc.expErr, err)
				}
				return
			}

			// test does not expect an error
			if err != nil {
				t.Errorf("Unexpected error: %q", err)
			}

			if out.String() != tc.out {
				t.Errorf("Expected output: %q. Got %q.", tc.out, out.String())
			}
		})
	}
}

// Test helper function that uses the `git` command to create a Bare Git
// repository that works like an external Git service like GitHub. A bare
// git repository is a repository that contains only the git data but no
// working directory, so it cannot be used to make local modifications to
// the code. This makes it well suited to serve as a remote repository.
// Returns cleanup function to delete the temporary directory and the local
// .git subdirectory in the target project directory.
func setupGit(t *testing.T, proj string) func() {
	t.Helper()

	// Verify user has git installed
	gitExec, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}

	// Create a temporary directory for the simulated remote git repository
	tempDir, err := os.MkdirTemp("", "gocitest")
	if err != nil {
		t.Fatal(err)
	}

	// Full path of the target project directory
	projPath, err := filepath.Abs(proj)
	if err != nil {
		t.Fatal(err)
	}

	// URI of the simulated remote git repository
	remoteURI := fmt.Sprintf("file://%s", tempDir)

	gitCmdList := []struct {
		args []string
		dir  string
		env  []string
	}{
		{[]string{"init", "--bare"}, tempDir, nil},
		{[]string{"init"}, projPath, nil},
		{[]string{"remote", "add", "origin", remoteURI}, projPath, nil},
		{[]string{"add", "."}, projPath, nil},
		{[]string{"commit", "-m", "test"}, projPath,
			[]string{
				"GIT_COMMITTER_NAME=test",
				"GIT_COMMITTER_EMAIL=test@example.com",
				"GIT_AUTHOR_NAME=test",
				"GIT_AUTHOR_EMAIL=test@example.com",
			}},
	}

	for _, g := range gitCmdList {
		gitCmd := exec.Command(gitExec, g.args...)
		gitCmd.Dir = g.dir

		if g.env != nil {
			// Inject env variables into the external command environment
			gitCmd.Env = append(os.Environ(), g.env...)
		}

		if err := gitCmd.Run(); err != nil {
			t.Fatal(err)
		}
	}

	return func() {
		os.RemoveAll(tempDir)
		os.RemoveAll(filepath.Join(projPath, ".git"))
	}
}

// This function creates a mock command for testing code that runs external commands.
// Instead of running the real executable, it runs the test binary itself with the flag -test.run=TestHelperProcess.
// The TestHelperProcess function simulates the command's behavior. To prevent Go from running this function during normal tests,
// we use the environment variable GO_WANT_HELPER_PROCESS. Only when this variable is set does TestHelperProcess execute its mock logic.
// This allows us to safely simulate external commands (like git) in tests without actually running them.
func mockCmdContext(ctx context.Context, exe string, args ...string) *exec.Cmd {
	// Arguments to be passed to the command
	cs := []string{"-test.run=TestHelperProcess"}
	cs = append(cs, exe)     // executable to run
	cs = append(cs, args...) // command arguments

	cmd := exec.CommandContext(ctx, os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// Mock function to simulate a command that timesout
func mockCmdTimeout(ctx context.Context, exe string, args ...string) *exec.Cmd {
	cmd := mockCmdContext(ctx, exe, args...)
	cmd.Env = append(cmd.Env, "GO_HELPER_TIMEOUT=1")
	return cmd
}

// Simulate the command
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return // prevent execution if not called from the mock command
	}

	if os.Getenv("GO_HELPER_TIMEOUT") == "1" {
		time.Sleep(15 * time.Second) // simulate long running process
	}

	if os.Args[2] == "git" {
		fmt.Fprintln(os.Stdout, "Everything up-to-date")
		os.Exit(0)
	}

	os.Exit(1)
}
