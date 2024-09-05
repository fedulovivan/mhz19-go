package types

type EngineAsSupplier interface {
	SetDevicesService(s DevicesService)
	DevicesService() DevicesService
	SetMessagesService(s MessagesService)
	MessagesService() MessagesService
	SetProviders(s ...ChannelProvider)
	Provider(ct ChannelType) ChannelProvider
}

type Engine interface {
	EngineAsSupplier
	InvokeActionFunc(mm []Message, a Action, r Rule, tid string)
	MatchesCondition(mt MessageTuple, c Condition, r Rule, tid string) bool
	InvokeConditionFunc(mt MessageTuple, fn CondFn, args Args, r Rule, tid string) bool
	MatchesListSome(mt MessageTuple, cc []Condition, r Rule, tid string) bool
	MatchesListEvery(mt MessageTuple, cc []Condition, r Rule, tid string) bool
	ExecuteActions(mm []Message, r Rule, tid string)
	HandleMessage(m Message, rules []Rule)
	Start()
	Stop()
	SetLogTag(f LogTagFn)
	SetLdmService(r LdmService)
	AppendRules(rules ...Rule)
}
