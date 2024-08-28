package stats

import (
	"time"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

var logTag = logger.MakeTag(logger.STATS)

type statsApi struct {
	service StatsService
}

func NewApi(router *routing.Router, service StatsService) {
	api := statsApi{
		service,
	}
	group := router.Group("/stats")
	group.Get("", api.get)
}

func (api statsApi) get(c *routing.Context) error {
	defer utils.TimeTrack(logTag, time.Now(), "api:get")
	data, err := api.service.Get()
	if err != nil {
		return err
	}
	return c.Write(data)
}
