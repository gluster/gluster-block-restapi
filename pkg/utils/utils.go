package utils

import (
	"bytes"
	"os/exec"
	"syscall"
)

// ExecuteCommandError represents command execution error
type ExecuteCommandError struct {
	ExitStatus int
	Errstr     string
	Err        error
}

func (e *ExecuteCommandError) Error() string {
	errstr := e.Errstr
	if errstr != "" {
		errstr = "; " + errstr
	}
	return e.Err.Error() + errstr
}

func execStderrCombined(err error, stderr *bytes.Buffer) error {
	if err == nil {
		return nil
	}

	execErr := ExecuteCommandError{
		ExitStatus: -1,
		Errstr:     stderr.String(),
		Err:        err,
	}

	exiterr, ok := err.(*exec.ExitError)
	if ok {
		status, ok := exiterr.Sys().(syscall.WaitStatus)
		if ok {
			execErr.ExitStatus = status.ExitStatus()
		}
	}

	return &execErr
}

// ExecuteCommandOutput runs the command and adds additional error information
func ExecuteCommandOutput(cmdName string, arg ...string) ([]byte, error) {
	cmd := exec.Command(cmdName, arg...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()

	if err != nil {
		return out, execStderrCombined(err, &stderr)
	}

	return out, nil
}

// ExecuteCommandRun runs the command and adds additional
// error information
func ExecuteCommandRun(cmdName string, arg ...string) error {
	cmd := exec.Command(cmdName, arg...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	return execStderrCombined(cmd.Run(), &stderr)
}
