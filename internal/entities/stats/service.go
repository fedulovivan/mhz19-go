package stats

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var _ types.StatsService = (*service)(nil)

type service struct {
	repository StatsRepository
}

func (s service) Get() (res types.TableStats, err error) {
	res, err = s.repository.Get()
	if err != nil {
		return
	}
	return
}

func NewService(r StatsRepository) service {
	return service{
		repository: r,
	}
}
