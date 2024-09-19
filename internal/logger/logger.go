package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	"github.com/lmittmann/tint"
)

var seq = utils.NewSeq(0)

func Init() {
	if app.Config.IsDev {
		w := os.Stderr
		slog.SetDefault(slog.New(
			tint.NewHandler(w, &tint.Options{
				Level:      app.Config.LogLevel,
				TimeFormat: "15:04:05.000",
			}),
		))
	} else {
		slog.SetLogLoggerLevel(app.Config.LogLevel)
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	}
}

type TagName string

const (
	MAIN     TagName = "[main]      "
	ENGINE   TagName = "[engine]    "
	DB       TagName = "[db]        "
	REST     TagName = "[rest]      "
	ACTIONS  TagName = "[actions]   "
	ARGS     TagName = "[args]      "
	LDM      TagName = "[a_ldm]     "
	RULES    TagName = "[a_rules]   "
	STATS    TagName = "[a_stats]   "
	MESSAGES TagName = "[a_messages]"
	DEVICES  TagName = "[a_devices] "
	TBOT     TagName = "[p_tbot]    "
	DNSSD    TagName = "[p_dnssd]   "
	MQTT     TagName = "[p_mqtt]    "
	BURIED   TagName = "[p_buried]  "
)

type Tag interface {
	With(string, ...any) Tag
	WithTid() Tag
	F(format string, a ...any) string
}

type tag struct {
	tags []string
}

func NewTag(first TagName) Tag {
	return &tag{
		tags: []string{string(first)},
	}
}

func (t *tag) With(format string, a ...any) Tag {
	res := *t
	res.tags = append(res.tags, fmt.Sprintf(format, a...))
	return &res
}

func (t *tag) WithTid() Tag {
	return t.With("Tid#%v", seq.Inc())
}

func (t *tag) F(format string, a ...any) string {
	return strings.Join(
		append(t.tags, fmt.Sprintf(format, a...)),
		" ",
	)
}

// if app.Config.IsDev {
// 	// in development pad tag with spaces for extra nice output
// 	return func(message string) string {
// 		return fmt.Sprintf("%-10s", "["+tag+"]") + " " + message
// 	}
// }
