package last_device_message

import (
	"fmt"
	"sort"
	"sync"

	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

var instance LdmRepository

type LdmRepository interface {
	MakeKey(deviceClass types.DeviceClass, deviceId types.DeviceId) types.LdmKey
	Get(key types.LdmKey) types.Message
	Set(key types.LdmKey, m types.Message)
	GetAll() []types.Message
	GetByDeviceId(deviceId types.DeviceId) types.Message
}

type repository struct {
	mu   sync.RWMutex
	data map[types.LdmKey]types.Message
	// unsafe cache of device id to Key map, used only in GetByDeviceId
	// panics on key collision, see implementation in "Set"
	device_id_to_key_unsafemap map[types.DeviceId]types.LdmKey
}

func RepositoryInstance() LdmRepository {
	if instance == nil {
		instance = &repository{
			data:                       make(map[types.LdmKey]types.Message),
			device_id_to_key_unsafemap: make(map[types.DeviceId]types.LdmKey),
		}
	}
	return instance
}

func (r *repository) MakeKey(deviceClass types.DeviceClass, deviceId types.DeviceId) types.LdmKey {
	return types.LdmKey(types.DEVICE_CLASS_NAMES[deviceClass] + "-" + string(deviceId))
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

// "private" getter not protected by lock
func (r *repository) get_unsafe(key types.LdmKey) types.Message {
	return r.data[key]
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
}
