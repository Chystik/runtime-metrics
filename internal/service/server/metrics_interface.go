package metricsservice

import (
	"context"

	"github.com/Chystik/runtime-metrics/internal/models"
)

type MetricsRepository interface {
	UpdateGauge(context.Context, models.Metric) error
	UpdateCounter(context.Context, models.Metric) error
	UpdateAll(context.Context, []models.Metric) error
	Get(context.Context, models.Metric) (models.Metric, error)
	GetAll(context.Context) ([]models.Metric, error)
}
