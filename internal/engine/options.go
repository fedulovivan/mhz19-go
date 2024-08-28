package engine

type Options struct {
	logTag         func(m string) string
	providers      []Provider
	rules          []Rule
	messageService MessagesService
}

func NewOptions() Options {
	return Options{
		logTag: func(m string) string { return "" },
	}
}

func (o *Options) SetLogTag(f func(m string) string) {
	o.logTag = f
}
func (o *Options) SetProviders(s ...Provider) {
	o.providers = s
}
func (o *Options) SetMessagesService(s MessagesService) {
	o.messageService = s
}
func (o *Options) SetRules(rr []Rule) {
	o.rules = rr
}
