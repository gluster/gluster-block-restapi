package main

import (
	"fmt"
	"io"
	stdlog "log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

// LogWriter represents writer object
var LogWriter io.WriteCloser

func openLogFile(filepath string) (io.WriteCloser, error) {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func setLogOutput(w io.Writer) {
	log.SetOutput(w)
	stdlog.SetOutput(log.StandardLogger().Writer())
}

func initLogger(logdir string, logfile string, loglevel string) error {
	// Close the previously opened log file
	if LogWriter != nil {
		err := LogWriter.Close()
		if err != nil {
			return err
		}
		LogWriter = nil
	}

	log.AddHook(&SourceHook{})
	level, err := log.ParseLevel(strings.ToLower(loglevel))
	if err != nil {
		setLogOutput(os.Stderr)
		log.WithError(err).Debug("Failed to parse log level")
		return err
	}
	log.SetLevel(level)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05.000000"})

	if strings.ToLower(logfile) == "stderr" || logfile == "-" {
		setLogOutput(os.Stderr)
	} else if strings.ToLower(logfile) == "stdout" {
		setLogOutput(os.Stdout)
	} else {
		logFilePath := path.Join(logdir, logfile)
		logFile, err := openLogFile(logFilePath)
		if err != nil {
			setLogOutput(os.Stderr)
			log.WithError(err).Debugf("Failed to open log file %s", logFilePath)
			return err
		}
		setLogOutput(logFile)
		LogWriter = logFile
	}
	return nil
}

const (
	// sourceField is the field name used for logging source location.
	sourceField = "source"
	repo        = "github.com/gluster/gluster-block-restapi"
)

// ref: github.com/gluster/glusterd2/pkg/logging/hook.go

// SourceHook provides information about source location of logging
// It implements logrus.Hook interface
type SourceHook struct {
}

// Levels returns all logrus levels. The hook is fired only for those log
// levels returned by this function.
func (s *SourceHook) Levels() []log.Level {
	return log.AllLevels
}

// Fire adds file name, function name and line number to the log entry.
func (s *SourceHook) Fire(entry *log.Entry) error {
	pcs := make([]uintptr, 3)
	n := runtime.Callers(6, pcs)
	if n == 0 {
		return nil
	}

	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		if strings.Contains(frame.File, repo) && !strings.Contains(frame.File, "vendor") {
			entry.Data[sourceField] = fmt.Sprintf("[%s:%d:%s]", filepath.Base(frame.File), frame.Line, filepath.Base(frame.Function))
			break
		}
		if !more {
			break
		}
	}

	return nil
}
