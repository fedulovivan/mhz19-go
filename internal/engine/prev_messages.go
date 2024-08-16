package engine

import "sync"

var prevMessages map[string]Message

var lock sync.RWMutex

func makeMessageKey(dc DeviceClass, deviceId DeviceId) string {
	return string(dc) + "-" + string(deviceId)
}

func PrevMessageGet(key string) Message {
	lock.RLock()
	defer lock.RUnlock()
	return prevMessages[key]
}

func PrevMessagePut(key string, m Message) {
	lock.Lock()
	defer lock.Unlock()
	prevMessages[key] = m
}
