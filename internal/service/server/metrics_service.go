package metricsservice

import (
	"github.com/Chystik/runtime-metrics/internal/infrastructure/repository"
	"github.com/Chystik/runtime-metrics/internal/models"
)

type MetricsService interface {
	UpdateGauge(models.Metric)
	UpdateCounter(models.Metric)
}

type metricsService struct {
	metricsRepo repository.MetricsRepository
}

func New(mr repository.MetricsRepository) MetricsService {
	return &metricsService{metricsRepo: mr}
}

func (ss *metricsService) UpdateGauge(metric models.Metric) {
	ss.metricsRepo.UpdateGauge(metric)
}

func (ss *metricsService) UpdateCounter(metric models.Metric) {
	val, err := ss.metricsRepo.Get(metric.Name)
	if err != nil {
		return
	}

	val.Counter += metric.Counter
	ss.metricsRepo.UpdateCounter(val)
}
