package service

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

type MetricsRepository interface {
	UpdateGauge(context.Context, models.Metric) error
	UpdateCounter(context.Context, models.Metric) error
	UpdateAll(context.Context, []models.Metric) error
	Get(context.Context, models.Metric) (models.Metric, error)
	GetAll(context.Context) ([]models.Metric, error)
}

type MetricsStorage interface {
	Read() error
	Write() error
	CloseFile() error
}
