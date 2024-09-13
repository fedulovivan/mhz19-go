package stats

import (
	"time"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

var tag = logger.NewTag(logger.STATS)

type statsApi struct {
	service types.StatsService
}

func NewApi(base *routing.RouteGroup, service types.StatsService) {
	api := statsApi{
		service,
	}
	group := base.Group("/stats")
	group.Get("", api.get)
}

func (api statsApi) get(c *routing.Context) error {
	defer utils.TimeTrack(tag.F, time.Now(), "api:get")
	data, err := api.service.Get()
	if err != nil {
		return err
	}
	return c.Write(data)
}
