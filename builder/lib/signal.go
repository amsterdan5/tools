//go:build !windows
// +build !windows

package tools_lib

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/amsterdan/tools/logger"
)

func regExitSignals() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		logger.Infof("got signal %v, exit", sig)
		os.Exit(11)
	}()
}
