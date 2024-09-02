package engine

import (
	"sync"

	"github.com/fedulovivan/mhz19-go/internal/types"
)

var prevMessages map[string]types.Message

var lock sync.RWMutex

func makeMessageKey(dc types.DeviceClass, deviceId types.DeviceId) string {
	return string(dc) + "-" + string(deviceId)
}

func PrevMessageGet(key string) types.Message {
	lock.RLock()
	defer lock.RUnlock()
	return prevMessages[key]
}

func PrevMessagePut(key string, m types.Message) {
	lock.Lock()
	defer lock.Unlock()
	prevMessages[key] = m
}
