package tbot_provider

import (
	"fmt"
	"log/slog"
)

type slogAdapter struct{}

func (l slogAdapter) Println(v ...any) {
	switch v0 := v[0].(type) {
	case string:
		slog.Debug(logTag(v0), "more", len(v)-1)
	case error:
		slog.Error(logTag(v0.Error()), "more", len(v)-1)
	default:
		slog.Error(logTag(fmt.Sprintf(
			"slogAdapter.Println() skipped, its first argument expected to be a string, but got %T with value %v",
			v[0], v[0],
		)))
	}
}
func (l slogAdapter) Printf(format string, v ...any) {
	last := format[len(format)-1]
	var nl byte = 10 /* \n */
	if last == nl {
		format = format[:len(format)-1]
	}
	slog.Debug(logTag(fmt.Sprintf(format, v...)))
}
