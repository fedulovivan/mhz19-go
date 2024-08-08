package logger

import (
	"log"
	"log/slog"
	"os"

	"github.com/fedulovivan/mhz19-go/internal/registry"
	"github.com/lmittmann/tint"
)

func init() {
	if registry.Config.IsDev {
		w := os.Stderr
		slog.SetDefault(slog.New(
			tint.NewHandler(w, &tint.Options{
				Level:      registry.Config.LogLevel,
				TimeFormat: "15:04:05.000",
				// TimeFormat: time.TimeOnly,
			}),
		))
	} else {
		slog.SetLogLoggerLevel(registry.Config.LogLevel)
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	}
}

func MakeTag(tag string) func(m string) string {
	return func(message string) string {
		return "[" + tag + "]" + " " + message
	}
}
