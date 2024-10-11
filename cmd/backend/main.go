package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/counters"
	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/entities/devices"
	"github.com/fedulovivan/mhz19-go/internal/entities/ldm"
	"github.com/fedulovivan/mhz19-go/internal/entities/messages"
	"github.com/fedulovivan/mhz19-go/internal/entities/rules"
	schema_version "github.com/fedulovivan/mhz19-go/internal/entities/schema-version"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/rest"
	"github.com/fedulovivan/mhz19-go/internal/types"

	buried_devices "github.com/fedulovivan/mhz19-go/internal/providers/buried_devices"
	dnssd "github.com/fedulovivan/mhz19-go/internal/providers/dnssd"
	mqtt "github.com/fedulovivan/mhz19-go/internal/providers/mqtt"
	"github.com/fedulovivan/mhz19-go/internal/providers/shim_provider"
	tbot "github.com/fedulovivan/mhz19-go/internal/providers/tbot"

	_ "net/http/pprof"
)

var tag = logger.NewTag(logger.MAIN)

func main() {

	// bootstrap application
	app.InitConfig()
	logger.Init()

	// acquire and validate schema version
	repo := schema_version.NewRepository(
		db.DbSingleton(),
	)
	version, err := repo.GetVersion()
	if err != nil {
		panic(err.Error())
	}
	app.ValidateSchemaVersion(version)

	// configure various engine dependencies and start it
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
	shimProvider := shim_provider.NewProvider()
	e.SetProviders(
		mqtt.NewProvider(),
		tbot.NewProvider(),
		dnssd.NewProvider(),
		buried_devices.NewProvider(
			ldm.NewService(ldm.RepoSingleton()),
			devicesService,
		),
		shimProvider,
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
	go func() {
		slog.Debug(tag.F("Fetching rules..."))
		var timer *time.Timer
		timer = time.AfterFunc(app.Config.RulesFetchingLimit, func() {
			slog.Error(tag.F("Fetching rules is still running..."))
			timer.Reset(app.Config.RulesFetchingLimit)
		})
		dbRules, err := rulesService.Get()
		timer.Stop()
		if err == nil {
			if len(dbRules) > 0 {
				slog.Debug(tag.F("Rules fetching done"), "rules", len(dbRules))
				e.AppendRules(dbRules...)
			} else {
				slog.Warn(tag.F("No mapping rules in database"))
			}
		} else {
			slog.Error(tag.F("Failed to load rules from db"), "err", err.Error())
			counters.Inc(counters.ERRORS_ALL)
		}
	}()
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
	rest.Init(shimProvider)

	// notify we are in the development mode
	if app.Config.IsDev {
		slog.Debug(tag.F("Running in developlment mode"))
	}

	// publish "Application started" message
	// TODO detect bot(s) are connected, instead of using dumb timeout
	time.AfterFunc(time.Second*10, func() {
		shimProvider.Push(types.NewSystemMessage("Application started"))
	})

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
