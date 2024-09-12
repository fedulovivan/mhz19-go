package messages

import (
	"time"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

type messagesApi struct {
	service types.MessagesService
	logTag  types.LogTagFn
}

func NewApi(base *routing.RouteGroup, service types.MessagesService) {
	logTag := logger.MakeTag(logger.MESSAGES)
	api := messagesApi{
		service,
		logTag,
	}
	group := base.Group("/messages")
	group.Get("", api.get)
	group.Get("/device/<deviceId>", api.getByDeviceId)
}

func (api messagesApi) get(c *routing.Context) error {
	defer utils.TimeTrack(api.logTag, time.Now(), "api:get")
	data, err := api.service.Get()
	if err != nil {
		return err
	}
	return c.Write(data)
}

func (api messagesApi) getByDeviceId(c *routing.Context) error {
	defer utils.TimeTrack(api.logTag, time.Now(), "api:getByDeviceId")
	data, err := api.service.GetByDeviceId(c.Param("deviceId"))
	if err != nil {
		return err
	}
	return c.Write(data)
}
