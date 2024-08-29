package utils

import (
	"fmt"
	"log/slog"
	"time"
)

// inially from https://coderwall.com/p/cp5fya/measuring-execution-time-in-go
// + extra goodies
func TimeTrack(logTag func(m string) string, start time.Time, name string) {
	elapsed := time.Since(start)
	wayFast := elapsed < time.Millisecond*10
	badlySlow := elapsed > time.Second*1
	badge := ""
	if wayFast {
		badge = "✨ "
	} else if badlySlow {
		badge = "🧨 "
	}
	m := logTag(fmt.Sprintf("%v%s took %s", badge, name, elapsed))
	if badlySlow {
		slog.Warn(m)
	} else {
		slog.Debug(m)
	}
}
