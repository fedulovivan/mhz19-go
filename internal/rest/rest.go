package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/entities/devices"
	"github.com/fedulovivan/mhz19-go/internal/entities/ldm"
	"github.com/fedulovivan/mhz19-go/internal/entities/messages"
	"github.com/fedulovivan/mhz19-go/internal/entities/rules"
	"github.com/fedulovivan/mhz19-go/internal/entities/stats"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/slash"
)

var logTag = logger.MakeTag(logger.REST)

var server http.Server

func Init() {

	router := routing.New()

	router.Use(
		slash.Remover(http.StatusMovedPermanently),
		content.TypeNegotiator(content.JSON),
	)

	// rules
	rules.NewApi(
		router,
		rules.ServiceSingleton(
			rules.NewRepository(
				db.DbSingleton(),
			),
		),
	)

	// stats
	stats.NewApi(
		router,
		stats.NewService(
			stats.NewRepository(
				db.DbSingleton(),
			),
		),
	)

	// devices
	devices.NewApi(
		router,
		devices.NewService(
			devices.NewRepository(
				db.DbSingleton(),
			),
		),
	)

	// messages
	messages.NewApi(
		router,
		messages.NewService(
			messages.NewRepository(
				db.DbSingleton(),
			),
		),
	)

	// last device message
	ldm.NewApi(
		router,
		ldm.NewService(
			ldm.RepoSingleton(),
		),
	)

	http.Handle("/", router)
	go func() {
		addr := fmt.Sprintf(":%v", app.Config.RestApiPort)
		slog.Debug(logTag("server is running at " + addr))
		server = http.Server{Addr: addr}
		err := server.ListenAndServe()
		slog.Warn(logTag(err.Error()))
	}()
}

func Stop() {
	slog.Debug(logTag("Stopping rest..."))
	err := server.Shutdown(context.Background())
	if err != nil {
		slog.Error(logTag(err.Error()))
	}
}
