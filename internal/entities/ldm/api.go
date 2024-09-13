package ldm

import (
	"time"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

var tag = logger.NewTag(logger.LDM)

type ldmApi struct {
	service types.LdmService
}

func NewApi(base *routing.RouteGroup, service types.LdmService) {
	api := ldmApi{
		service,
	}
	group := base.Group("/last-device-messages")
	group.Get("", api.get)
	group.Get("/<deviceId>", api.getByDeviceId)
}

func (api ldmApi) get(c *routing.Context) error {
	defer utils.TimeTrack(tag.F, time.Now(), "api:get")
	data := api.service.GetAll()
	return c.Write(data)
}

func (api ldmApi) getByDeviceId(c *routing.Context) error {
	defer utils.TimeTrack(tag.F, time.Now(), "api:getByDeviceId")
	data := api.service.GetByDeviceId(
		types.DeviceId(c.Param("deviceId")),
	)
	return c.Write(data)
}
