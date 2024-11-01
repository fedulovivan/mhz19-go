package messages

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

type messagesApi struct {
	service types.MessagesService
	tag     utils.Tag
}

func NewApi(base *routing.RouteGroup, service types.MessagesService) {
	logTag := utils.NewTag(logger.MESSAGES)
	api := messagesApi{
		service,
		logTag,
	}
	group := base.Group("/messages")
	group.Get("", api.get)
	group.Get("/device/<deviceId>", api.getByDeviceId)
	group.Get("/temperature/<deviceId>", api.getTemperature)
}

func (api messagesApi) get(c *routing.Context) error {
	defer utils.TimeTrack(api.tag.F, time.Now(), "api:get")
	data, err := api.service.Get()
	if err != nil {
		return err
	}
	return c.Write(data)
}

func (api messagesApi) getTemperature(c *routing.Context) error {
	defer utils.TimeTrack(api.tag.F, time.Now(), "api:getTemperature")
	deviceId := types.DeviceId(c.Param("deviceId"))
	data, err := api.service.GetWithTemperature(deviceId)
	if err != nil {
		return err
	}
	c.Response.Header().Set("Content-Type", "text/plain")
	for _, message := range data {
		_, err := writeCsvRowToResponse(
			c.Response,
			message.Temperature,
			message.Timestamp,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeCsvRowToResponse(
	writer http.ResponseWriter,
	temperature float64,
	timestamp time.Time,
) (int, error) {
	tsString := timestamp.Format(time.RFC3339)
	tempString := strconv.FormatFloat(temperature, 'f', 2, 32)
	return fmt.Fprint(
		writer,
		tsString+","+tempString+"\n",
	)
}

func (api messagesApi) getByDeviceId(c *routing.Context) error {
	defer utils.TimeTrack(api.tag.F, time.Now(), "api:getByDeviceId")
	tocsv := c.Query("tocsv") == "1"
	deviceId := types.DeviceId(c.Param("deviceId"))
	data, err := api.service.GetByDeviceId(deviceId)
	if err != nil {
		return err
	}
	if tocsv {
		c.Response.Header().Set("Content-Type", "text/plain")
		for _, row := range data {
			if payload, ok := row.Payload.(map[string]any); ok {
				if temperature, ok := payload["temperature"].(float64); ok {
					_, err := writeCsvRowToResponse(
						c.Response,
						temperature,
						row.Timestamp,
					)
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
	return c.Write(data)
}
