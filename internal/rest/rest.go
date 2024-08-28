package rest

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/rules"
	"github.com/fedulovivan/mhz19-go/internal/stats"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/slash"
)

var logTag = logger.MakeTag(logger.REST)

func Init() {

	router := routing.New()
	database := db.Init()

	router.Use(
		slash.Remover(http.StatusMovedPermanently),
		content.TypeNegotiator(content.JSON),
	)

	rules.NewApi(
		router,
		rules.NewService(
			rules.NewRepository(
				database,
			),
		),
	)

	stats.NewApi(
		router,
		stats.NewService(
			stats.NewRepository(
				database,
			),
		),
	)

	http.Handle("/", router)
	go func() {
		addr := fmt.Sprintf(":%v", app.Config.RestApiPort)
		slog.Debug(logTag("server is running at " + addr))
		err := http.ListenAndServe(addr, nil)
		slog.Error(err.Error())
	}()
}

func Stop() {
	// TODO see "go routing.GracefulShutdown"
	slog.Debug(logTag("Stopping rest... (for now is no op)"))
}
