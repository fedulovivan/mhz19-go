package devices

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

type devicesApi struct {
	service types.DevicesService
	tag     utils.Tag
}

func NewApi(base *routing.RouteGroup, service types.DevicesService) {
	logTag := utils.NewTag(logger.DEVICES)
	api := devicesApi{
		service,
		logTag,
	}
	group := base.Group("/devices")
	group.Get("", api.get)
	group.Put("", api.create)
	group.Get("/class/<deviceClass>", api.getByDeviceClass)
	group.Get("/<nativeId>", api.getByDeviceId)
	group.Post("/<nativeId>", api.updateName) // Deprecated: for backward compatibility with frontend
	group.Post("/<nativeId>/name", api.updateName)
	group.Post("/<nativeId>/buried-timeout", api.updateBuriedTimeout)
	group.Delete("/<id>", api.delete)
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

func (api devicesApi) create(c *routing.Context) error {
	defer utils.TimeTrack(api.tag.F, time.Now(), "api:create")
	origin := "api"
	device := types.Device{
		Origin: &origin,
	}
	err := c.Read(&device)
	if err != nil {
		return err
	}
	id, err := api.service.UpsertAll([]types.Device{device})
	if err != nil {
		return err
	}
	return c.Write(map[string]any{"ok": true, "deviceId": id})
}

func (api devicesApi) updateBuriedTimeout(c *routing.Context) error {
	defer utils.TimeTrack(api.tag.F, time.Now(), "api:updateBuriedTimeout")
	device := types.Device{
		DeviceId: types.DeviceId(c.Param("nativeId")),
	}
	payload := new(struct {
		Action string `json:"action"`
		Value  int64  `json:"value"`
	})
	err := c.Read(payload)
	if err != nil {
		return err
	}
	switch payload.Action {
	case "set":
		device.BuriedTimeout = &types.BuriedTimeout{Duration: time.Second * time.Duration(payload.Value)}
	case "off":
		device.BuriedTimeout = &types.BuriedTimeout{Duration: time.Second * 0}
	case "reset":
		// default nil will be used
	default:
		return fmt.Errorf("Unknown action=[%s]", payload.Action)
	}
	err = api.service.UpdateBuriedTimeout(device)
	if err != nil {
		return err
	}
	return c.Write(map[string]any{"ok": true})
}

func (api devicesApi) updateName(c *routing.Context) error {
	defer utils.TimeTrack(api.tag.F, time.Now(), "api:updateName")
	device := types.Device{
		DeviceId: types.DeviceId(c.Param("nativeId")),
	}
	err := c.Read(&device)
	if err != nil {
		return err
	}
	err = api.service.UpdateName(device)
	if err != nil {
		return err
	}
	return c.Write(map[string]any{"ok": true})
}

func (api devicesApi) getByDeviceId(c *routing.Context) (err error) {
	defer utils.TimeTrack(api.tag.F, time.Now(), "api:getByDeviceId")
	data, err := api.service.GetOne(types.DeviceId(c.Param("nativeId")))
	if err != nil {
		return
	}
	return c.Write(data)
}

func (api devicesApi) delete(c *routing.Context) error {
	defer utils.TimeTrack(api.tag.F, time.Now(), "api:delete")
	ruleId, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return err
	}
	err = api.service.Delete(ruleId)
	if err != nil {
		return err
	}
	return c.Write(map[string]any{"ok": true})
}
