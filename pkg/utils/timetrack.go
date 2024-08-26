package utils

import (
	"fmt"
	"log/slog"
	"time"
)

// from https://coderwall.com/p/cp5fya/measuring-execution-time-in-go
func TimeTrack(logTag func(m string) string, start time.Time, name string) {
	elapsed := time.Since(start)
	slog.Debug(logTag(fmt.Sprintf("%s took %s", name, elapsed)))
}
