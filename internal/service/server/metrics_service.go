package metricsservice

import (
	"context"

	"github.com/Chystik/runtime-metrics/internal/models"
)

type MetricsService interface {
	UpdateGauge(context.Context, models.Metric) error
	UpdateCounter(context.Context, models.Metric) error
	UpdateAll(context.Context, []models.Metric) error
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
	return ss.metricsRepo.UpdateCounter(ctx, metric)
}

func (ss *metricsService) GetMetric(ctx context.Context, metric models.Metric) (models.Metric, error) {
	return ss.metricsRepo.Get(ctx, metric)
}

func (ss *metricsService) GetAllMetrics(ctx context.Context) ([]models.Metric, error) {
	return ss.metricsRepo.GetAll(ctx)
}

func (ss *metricsService) UpdateAll(ctx context.Context, metrics []models.Metric) error {
	return ss.metricsRepo.UpdateAll(ctx, metrics)
}
