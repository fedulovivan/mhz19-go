package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	"github.com/lmittmann/tint"
)

var nsSequences = make(map[string]utils.Seq)
var secsMu = new(sync.Mutex)

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
	MOD_MAIN     string = "main"
	MOD_ENGINE   string = "engine"
	MOD_DB       string = "db"
	MOD_REST     string = "rest"
	MOD_ACTIONS  string = "actions"
	MOD_CONDS    string = "conds"
	MOD_ARGS     string = "args"
	MOD_LDM      string = "a_ldm"
	MOD_RULES    string = "a_rules"
	MOD_STATS    string = "a_stats"
	MOD_MESSAGES string = "a_messages"
	MOD_DEVICES  string = "a_devices"
	MOD_DICTS    string = "a_dicts"
	MOD_TBOT     string = "p_tbot"
	MOD_DNSSD    string = "p_dnssd"
	MOD_MQTT     string = "p_mqtt"
	MOD_BURIED   string = "p_buried"
)

const (
	MAIN     TagName = "[main]      "
	ENGINE   TagName = "[engine]    "
	DB       TagName = "[db]        "
	REST     TagName = "[rest]      "
	ACTIONS  TagName = "[actions]   "
	CONDS    TagName = "[conds]     "
	ARGS     TagName = "[args]      "
	LDM      TagName = "[a_ldm]     "
	RULES    TagName = "[a_rules]   "
	STATS    TagName = "[a_stats]   "
	MESSAGES TagName = "[a_messages]"
	DEVICES  TagName = "[a_devices] "
	DICTS    TagName = "[a_dicts]   "
	TBOT     TagName = "[p_tbot]    "
	DNSSD    TagName = "[p_dnssd]   "
	MQTT     TagName = "[p_mqtt]    "
	BURIED   TagName = "[p_buried]  "
)

type Tag interface {
	With(string, ...any) Tag
	WithTid(string) Tag
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
	tagscopy := make([]string, len(t.tags)+1)
	copy(tagscopy, t.tags)
	tagscopy[len(t.tags)] = fmt.Sprintf(format, a...)
	return &tag{
		tags: tagscopy,
	}
}

func (t *tag) WithTid(ns string) Tag {
	secsMu.Lock()
	defer secsMu.Unlock()
	if _, exist := nsSequences[ns]; !exist {
		nsSequences[ns] = utils.NewSeq(0)
	}
	return t.With("%s#%v", ns, nsSequences[ns].Inc())
}

func (t *tag) F(format string, a ...any) string {
	// fmt.Println("tags:", t.tags)
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
