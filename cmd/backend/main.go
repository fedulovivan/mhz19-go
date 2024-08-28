package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/rest"

	mqtt "github.com/fedulovivan/mhz19-go/internal/providers/mqtt"
	tbot "github.com/fedulovivan/mhz19-go/internal/providers/tbot"
)

var logTag = logger.MakeTag(logger.MAIN)

func main() {

	// bootstrap
	app.Init()
	logger.Init()
	rest.Init()

	// start engine
	engineOptions := engine.NewOptions()
	engineOptions.SetLogTag(logger.MakeTag(logger.ENGINE))
	engineOptions.SetProviders(mqtt.Provider, tbot.Provider)
	engineOptions.SetMessagesService(
		engine.NewService(
			engine.NewRepository(
				db.Init(),
			),
		),
	)
	engineOptions.SetRules(engine.GetTestStaticRules())
	engineInstance := engine.NewEngine(engineOptions)
	engineInstance.Start()

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

	// stop engine and all underlying providers
	engineInstance.Stop()

	// stop rest
	rest.Stop()

	slog.Info(logTag("All done, bye-bye"))
}
