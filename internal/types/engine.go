package types

type Engine interface {
	GetOptions() EngineOptions
	FindProvider(ct ChannelType) ChannelProvider
	MatchesCondition(mt MessageTuple, c Condition, r Rule, tid string) bool
	InvokeConditionFunc(mt MessageTuple, fn CondFn, args Args, r Rule, tid string) bool
	MatchesListSome(mt MessageTuple, cc []Condition, r Rule, tid string) bool
	MatchesListEvery(mt MessageTuple, cc []Condition, r Rule, tid string) bool
	ExecuteActions(mm []Message, r Rule, tid string)
	HandleMessage(m Message, rules []Rule)
	Start()
	Stop()
}

type EngineOptions interface {
	SetLogTag(f LogTagFn)
	SetProviders(s ...ChannelProvider)
	SetMessagesService(s MessagesService)
	SetDevicesService(s DevicesService)
	SetLdmService(r LdmService)
	SetRules(rules ...Rule)
	LogTag() LogTagFn
	Providers() []ChannelProvider
	MessagesService() MessagesService
	DevicesService() DevicesService
	LdmService() LdmService
	Rules() []Rule
}
