package engine

type Options struct {
	logTag   func(m string) string
	services []Service
}

func NewOptions() Options {
	return Options{
		logTag: func(m string) string { return "" },
	}
}

func (o *Options) SetLogTag(f func(m string) string) {
	o.logTag = f
}
func (o *Options) SetServices(s ...Service) {
	o.services = s
}

// start with defaults
var opts Options = NewOptions()
