package metricsservice

import (
	"context"

	"github.com/Chystik/runtime-metrics/internal/models"
	"github.com/Chystik/runtime-metrics/internal/service"
)

type metricsService struct {
	metricsRepo service.MetricsRepository
}

func New(mr service.MetricsRepository) *metricsService {
	return &metricsService{metricsRepo: mr}
}

func (ss *metricsService) UpdateGauge(ctx context.Context, metric models.Metric) error {
	return ss.metricsRepo.UpdateGauge(ctx, metric)
}

func (ss *metricsService) UpdateCounter(ctx context.Context, metric models.Metric) error {
	return ss.metricsRepo.UpdateCounter(ctx, metric)
}

func (ss *metricsService) Get(ctx context.Context, metric models.Metric) (models.Metric, error) {
	return ss.metricsRepo.Get(ctx, metric)
}

func (ss *metricsService) GetAll(ctx context.Context) ([]models.Metric, error) {
	return ss.metricsRepo.GetAll(ctx)
}

func (ss *metricsService) UpdateList(ctx context.Context, metrics []models.Metric) error {
	return ss.metricsRepo.UpdateList(ctx, metrics)
}
