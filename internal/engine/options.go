package engine

type Options struct {
	logTag   func(m string) string
	services []Provider
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
	o.services = s
}

// start with defaults
var opts Options = NewOptions()
