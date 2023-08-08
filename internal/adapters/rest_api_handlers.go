package adapters

import (
	handlers "github.com/Chystik/runtime-metrics/internal/adapters/rest_api_handlers"
	metricsservice "github.com/Chystik/runtime-metrics/internal/service/server"
)

type ServerHandlers interface {
	handlers.MetricsHandlers
}

type serverHandlers struct {
	handlers.MetricsHandlers
}

func NewServerHandlers(ms metricsservice.MetricsService) ServerHandlers {
	return &serverHandlers{
		handlers.NewMetricsHandlers(ms),
	}
}
