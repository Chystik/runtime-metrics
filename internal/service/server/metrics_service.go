package metricsservice

import "github.com/Chystik/runtime-metrics/internal/models"

type MetricsService interface {
	UpdateGauge(models.Metric)
	UpdateCounter(models.Metric)
	GetMetric(name string) (models.Metric, error)
	GetAllMetrics() []models.Metric
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
	m, err := ss.metricsRepo.Get(metric.Name)
	if err != nil {
		m.Name = metric.Name
		m.Counter = metric.Counter
	} else {
		m.Counter += metric.Counter
	}
	ss.metricsRepo.UpdateCounter(m)
}

func (ss *metricsService) GetMetric(name string) (models.Metric, error) {
	m, err := ss.metricsRepo.Get(name)
	if err != nil {
		return models.Metric{}, err
	}
	return m, nil
}

func (ss *metricsService) GetAllMetrics() []models.Metric {
	return ss.metricsRepo.GetAll()
}
