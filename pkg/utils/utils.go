package utils

import (
	"bytes"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
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

	log.Debug("executing command: ", cmdName, " ", strings.Join(arg, " "))
	out, err := cmd.Output()

	if err != nil {
		log.WithError(err).Error("got error in executing command")
		return out, execStderrCombined(err, &stderr)
	}

	log.Debug("command executed successfully")
	return out, nil
}

// ExecuteCommandRun runs the command and adds additional
// error information
func ExecuteCommandRun(cmdName string, arg ...string) error {
	cmd := exec.Command(cmdName, arg...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	log.Debug("executing command: ", cmdName, " ", strings.Join(arg, " "))
	if err := execStderrCombined(cmd.Run(), &stderr); err != nil {
		log.WithError(err).Error("got error in executing command")
		return err
	}
	log.Debug("command executed successfully")
	return nil
}
