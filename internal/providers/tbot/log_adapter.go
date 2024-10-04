package tbot_provider

import (
	"fmt"
	"log/slog"

	"github.com/fedulovivan/mhz19-go/internal/counters"
)

type slogAdapter struct{}

func (l slogAdapter) Println(v ...any) {
	switch v0 := v[0].(type) {
	case string:
		slog.Debug(tag.F(v0), "more", len(v)-1)
	case error:
		slog.Error(tag.F(v0.Error()), "more", len(v)-1)
		counters.Inc(counters.ERRORS_ALL)
	default:
		slog.Error(tag.F(
			"slogAdapter.Println() skipped, its first argument expected to be a string, but got %T with value %v",
			v[0], v[0],
		))
		counters.Inc(counters.ERRORS_ALL)
	}
}
func (l slogAdapter) Printf(format string, v ...any) {
	last := format[len(format)-1]
	var nl byte = 10 /* \n */
	if last == nl {
		format = format[:len(format)-1]
	}
	slog.Debug(tag.F(fmt.Sprintf(format, v...)))
}
