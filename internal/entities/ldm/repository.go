package ldm

import (
	"fmt"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/app"
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
	Buried() chan types.LdmKey
	GetByDeviceId(deviceId types.DeviceId) types.Message
	AppendBuriedBlacklist(...types.LdmKey)
}

var _ LdmRepository = (*repository)(nil)

type repository struct {

	// unsafe cache of device id to Key map, used only in GetByDeviceId
	// panics on key collision, see implementation in "Set"
	device_id_to_key_unsafemap map[types.DeviceId]types.LdmKey

	data             map[types.LdmKey]types.Message
	buried_timers    map[types.LdmKey]*time.Timer
	buried_chan      chan types.LdmKey
	buried_blacklist []types.LdmKey

	mu sync.RWMutex
}

func RepoSingleton() LdmRepository {
	if instance == nil {
		instance = &repository{
			data:                       make(map[types.LdmKey]types.Message),
			buried_chan:                make(chan types.LdmKey, 100),
			buried_timers:              make(map[types.LdmKey]*time.Timer),
			device_id_to_key_unsafemap: make(map[types.DeviceId]types.LdmKey),
		}
	}
	return instance
}

func (r *repository) AppendBuriedBlacklist(keys ...types.LdmKey) {
	r.buried_blacklist = append(r.buried_blacklist, keys...)
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

func (r *repository) Buried() chan types.LdmKey {
	return r.buried_chan
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

	// implementation of "buried devices" aka "have not seen for a while" feature
	if !slices.Contains(r.buried_blacklist, key) {
		if timer, ok := r.buried_timers[key]; ok {
			timer.Reset(app.Config.BuriedTimeout)
		} else {
			r.buried_timers[key] = time.AfterFunc(
				app.Config.BuriedTimeout,
				func() {
					r.buried_chan <- key
				},
			)
		}
	}

}
