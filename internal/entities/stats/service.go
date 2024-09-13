package stats

import (
	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var _ types.StatsService = (*statsService)(nil)

type statsService struct {
	repository StatsRepository
}

func (s statsService) Get() (res types.StatsGetResult, err error) {
	res, err = s.repository.Get()
	if err != nil {
		return
	}
	res.InjectAppStats(app.StatsSingleton())
	return
}

func NewService(r StatsRepository) types.StatsService {
	return statsService{
		repository: r,
	}
}
