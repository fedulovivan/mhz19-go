package stats

type StatsService interface {
	Get() (GetResult, error)
}

type GetResult struct {
	Rules    int32 `json:"rules"`
	Devices  int32 `json:"devices"`
	Messages int32 `json:"messages"`
}

type statsService struct {
	repository StatsRepository
}

func (s statsService) Get() (GetResult, error) {
	return s.repository.Get()
}

func NewService(r StatsRepository) StatsService {
	return statsService{
		repository: r,
	}
}
