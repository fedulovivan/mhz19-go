package rules

import (
	"strconv"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

var logTag = logger.MakeTag(logger.RULES)

type rulesApi struct {
	service RulesService
}

func NewApi(router *routing.Router, service RulesService) {
	api := rulesApi{
		service,
	}
	group := router.Group("/rules")
	group.Get("", api.getAll)
	group.Get("/<id>", api.getOne)
	group.Put("", api.create)
}

func (api rulesApi) create(c *routing.Context) error {
	defer utils.TimeTrack(logTag, time.Now(), "api:create")
	rule := engine.Rule{}
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
	defer utils.TimeTrack(logTag, time.Now(), "api:getOne")
	ruleId, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return err
	}
	rule, err := api.service.GetOne(int32(ruleId))
	if err != nil {
		return err
	}
	return c.Write(rule)
}

func (api rulesApi) getAll(c *routing.Context) error {
	defer utils.TimeTrack(logTag, time.Now(), "api:getAll")
	rules, err := api.service.Get()
	if err != nil {
		return err
	}
	return c.Write(rules)
}
