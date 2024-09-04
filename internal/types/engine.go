package types

type Engine interface {
	FindProvider(ct ChannelType) ChannelProvider
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
	SetProviders(s ...ChannelProvider)
	SetMessagesService(s MessagesService)
	SetDevicesService(s DevicesService)
	SetLdmService(r LdmService)
	AppendRules(rules ...Rule)

	LogTag() LogTagFn
	MessagesService() MessagesService
	DevicesService() DevicesService
}
