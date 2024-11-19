package ldm

import (
	"fmt"
	"sort"
	"sync"

	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

var instance *repo

type LdmRepository interface {
	NewKey(deviceClass types.DeviceClass, deviceId types.DeviceId) types.LdmKey
	Get(key types.LdmKey) types.Message
	Has(key types.LdmKey) bool
	Set(key types.LdmKey, m types.Message)
	GetAll() []types.Message
	OnSet() chan types.LdmKey
	GetByDeviceId(deviceId types.DeviceId) (types.Message, error)
}

var _ LdmRepository = (*repo)(nil)

type repo struct {
	// unsafe (unreliable to be correct) cache of device id to Key map, used only in GetByDeviceId
	// panics on key collision, see implementation in "Set"
	device_id_to_key_unsafemap map[types.DeviceId]types.LdmKey
	onset                      chan types.LdmKey
	data                       map[types.LdmKey]types.Message
	mu                         sync.RWMutex
}

func RepoSingleton() *repo {
	if instance == nil {
		instance = &repo{
			device_id_to_key_unsafemap: make(map[types.DeviceId]types.LdmKey),
			onset:                      make(chan types.LdmKey),
			data:                       make(map[types.LdmKey]types.Message),
		}
	}
	return instance
}

func (r *repo) NewKey(deviceClass types.DeviceClass, deviceId types.DeviceId) types.LdmKey {
	return types.LdmKey{
		DeviceClass: deviceClass,
		DeviceId:    deviceId,
	}
}

func (r *repo) GetAll() (result []types.Message) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result = utils.Values(r.data)
	sort.Slice(
		result,
		func(i, j int) bool {
			return result[i].Timestamp.After(result[j].Timestamp)
		},
	)
	return
}

func (r *repo) GetByDeviceId(deviceId types.DeviceId) (res types.Message, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key, exist := r.device_id_to_key_unsafemap[deviceId]
	if !exist {
		err = fmt.Errorf("no last message for %s", deviceId)
		return
	}
	return r.get_unsafe(key), nil
}

func (r *repo) Get(key types.LdmKey) types.Message {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.get_unsafe(key)
}

func (r *repo) Has(key types.LdmKey) (flag bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, flag = r.data[key]
	return
}

// "private" getter not protected by lock
func (r *repo) get_unsafe(key types.LdmKey) types.Message {
	return r.data[key]
}

func (r *repo) OnSet() chan types.LdmKey {
	return r.onset
}

func (r *repo) Set(newkey types.LdmKey, m types.Message) {
	r.mu.Lock()
	defer r.mu.Unlock()
	existingkey, alreadyexist := r.device_id_to_key_unsafemap[m.DeviceId]
	if alreadyexist && existingkey != newkey {
		panic(fmt.Sprintf(
			"device_id_to_key_unsafemap key collision: device id %v is referred by several keys: %v and %v",
			m.DeviceId,
			existingkey,
			newkey,
		))
	} else if !alreadyexist {
		r.device_id_to_key_unsafemap[m.DeviceId] = newkey
	}
	r.data[newkey] = m
	// this blocks now, if no receiver for onset is registered
	r.onset <- newkey
}
