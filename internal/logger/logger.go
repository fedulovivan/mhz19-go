package logger

import (
	"log"
	"log/slog"
	"os"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	"github.com/lmittmann/tint"
)

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
	MAIN     utils.TagName = "[main]      "
	ENGINE   utils.TagName = "[engine]    "
	DB       utils.TagName = "[db]        "
	REST     utils.TagName = "[rest]      "
	ACTIONS  utils.TagName = "[actions]   "
	CONDS    utils.TagName = "[conds]     "
	ARGS     utils.TagName = "[args]      "
	LDM      utils.TagName = "[a_ldm]     "
	RULES    utils.TagName = "[a_rules]   "
	STATS    utils.TagName = "[a_stats]   "
	MESSAGES utils.TagName = "[a_messages]"
	DEVICES  utils.TagName = "[a_devices] "
	DICTS    utils.TagName = "[a_dicts]   "
	TBOT     utils.TagName = "[p_tbot]    "
	DNSSD    utils.TagName = "[p_dnssd]   "
	MQTT     utils.TagName = "[p_mqtt]    "
	BURIED   utils.TagName = "[p_buried]  "
)

// if app.Config.IsDev {
// 	// in development pad tag with spaces for extra nice output
// 	return func(message string) string {
// 		return fmt.Sprintf("%-10s", "["+tag+"]") + " " + message
// 	}
// }
