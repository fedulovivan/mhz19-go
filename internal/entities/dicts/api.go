package dicts

import (
	"time"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

var tag = logger.NewTag(logger.DICTS)

type api struct {
	service types.DictsService
}

func NewApi(base *routing.RouteGroup, service types.DictsService) {
	api := api{
		service,
	}
	group := base.Group("/dicts")
	group.Get("", api.all)
	group.Get("/<type>", api.get)
}

func (api api) all(c *routing.Context) (err error) {
	defer utils.TimeTrack(tag.F, time.Now(), "api:all")
	out, err := api.service.All()
	if err != nil {
		return
	}
	return c.Write(out)
}

func (api api) get(c *routing.Context) (err error) {
	defer utils.TimeTrack(tag.F, time.Now(), "api:get")
	out, err := api.service.Get(
		types.DictType(c.Param("type")),
	)
	if err != nil {
		return
	}
	return c.Write(out)
}
