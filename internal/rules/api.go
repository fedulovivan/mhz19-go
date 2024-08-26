package rules

import (
	"fmt"
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
	defer utils.TimeTrack(logTag, time.Now(), "CreateRule")
	rule := engine.Rule{}
	err := c.Read(&rule)
	if err != nil {
		return err
	}
	err = api.service.Create(rule)
	if err != nil {
		return err
	}
	return c.Write(map[string]bool{"ok": true})
}

func (api rulesApi) getOne(c *routing.Context) error {
	defer utils.TimeTrack(logTag, time.Now(), "GetRule")
	return fmt.Errorf(
		"not yet implemented /rules/%v",
		c.Param("id"),
	)
}

func (api rulesApi) getAll(c *routing.Context) error {
	defer utils.TimeTrack(logTag, time.Now(), "GetRules")
	rules, err := api.service.Get()
	if err != nil {
		return err
	}
	return c.Write(rules)
}
