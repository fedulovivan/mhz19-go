package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

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

var tag = logger.NewTag(logger.REST)

var server http.Server

// borrowed with simplifications from qiangxue/go-rest-api
// see original at https://github.com/qiangxue/go-rest-api/blob/master/internal/errors/middleware.go
func errorHandler(c *routing.Context) (err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			var ok bool
			if err, ok = rerr.(error); !ok {
				err = fmt.Errorf("%v", rerr)
			}
			slog.Error(tag.F("recovered from panic"))
			fmt.Println(string(debug.Stack()))
		}
		if err != nil {
			slog.Error(tag.F("errorHandler:"), "path", c.Request.URL.Path, "err", err.Error())
			res := map[string]any{
				"is_error": true,
				"error":    err.Error(),
			}
			if err = c.Write(res); err != nil {
				slog.Error(tag.F("failed writing error response"), "err", err)
			}
			c.Abort()
			err = nil
		}
	}()
	return c.Next()
}

func requestCounter(c *routing.Context) error {
	go app.StatsSingleton().ApiRequests.Inc()
	return c.Next()
}

func Init() {

	router := routing.New()
	router.Use(
		slash.Remover(http.StatusMovedPermanently),
	)

	router.Get("/", func(ctx *routing.Context) error {
		return ctx.Write(fmt.Sprintf(
			"server root, rest api is available at %v",
			app.Config.RestApiPath,
		))
	})
	router.Get(app.Config.RestApiPath, func(ctx *routing.Context) error {
		return ctx.Write("rest api root")
	})

	basegroup := router.Group(app.Config.RestApiPath)
	basegroup.Use(
		errorHandler,
		content.TypeNegotiator(content.JSON),
		requestCounter,
	)

	// rules
	rules.NewApi(
		basegroup,
		rules.ServiceSingleton(
			rules.NewRepository(
				db.DbSingleton(),
			),
		),
	)

	// stats
	stats.NewApi(
		basegroup,
		stats.NewService(
			stats.NewRepository(
				db.DbSingleton(),
			),
		),
	)

	// devices
	devices.NewApi(
		basegroup,
		devices.NewService(
			devices.NewRepository(
				db.DbSingleton(),
			),
		),
	)

	// messages
	messages.NewApi(
		basegroup,
		messages.NewService(
			messages.NewRepository(
				db.DbSingleton(),
			),
		),
	)

	// last device message
	ldm.NewApi(
		basegroup,
		ldm.NewService(
			ldm.RepoSingleton(),
		),
	)

	http.Handle("/", router)
	go func() {
		addr := fmt.Sprintf(":%v", app.Config.RestApiPort)
		slog.Debug(tag.F("server is running at " + addr))
		server = http.Server{Addr: addr}
		err := server.ListenAndServe()
		slog.Debug(tag.F(err.Error()))
	}()
}

func Stop() {
	slog.Debug(tag.F("Stopping rest..."))
	err := server.Shutdown(context.Background())
	if err != nil {
		slog.Error(tag.F(err.Error()))
	}
}
