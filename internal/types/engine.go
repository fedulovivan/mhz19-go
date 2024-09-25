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
	InvokeActionFunc(mm []Message, a Action, tag logger.Tag)
	MatchesCondition(mtcb MessageTupleFn, c Condition, tag logger.Tag) bool
	InvokeConditionFunc(mt MessageTuple, fn CondFn, not bool, args Args, tag logger.Tag) bool
	MatchesListSome(mtcb MessageTupleFn, cc []Condition, tag logger.Tag) bool
	MatchesListEvery(mtcb MessageTupleFn, cc []Condition, tag logger.Tag) bool
	ExecuteActions(mm []Message, r Rule, tag logger.Tag)
	HandleMessage(m Message, rules []Rule)
	SetLdmService(r LdmService)
	AppendRules(rules ...Rule)
	DeleteRule(ruleId int)
	Start()
	Stop()
}
