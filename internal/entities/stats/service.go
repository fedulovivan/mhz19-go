package stats

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var _ types.StatsService = (*statsService)(nil)

type statsService struct {
	repository StatsRepository
}

func (s statsService) Get() (res types.TableStats, err error) {
	res, err = s.repository.Get()
	if err != nil {
		return
	}
	return
}

func NewService(r StatsRepository) types.StatsService {
	return statsService{
		repository: r,
	}
}
