package metricsservice

import (
	"errors"

	memstorage "github.com/Chystik/runtime-metrics/internal/infrastructure/repository/mem_storage"
	"github.com/Chystik/runtime-metrics/internal/models"
)

type MetricsService interface {
	UpdateGauge(models.Metric)
	UpdateCounter(models.Metric)
	GetMetric(Id string) (models.Metric, error)
	GetAllMetrics() []models.Metric
}

type metricsService struct {
	metricsRepo MetricsRepository
}

func New(mr MetricsRepository) *metricsService {
	return &metricsService{metricsRepo: mr}
}

func (ss *metricsService) UpdateGauge(metric models.Metric) {
	ss.metricsRepo.UpdateGauge(metric)
}

func (ss *metricsService) UpdateCounter(metric models.Metric) {
	m, err := ss.metricsRepo.Get(metric.ID)
	if errors.Is(err, memstorage.ErrNotFoundMetric) {
		m = metric
	} else {
		*m.Delta = *m.Delta + *metric.Delta
	}
	ss.metricsRepo.UpdateCounter(m)
}

func (ss *metricsService) GetMetric(Id string) (models.Metric, error) {
	return ss.metricsRepo.Get(Id)
}

func (ss *metricsService) GetAllMetrics() []models.Metric {
	return ss.metricsRepo.GetAll()
}
