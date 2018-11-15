package utils

import (
	"os"
	"os/signal"
	"syscall"
)

var shutdownSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}

// SetSignalHandler returns a channel which will get closed when first shutdownSignal
// will get received. Program will terminate if second signal got caught
func SetSignalHandler() <-chan struct{} {
	stop := make(chan struct{})
	c := make(chan os.Signal, len(shutdownSignals))
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
