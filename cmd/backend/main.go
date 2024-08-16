package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	mqtt "github.com/fedulovivan/mhz19-go/internal/service/mqtt"
	tbot "github.com/fedulovivan/mhz19-go/internal/service/tbot"
)

func init() {
	app.RecordStartTime()
}

var logTag = logger.MakeTag("MAIN")

func main() {

	// var l []any = []any{"111", "222"}
	// var v any = "11"
	// fmt.Println(slices.Contains(l, v))

	// start engine with list of connected services
	eopts := engine.NewOptions()
	eopts.SetLogTag(logger.MakeTag("ENGN"))
	eopts.SetServices(mqtt.Service, tbot.Service)
	engine.Start(eopts)

	// notify we are in the development mode
	if app.Config.IsDev {
		slog.Debug(logTag("Running in developlment mode"))
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
	slog.Debug(logTag("App termination signal received"))

	// stop engine and all underlying services
	engine.Stop()

	slog.Info(logTag("All done, bye-bye"))
}
