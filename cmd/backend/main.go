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

	mqtt "github.com/fedulovivan/mhz19-go/internal/providers/mqtt"
	tbot "github.com/fedulovivan/mhz19-go/internal/providers/tbot"
)

var logTag = logger.MakeTag(logger.MAIN)

func main() {

	// Make a channel for results and start listening
	// mdns.DefaultParams()
	// entriesCh := make(chan *mdns.ServiceEntry, 4)
	// go func() {
	// 	for entry := range entriesCh {
	// 		fmt.Printf("Got new entry: %v\n", entry)
	// 	}
	// }()
	// err := mdns.Query(&mdns.QueryParam{
	// 	Service:     "_ewelink",
	// 	DisableIPv6: true,
	// 	Entries:     entriesCh,
	// })
	// if err != nil {
	// 	slog.Error(err.Error())
	// }
	// err :=
	// err := mdns.Lookup("_ewelink._tcp.local", entriesCh)
	// 	fmt.Println(err)
	// }
	// close(entriesCh)

	// bootstrap application
	app.ConfigInit()
	logger.Init()
	rest.Init()

	// configure and start engine
	eopts := engine.NewOptions()
	eopts.SetLogTag(logger.MakeTag(logger.ENGINE))
	eopts.SetProviders(mqtt.Provider, tbot.Provider)
	eopts.SetMessagesService(
		messages.NewService(
			messages.NewRepository(
				db.Instance(),
			),
		),
	)
	eopts.SetDevicesService(
		devices.NewService(
			devices.NewRepository(
				db.Instance(),
			),
		),
	)
	eopts.SetLdmService(
		ldm.NewService(
			ldm.RepositoryInstance(),
		),
	)
	eopts.SetRules(
		engine.GetStaticRules()...,
	)
	einst := engine.NewEngine(eopts)
	einst.Start()

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
	einst.Stop()

	// stop rest
	rest.Stop()

	slog.Info(logTag("All done, bye-bye"))
}
