package utils

import (
	"fmt"
	"log/slog"
	"time"
)

// initially from https://coderwall.com/p/cp5fya/measuring-execution-time-in-go
// + extra goodies
func TimeTrack(logTag func(format string, a ...any) string, start time.Time, name string) {
	elapsed := time.Since(start)
	badlySlow := elapsed > time.Second*1
	if logTag == nil {
		logTag = fmt.Sprintf
	}
	m := logTag("%s took %s", name, elapsed)
	if badlySlow {
		slog.Warn(m)
	} else {
		slog.Debug(m)
	}
}

// elapsed := time.Since(start)
// wayFast := elapsed < time.Millisecond*5
// badlySlow := elapsed > time.Second*1
// badge := ""
// if wayFast {
// 	badge = "✨ "
// } else if badlySlow {
// 	badge = "🧨 "
// }
// m := logTag("%v%s took %s", badge, name, elapsed)
// if badlySlow {
// 	slog.Warn(m)
// } else {
// 	slog.Debug(m)
// }
