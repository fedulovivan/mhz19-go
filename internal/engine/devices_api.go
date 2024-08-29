package engine

import (
	"time"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

type devicesApi struct {
	service DevicesService
	logTag  logger.LogTagFn
}

func NewDevicesApi(router *routing.Router, service DevicesService) {
	logTag := logger.MakeTag(logger.DEVICES)
	api := devicesApi{
		service,
		logTag,
	}
	group := router.Group("/devices")
	group.Get("", api.get)
}

func (api devicesApi) get(c *routing.Context) error {
	defer utils.TimeTrack(api.logTag, time.Now(), "api:get")
	data, err := api.service.Get()
	if err != nil {
		return err
	}
	return c.Write(data)
}
