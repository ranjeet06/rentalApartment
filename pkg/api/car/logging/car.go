package car

import (
	"github.com/ribice/gorsk/pkg/api/car/service"
)

type LogService struct {
	service.Service
}

func New(svc service.Service) *LogService {
	return &LogService{
		Service: svc,
	}
}
