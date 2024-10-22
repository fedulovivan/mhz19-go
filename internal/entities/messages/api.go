package messages

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

type messagesApi struct {
	service types.MessagesService
	tag     logger.Tag
}

func NewApi(base *routing.RouteGroup, service types.MessagesService) {
	logTag := logger.NewTag(logger.MESSAGES)
	api := messagesApi{
		service,
		logTag,
	}
	group := base.Group("/messages")
	group.Get("", api.get)
	group.Get("/device/<deviceId>", api.getByDeviceId)
}

func (api messagesApi) get(c *routing.Context) error {
	defer utils.TimeTrack(api.tag.F, time.Now(), "api:get")
	data, err := api.service.Get()
	if err != nil {
		return err
	}
	return c.Write(data)
}

func (api messagesApi) getByDeviceId(c *routing.Context) error {
	defer utils.TimeTrack(api.tag.F, time.Now(), "api:getByDeviceId")
	tocsv := c.Query("tocsv") == "1"
	data, err := api.service.GetByDeviceId(c.Param("deviceId"))
	if err != nil {
		return err
	}
	if tocsv {
		for _, row := range data {
			timestamp := row.Timestamp.Format(time.RFC3339)
			tstring := ""
			if payload, ok := row.Payload.(map[string]any); ok {
				if temperature, ok := payload["temperature"].(float64); ok {
					tstring = strconv.FormatFloat(temperature, 'f', 2, 32)
				}
			}
			_, err := fmt.Fprint(
				c.Response,
				timestamp+","+tstring+"\n",
			)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return c.Write(data)
}
