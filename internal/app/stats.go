package app

import "github.com/fedulovivan/mhz19-go/pkg/utils"

type AppStatCounters struct {
	EngineMessagesReceived utils.Seq
	EngineRulesMatched     utils.Seq
	ApiRequests            utils.Seq
}

var instance *AppStatCounters

func NewStats() *AppStatCounters {
	return &AppStatCounters{
		EngineMessagesReceived: utils.NewSeq(0),
		EngineRulesMatched:     utils.NewSeq(0),
		ApiRequests:            utils.NewSeq(0),
	}
}

func StatsSingleton() *AppStatCounters {
	if instance == nil {
		instance = NewStats()
	}
	return instance
}
