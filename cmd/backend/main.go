package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/entities/devices"
	"github.com/fedulovivan/mhz19-go/internal/entities/ldm"
	"github.com/fedulovivan/mhz19-go/internal/entities/messages"
	"github.com/fedulovivan/mhz19-go/internal/entities/rules"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/rest"

	buried_devices "github.com/fedulovivan/mhz19-go/internal/providers/buried_devices"
	dnssd "github.com/fedulovivan/mhz19-go/internal/providers/dnssd"
	mqtt "github.com/fedulovivan/mhz19-go/internal/providers/mqtt"
	p_rest "github.com/fedulovivan/mhz19-go/internal/providers/rest"
	tbot "github.com/fedulovivan/mhz19-go/internal/providers/tbot"
)

var tag = logger.NewTag(logger.MAIN)

func main() {

	// bootstrap application
	app.InitConfig()
	logger.Init()

	// configure engine dependencies and start it
	rulesService := rules.ServiceSingleton(
		rules.NewRepository(
			db.DbSingleton(),
		),
	)
	devicesService := devices.NewService(
		devices.NewRepository(
			db.DbSingleton(),
		),
	)
	e := engine.NewEngine()
	e.SetLogTag(logger.NewTag(logger.ENGINE))
	restProvider := p_rest.NewProvider()
	e.SetProviders(
		mqtt.NewProvider(),
		tbot.NewProvider(),
		dnssd.NewProvider(),
		buried_devices.NewProvider(
			ldm.NewService(ldm.RepoSingleton()),
			devicesService,
		),
		restProvider,
	)
	e.SetMessagesService(
		messages.NewService(
			messages.NewRepository(
				db.DbSingleton(),
			),
		),
	)
	e.SetDevicesService(
		devicesService,
	)
	e.SetLdmService(
		ldm.NewService(
			ldm.RepoSingleton(),
		),
	)
	dbRules, err := rulesService.Get()
	if err == nil {
		if len(dbRules) > 0 {
			e.AppendRules(dbRules...)
		} else {
			slog.Warn(tag.F("No mapping rules in database"))
		}
	} else {
		slog.Error(tag.F("Failed to load rules from db"), "err", err.Error())
	}
	go func() {
		for rule := range rulesService.OnCreated() {
			e.AppendRules(rule)
		}
	}()
	go func() {
		for ruleId := range rulesService.OnDeleted() {
			e.DeleteRule(ruleId)
		}
	}()
	e.Start()

	// init rest
	rest.Init(restProvider)

	// notify we are in the development mode
	if app.Config.IsDev {
		slog.Debug(tag.F("Running in developlment mode"))
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
	slog.Debug(tag.F("App termination signal received"))

	// stop engine and all underlying providers
	e.Stop()

	// stop rest
	rest.Stop()

	slog.Info(tag.F("All done, bye-bye"))
}

// e.AppendRules(engine.GetStaticRules()...)
