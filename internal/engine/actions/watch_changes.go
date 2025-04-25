package actions

import (
	"fmt"
	"sync"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/arguments"

	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

type Key struct {
	// TODO extend with other key components, look its not enough to bound time to device id only
	DeviceId types.DeviceId
}

type TimerWithValue struct {
	Timer *time.Timer
	Value any
	Since time.Time
}

var timersMu sync.Mutex
var timers = make(map[Key]TimerWithValue)

// var timeout = time.Second * 10

// args: Timeout, Value, From
var WatchChanges types.ActionImpl = func(
	compound types.MessageCompound,
	args types.Args,
	mapping types.Mapping,
	e types.ServiceAndProviderSupplier,
	tag utils.Tag,
) (err error) {

	m := compound.Curr

	if m == nil {
		return
	}

	key := Key{m.DeviceId}

	reader := arguments.NewReader(
		compound.Curr, args, mapping, nil, e, tag,
	)

	currValue := reader.Get("Value")
	ifValue := reader.Get("From")

	// TODO use same approach as in "type Throttle struct"
	timeoutFloat, err := arguments.GetTyped[float64](&reader, "Timeout")
	timeout := time.Second * time.Duration(timeoutFloat)

	timersMu.Lock()
	defer timersMu.Unlock()

	if prev, ok := timers[key]; ok { // stop and delete timer of value has changed
		if prev.Value != currValue {
			prev.Timer.Stop()
			delete(timers, key)
		}
	} else if ifValue == currValue { // run timer if Value equal to From
		timers[key] = TimerWithValue{
			Timer: time.AfterFunc(timeout, func() {
				timersMu.Lock()
				defer timersMu.Unlock()
				p := e.GetProvider(types.PROVIDER_SHIM_PROVIDER)
				p.Push(types.NewSystemMessage(
					fmt.Sprintf("No changes for last %s", timeout),
					types.DEVICE_ID_FOR_THE_WATCHER_MESSAGE,
				))
				delete(timers, key)
			}),
			Value: currValue,
			Since: time.Now(),
		}
	}
	return
}
