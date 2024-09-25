package dicts

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var _ types.DictsService = (*service)(nil)

type service struct {
	repository DictsRepository
}

func NewService(r DictsRepository) types.DictsService {
	return service{
		repository: r,
	}
}

func BuildItems(in []DbDictItem) (res []types.DictItem) {
	res = make([]types.DictItem, len(in))
	for i, item := range in {
		res[i] = types.DictItem{
			Id:   int(item.Id),
			Name: item.Name,
		}
	}
	return
}

func (s service) Get(dtype types.DictType) (res []types.DictItem, err error) {
	data, err := s.repository.Get(dtype)
	if err != nil {
		return
	}
	return BuildItems(data), nil
}

func (s service) All() (res map[types.DictType][]types.DictItem, err error) {
	actions,
		conditions,
		channels,
		deviceClasses,
		err := s.repository.All()
	if err != nil {
		return
	}
	res = map[types.DictType][]types.DictItem{
		types.DICT_ACTIONS:        BuildItems(actions),
		types.DICT_CONDITIONS:     BuildItems(conditions),
		types.DICT_CHANNELS:       BuildItems(channels),
		types.DICT_DEVICE_CLASSES: BuildItems(deviceClasses),
	}
	return
}
