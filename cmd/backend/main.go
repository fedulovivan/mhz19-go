package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/devices"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	ldm "github.com/fedulovivan/mhz19-go/internal/last_device_message"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/messages"
	"github.com/fedulovivan/mhz19-go/internal/rest"
	"github.com/fedulovivan/mhz19-go/internal/rules"

	dnssd "github.com/fedulovivan/mhz19-go/internal/providers/dnssd"
	mqtt "github.com/fedulovivan/mhz19-go/internal/providers/mqtt"
	tbot "github.com/fedulovivan/mhz19-go/internal/providers/tbot"
)

var logTag = logger.MakeTag(logger.MAIN)

func main() {

	// bootstrap application
	app.InitConfig()
	logger.Init()
	rest.Init()

	// configure and start engine
	rulesService := rules.ServiceSingleton(
		rules.NewRepository(
			db.DbSingleton(),
		),
	)
	dbRules, _ := rulesService.Get()
	e := engine.NewEngine()
	e.SetLogTag(logger.MakeTag(logger.ENGINE))
	e.SetProviders(mqtt.Provider, tbot.Provider, dnssd.Provider)
	e.SetMessagesService(
		messages.NewService(
			messages.NewRepository(
				db.DbSingleton(),
			),
		),
	)
	e.SetDevicesService(
		devices.NewService(
			devices.NewRepository(
				db.DbSingleton(),
			),
		),
	)
	e.SetLdmService(
		ldm.NewService(
			ldm.RepoSingleton(),
		),
	)
	e.AppendRules(engine.GetStaticRules()...)
	e.AppendRules(dbRules...)
	go func() {
		for rule := range rulesService.OnCreated() {
			e.AppendRules(rule)
		}
	}()
	e.Start()

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
	e.Stop()

	// stop rest
	rest.Stop()

	slog.Info(logTag("All done, bye-bye"))
}
