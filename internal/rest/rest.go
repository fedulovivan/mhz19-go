package rest

import (
	"context"

	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/counters"
	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/entities/devices"
	"github.com/fedulovivan/mhz19-go/internal/entities/dicts"
	"github.com/fedulovivan/mhz19-go/internal/entities/ldm"
	"github.com/fedulovivan/mhz19-go/internal/entities/messages"
	push_message "github.com/fedulovivan/mhz19-go/internal/entities/push-message"
	"github.com/fedulovivan/mhz19-go/internal/entities/rules"
	stats_e "github.com/fedulovivan/mhz19-go/internal/entities/stats"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/cors"
	"github.com/go-ozzo/ozzo-routing/v2/slash"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var tag = utils.NewTag(logger.REST)

var server http.Server

// borrowed with simplifications from qiangxue/go-rest-api
// see original at https://github.com/qiangxue/go-rest-api/blob/master/internal/errors/middleware.go
func errorHandler(c *routing.Context) (err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			var ok bool
			if err, ok = rerr.(error); !ok {
				err = fmt.Errorf("panic: %v", rerr)
			}
			slog.Error(tag.F("recovered from panic"))
			fmt.Println(string(debug.Stack()))
		}
		if err != nil {
			slog.Error(tag.F("errorHandler:"), "method", c.Request.Method, "path", c.Request.URL.Path, "err", err.Error())
			counters.Inc(counters.ERRORS_ALL)
			counters.Errors.WithLabelValues(logger.MOD_REST).Inc()
			res := map[string]any{
				"is_error": true,
				"error":    err.Error(),
			}
			if err = c.WriteWithStatus(res, http.StatusInternalServerError); err != nil {
				slog.Error(tag.F("failed writing error response"), "err", err)
			}
			c.Abort()
			err = nil
		}
	}()
	return c.Next()
}

// func (rg *routing.RouteGroup) Path() string {
// 	return ""
// }

func requestCounter(c *routing.Context) error {
	go counters.Inc(counters.API_REQUESTS)
	go counters.ApiRequests.WithLabelValues(
		// c.Router().RouteGroup,
		// fmt.Sprintf("%v", c.Router().RouteGroup),
		c.Request.URL.Path,
		c.Request.Method,
	).Inc()
	return c.Next()
}

func Init(shimProvider types.ChannelProvider) {

	// set custom json writer for better performance
	// see var DataWriters
	// in ozzo-routing/content/type.go
	content.DataWriters[content.JSON] = &JSONDataWriterCustom{}

	router := routing.New()
	router.Use(
		slash.Remover(http.StatusMovedPermanently),
		cors.Handler(cors.AllowAll),
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

	router.Get("/metrics", routing.HTTPHandler(promhttp.Handler()))

	apibase := router.Group(app.Config.RestApiPath)
	apibase.Use(
		errorHandler,
		content.TypeNegotiator(content.JSON),
		requestCounter,
	)

	// rules
	rules.NewApi(
		apibase,
		rules.ServiceSingleton(
			rules.NewRepository(
				db.DbSingleton(),
			),
		),
	)

	// stats
	stats_e.NewApi(
		apibase,
		stats_e.NewService(
			stats_e.NewRepository(
				db.DbSingleton(),
			),
		),
	)

	// devices
	devices.NewApi(
		apibase,
		devices.NewService(
			devices.NewRepository(
				db.DbSingleton(),
			),
		),
	)

	// messages
	messages.NewApi(
		apibase,
		messages.NewService(
			messages.NewRepository(
				db.DbSingleton(),
			),
		),
	)

	// last device message
	ldm.NewApi(
		apibase,
		ldm.NewService(
			ldm.RepoSingleton(),
		),
	)

	// push engine message received via rest
	push_message.NewApi(
		apibase,
		shimProvider,
	)

	// dictionaries
	dicts.NewApi(
		apibase,
		dicts.NewService(
			dicts.NewRepository(
				db.DbSingleton(),
			),
		),
	)

	http.Handle("/", router)
	go func() {
		addr := fmt.Sprintf(":%v", app.Config.RestApiPort)
		slog.Debug(tag.F("Server is running at " + addr))
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
		counters.Inc(counters.ERRORS_ALL)
		counters.Errors.WithLabelValues(logger.MOD_REST).Inc()
	}
}

// random id for the load tests
// idspool := []types.DeviceId{
// 	types.DeviceId("lorem ipsum"),
// 	types.DeviceId("is simply dummy "),
// 	types.DeviceId("text of the "),
// 	types.DeviceId("printing and "),
// 	types.DeviceId("typesetting industry"),
// 	types.DeviceId("has been the industrys "),
// 	types.DeviceId("standard dummy text"),
// 	types.DeviceId("ever since the 1500s"),
// }
// id := idspool[rand.Intn(len(idspool))]
