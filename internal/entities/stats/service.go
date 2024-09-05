package stats

import "github.com/fedulovivan/mhz19-go/internal/types"

type statsService struct {
	repository StatsRepository
}

func (s statsService) Get() (types.StatsGetResult, error) {
	return s.repository.Get()
}

func NewService(r StatsRepository) types.StatsService {
	return statsService{
		repository: r,
	}
}
