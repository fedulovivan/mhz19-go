package messages

import (
	"time"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

type messagesApi struct {
	service MessagesService
	logTag  logger.LogTagFn
}

func NewApi(router *routing.Router, service MessagesService) {
	logTag := logger.MakeTag(logger.MESSAGES)
	api := messagesApi{
		service,
		logTag,
	}
	group := router.Group("/messages")
	group.Get("", api.get)
}

func (api messagesApi) get(c *routing.Context) error {
	defer utils.TimeTrack(api.logTag, time.Now(), "api:get")
	data, err := api.service.Get()
	if err != nil {
		return err
	}
	return c.Write(data)
}
