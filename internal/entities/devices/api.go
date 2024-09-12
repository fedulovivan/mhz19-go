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
	logTag  types.LogTagFn
}

func NewApi(base *routing.RouteGroup, service types.DevicesService) {
	logTag := logger.MakeTag(logger.DEVICES)
	api := devicesApi{
		service,
		logTag,
	}
	group := base.Group("/devices")
	group.Get("", api.get)
	group.Get("/class/<deviceClass>", api.getByDeviceClass)
}

func (api devicesApi) get(c *routing.Context) error {
	defer utils.TimeTrack(api.logTag, time.Now(), "api:get")
	data, err := api.service.Get()
	if err != nil {
		return err
	}
	return c.Write(data)
}

func (api devicesApi) getByDeviceClass(c *routing.Context) error {
	defer utils.TimeTrack(api.logTag, time.Now(), "api:getByDeviceClass")
	dc, err := strconv.Atoi(c.Param("deviceClass"))
	if err != nil {
		return err
	}
	data, err := api.service.GetByDeviceClass(types.DeviceClass(dc))
	if err != nil {
		return err
	}
	return c.Write(data)
}
