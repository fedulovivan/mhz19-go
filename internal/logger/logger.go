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
	MAIN     TagName = "[main]    "
	ENGINE   TagName = "[engine]  "
	MQTT     TagName = "[mqtt]    "
	TBOT     TagName = "[tbot]    "
	DB       TagName = "[db]      "
	REST     TagName = "[rest]    "
	RULES    TagName = "[rules]   "
	STATS    TagName = "[stats]   "
	LDM      TagName = "[ldm]     "
	MESSAGES TagName = "[messages]"
	DEVICES  TagName = "[devices] "
	DNSSD    TagName = "[dns-sd]  "
	ACTIONS  TagName = "[actions] "
	BURIED   TagName = "[buried]  "
)

type Tag interface {
	With(string, ...any) Tag
	WithTid() Tag
	F(string) string
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

func (t *tag) F(message string) string {
	return strings.Join(
		append(t.tags, message),
		" ",
	)
}
