package server

import (
	"os"
	"os/signal"
)

var onlyOneSignaleHandler = make(chan struct{})

var shutdownHandler chan os.Signal

// SetupSignalHandler 设置
func SetupSignalHandler() <-chan struct{} {
	// 同一个channle仅允许被close一次，再次close会panic.
	// 利用这个我确保onleyOneSignaleHandler只被close一次.
	close(onlyOneSignaleHandler)

	shutdownHandler = make(chan os.Signal, 2)
	stop := make(chan struct{})

	signal.Notify(shutdownHandler, shutdownSignals...)

	go func() {
		<-shutdownHandler
		close(stop)
		<-shutdownHandler
		os.Exit(1)
	}()

	return stop
}
