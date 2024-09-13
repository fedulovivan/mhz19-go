package devices

import (
	"strconv"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

type devicesApi struct {
	service types.DevicesService
	tag     logger.Tag
}

func NewApi(base *routing.RouteGroup, service types.DevicesService) {
	logTag := logger.NewTag(logger.DEVICES)
	api := devicesApi{
		service,
		logTag,
	}
	group := base.Group("/devices")
	group.Get("", api.get)
	group.Get("/class/<deviceClass>", api.getByDeviceClass)
	group.Get("/<deviceId>", api.getByDeviceId)
}

func (api devicesApi) get(c *routing.Context) (err error) {
	defer utils.TimeTrack(api.tag.F, time.Now(), "api:get")
	data, err := api.service.Get()
	if err != nil {
		return
	}
	return c.Write(data)
}

func (api devicesApi) getByDeviceClass(c *routing.Context) (err error) {
	defer utils.TimeTrack(api.tag.F, time.Now(), "api:getByDeviceClass")
	dc, err := strconv.Atoi(c.Param("deviceClass"))
	if err != nil {
		return
	}
	data, err := api.service.GetByDeviceClass(types.DeviceClass(dc))
	if err != nil {
		return
	}
	return c.Write(data)
}

func (api devicesApi) getByDeviceId(c *routing.Context) (err error) {
	defer utils.TimeTrack(api.tag.F, time.Now(), "api:getByDeviceId")
	data, err := api.service.GetOne(types.DeviceId(c.Param("deviceId")))
	if err != nil {
		return
	}
	return c.Write(data)
}
