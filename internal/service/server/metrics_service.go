package metricsservice

import (
	"context"
	"errors"

	memstorage "github.com/Chystik/runtime-metrics/internal/infrastructure/repository/mem_storage"
	"github.com/Chystik/runtime-metrics/internal/infrastructure/repository/postgres"
	"github.com/Chystik/runtime-metrics/internal/models"
)

type MetricsService interface {
	UpdateGauge(context.Context, models.Metric) error
	UpdateCounter(context.Context, models.Metric) error
	GetMetric(context.Context, models.Metric) (models.Metric, error)
	GetAllMetrics(context.Context) ([]models.Metric, error)
}

type metricsService struct {
	metricsRepo MetricsRepository
}

func New(mr MetricsRepository) *metricsService {
	return &metricsService{metricsRepo: mr}
}

func (ss *metricsService) UpdateGauge(ctx context.Context, metric models.Metric) error {
	return ss.metricsRepo.UpdateGauge(ctx, metric)
}

func (ss *metricsService) UpdateCounter(ctx context.Context, metric models.Metric) error {
	m, err := ss.metricsRepo.Get(ctx, metric)
	if errors.Is(err, memstorage.ErrNotFoundMetric) || errors.Is(err, postgres.ErrNotFoundMetric) {
		m = metric
	} else {
		*m.Delta = *m.Delta + *metric.Delta
	}
	return ss.metricsRepo.UpdateCounter(ctx, m)
}

func (ss *metricsService) GetMetric(ctx context.Context, metric models.Metric) (models.Metric, error) {
	return ss.metricsRepo.Get(ctx, metric)
}

func (ss *metricsService) GetAllMetrics(ctx context.Context) ([]models.Metric, error) {
	return ss.metricsRepo.GetAll(ctx)
}
