package types

import "github.com/fedulovivan/mhz19-go/internal/logger"

type EngineAsSupplier interface {
	SetDevicesService(s DevicesService)
	DevicesService() DevicesService
	SetMessagesService(s MessagesService)
	MessagesService() MessagesService
	SetProviders(s ...ChannelProvider)
	FindProvider(ct ChannelType) ChannelProvider
}

type Engine interface {
	EngineAsSupplier
	InvokeActionFunc(compound MessageCompound, a Action, tag logger.Tag)
	MatchesCondition(mtcb GetCompoundForOtherDeviceId, c Condition, tag logger.Tag) bool
	InvokeConditionFunc(mt MessageCompound, c Condition, tag logger.Tag) bool
	MatchesListSome(mtcb GetCompoundForOtherDeviceId, cc []Condition, tag logger.Tag) bool
	MatchesListEvery(mtcb GetCompoundForOtherDeviceId, cc []Condition, tag logger.Tag) bool
	ExecuteActions(compound MessageCompound, r Rule, tag logger.Tag)
	HandleMessage(m Message, rules []Rule)
	SetLdmService(r LdmService)
	AppendRules(rules ...Rule)
	DeleteRule(ruleId int)
	Start()
	Stop()
}
