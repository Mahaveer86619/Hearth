package services

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (hs *HealthService) Ping() string {
	return "pong"
}

func (hs *HealthService) Check() string {
	return "ok"
}
