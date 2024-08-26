package logger

import (
	"log"
	"log/slog"
	"os"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/lmittmann/tint"
)

func Init() {
	if app.Config.IsDev {
		w := os.Stderr
		slog.SetDefault(slog.New(
			tint.NewHandler(w, &tint.Options{
				Level:      app.Config.LogLevel,
				TimeFormat: "15:04:05.000",
				// TimeFormat: time.TimeOnly,
			}),
		))
	} else {
		slog.SetLogLoggerLevel(app.Config.LogLevel)
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	}
}

type TagName string

const (
	MAIN   TagName = "[main]   "
	ENGINE TagName = "[engine] "
	MQTT   TagName = "[mqtt]   "
	TBOT   TagName = "[tbot]   "
	MODEL  TagName = "[model]  "
	REST   TagName = "[rest]   "
	RULES  TagName = "[rules]  "
)

func MakeTag(tag TagName) func(m string) string {
	return func(message string) string {
		return string(tag) + " " + message
	}
}

// if app.Config.IsDev {
// 	// in development pad tag with spaces for extra nice output
// 	return func(message string) string {
// 		return fmt.Sprintf("%-10s", "["+tag+"]") + " " + message
// 	}
// }
