package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/mqtt"
	"github.com/fedulovivan/mhz19-go/internal/registry"
	"github.com/fedulovivan/mhz19-go/internal/tbot"
)

func init() {
	registry.RecordStartTime()
}

var withTag = logger.MakeTag("MAIN")

func main() {

	// start engine with list of connected services
	engine.Start(mqtt.Service, tbot.Service)

	// notify we are in the development mode
	if registry.Config.IsDev {
		slog.Debug(withTag("Running in developlment mode"))
	}

	// handle shutdown
	stopped := make(chan struct{})
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range s {
			close(stopped)
		}
	}()
	<-stopped
	slog.Debug(withTag("App termination signal received"))

	// stop engine and all underlying services
	engine.Stop()

	slog.Info(withTag("All done, bye-bye"))
}
