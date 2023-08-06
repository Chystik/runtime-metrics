package metricsservice

import "github.com/Chystik/runtime-metrics/internal/models"

type MetricsService interface {
	UpdateGauge(models.Metric)
	UpdateCounter(models.Metric)
}

type metricsService struct {
	metricsRepo MetricsRepository
}

func New(mr MetricsRepository) MetricsService {
	return &metricsService{metricsRepo: mr}
}

func (ss *metricsService) UpdateGauge(metric models.Metric) {
	ss.metricsRepo.UpdateGauge(metric)
}

func (ss *metricsService) UpdateCounter(metric models.Metric) {
	val, err := ss.metricsRepo.Get(metric.Name)
	if err != nil {
		val.Name = metric.Name
		val.Counter = metric.Counter
	} else {
		val.Counter += metric.Counter
	}
	ss.metricsRepo.UpdateCounter(val)
}
