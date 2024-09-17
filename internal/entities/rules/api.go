package rules

import (
	"strconv"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

var tag = logger.NewTag(logger.RULES)

type rulesApi struct {
	service types.RulesService
}

func NewApi(base *routing.RouteGroup, service types.RulesService) {
	api := rulesApi{
		service,
	}
	group := base.Group("/rules")
	group.Get("", api.getAll)
	group.Get("/<id>", api.getOne)
	group.Delete("/<id>", api.delete)
	group.Put("", api.create)
}

func (api rulesApi) create(c *routing.Context) error {
	defer utils.TimeTrack(tag.F, time.Now(), "api:create")
	rule := types.Rule{}
	err := c.Read(&rule)
	if err != nil {
		return err
	}
	ruleId, err := api.service.Create(rule)
	if err != nil {
		return err
	}
	return c.Write(map[string]any{"ok": true, "ruleId": ruleId})
}

func (api rulesApi) getOne(c *routing.Context) error {
	defer utils.TimeTrack(tag.F, time.Now(), "api:getOne")
	ruleId, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return err
	}
	rule, err := api.service.GetOne(int(ruleId))
	if err != nil {
		return err
	}
	return c.Write(rule)
}

func (api rulesApi) delete(c *routing.Context) error {
	defer utils.TimeTrack(tag.F, time.Now(), "api:delete")
	ruleId, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return err
	}
	err = api.service.Delete(int(ruleId))
	if err != nil {
		return err
	}
	return c.Write(map[string]any{"ok": true})
}

func (api rulesApi) getAll(c *routing.Context) error {
	defer utils.TimeTrack(tag.F, time.Now(), "api:getAll")
	rules, err := api.service.Get()
	if err != nil {
		return err
	}
	return c.Write(rules)
}
