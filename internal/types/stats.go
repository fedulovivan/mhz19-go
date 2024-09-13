package types

import (
	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

type StatsGetResult struct {
	Rules                  int32      `json:"rules"`
	Devices                int32      `json:"devices"`
	Messages               int32      `json:"messages"`
	EngineMessagesReceived int32      `json:"engineMessagesReceived"`
	EngineRulesMatched     int32      `json:"engineRulesMatched"`
	ApiRequests            int32      `json:"apiRequests"`
	Uptime                 app.Uptime `json:"uptime"`
	Memory                 string     `json:"memory"`
}

func (r *StatsGetResult) WithAppStats(stats *app.AppStatCounters) {
	r.EngineMessagesReceived = int32(stats.EngineMessagesReceived.Value())
	r.EngineRulesMatched = int32(stats.EngineRulesMatched.Value())
	r.ApiRequests = int32(stats.ApiRequests.Value())
	r.Uptime = app.GetUptime()
	r.Memory = utils.GetMemUsage()
}
