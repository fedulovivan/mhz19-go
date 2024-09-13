package ldm

import (
	"fmt"
	"sort"
	"sync"

	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

var instance LdmRepository

type LdmRepository interface {
	NewKey(deviceClass types.DeviceClass, deviceId types.DeviceId) types.LdmKey
	Get(key types.LdmKey) types.Message
	Has(key types.LdmKey) bool
	Set(key types.LdmKey, m types.Message)
	GetAll() []types.Message
	OnSet() chan types.LdmKey
	GetByDeviceId(deviceId types.DeviceId) types.Message
}

var _ LdmRepository = (*repository)(nil)

type repository struct {
	// unsafe cache of device id to Key map, used only in GetByDeviceId
	// panics on key collision, see implementation in "Set"
	device_id_to_key_unsafemap map[types.DeviceId]types.LdmKey
	onset                      chan types.LdmKey
	data                       map[types.LdmKey]types.Message
	mu                         sync.RWMutex
}

func RepoSingleton() LdmRepository {
	if instance == nil {
		instance = &repository{
			device_id_to_key_unsafemap: make(map[types.DeviceId]types.LdmKey),
			onset:                      make(chan types.LdmKey, 100),
			data:                       make(map[types.LdmKey]types.Message),
		}
	}
	return instance
}

func (r *repository) NewKey(deviceClass types.DeviceClass, deviceId types.DeviceId) types.LdmKey {
	return types.LdmKey{
		DeviceClass: deviceClass,
		DeviceId:    deviceId,
	}
}

func (r *repository) GetAll() (result []types.Message) {
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

func (r *repository) GetByDeviceId(deviceId types.DeviceId) types.Message {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := r.device_id_to_key_unsafemap[deviceId]
	return r.get_unsafe(key)
}

func (r *repository) Get(key types.LdmKey) types.Message {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.get_unsafe(key)
}

func (r *repository) Has(key types.LdmKey) (flag bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, flag = r.data[key]
	return
}

// "private" getter not protected by lock
func (r *repository) get_unsafe(key types.LdmKey) types.Message {
	return r.data[key]
}

func (r *repository) OnSet() chan types.LdmKey {
	return r.onset
}

func (r *repository) Set(key types.LdmKey, m types.Message) {
	r.mu.Lock()
	defer r.mu.Unlock()
	existingkey, alreadyexist := r.device_id_to_key_unsafemap[m.DeviceId]
	if alreadyexist && existingkey != key {
		panic(fmt.Sprintf(
			"unsafemap key collision: same device id %v is referred by different keys: %v and %v",
			m.DeviceId,
			existingkey,
			key,
		))
	} else if !alreadyexist {
		r.device_id_to_key_unsafemap[m.DeviceId] = key
	}
	r.data[key] = m
	r.onset <- key
}
