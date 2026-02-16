# goci - A Simple Go Continuous Integration Tool

A lightweight, command line implementation of a Continuous Integration (CI) tool for Go programs. This project demonstrates how to control and orchestrate multiple processes using Go.

**⚠️ Note:** This is a learning project and is **not production-ready**.

## Overview

`goci` automates the core steps of a CI pipeline for Go projects. It verifies that your code is valid, tested, properly formatted, and ready to be pushed to your remote repository.

## CI Pipeline Steps

The tool executes the following steps in sequence:

1. **Build** - Runs `go build` to verify the program structure is valid
2. **Test** - Runs `go test` to ensure the program behaves as intended
3. **Format** - Runs `gofmt` to ensure code conforms to Go formatting standards
4. **Push** - Runs `git push` to push code to the remote repository

If any step fails, the pipeline stops and reports the error.

## Requirements

- Go 1.16 or higher
- Git (for the push step)
- A Go project with tests

## How It Works

The tool executes each pipeline step using Go's `os/exec` package, capturing output and handling errors appropriately. Each step must succeed before moving to the next one.

## Learning Value

This project is a great starting point for understanding:

- How CI/CD pipelines work
- Process management in Go
- Building command-line automation tools
- Error handling and process orchestration
