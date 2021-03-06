package creds

import (
	"bytes"
	"os/exec"

	"github.com/mattn/go-shellwords"
)

// ExecResult returns the results of executing a command
type ExecResult struct {
	Error    error  `json:"error"`
	ExitCode int    `json:"exitcode"`
	StdOut   string `json:"stdout"`
	StdErr   string `json:"stderr"`
}

// StringToCommand converts a string to a command slice for use in Exec
func StringToCommand(raw string) ([]string, error) {
	return shellwords.Parse(raw)
}

// ExecString runs a simple command with the provided credentials
func (c Creds) ExecString(command string) ExecResult {
	args, err := StringToCommand(command)
	if err != nil {
		return ExecResult{Error: err}
	}
	return c.Exec(args)
}

// Exec runs a command with the provided credentials
func (c Creds) Exec(command []string) ExecResult {
	logger.InfoMsgf("executing command: %v", command)

	cmd := exec.Command(command[0])
	cmd.Args = command
	cmd.Env = c.ToEnviron()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	exitCode := 0
	err := cmd.Run()
	if err != nil {
		exitCode = -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}

	return ExecResult{
		Error:    err,
		ExitCode: exitCode,
		StdOut:   stdout.String(),
		StdErr:   stderr.String(),
	}
}
