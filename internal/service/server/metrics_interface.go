package metricsservice

import "github.com/Chystik/runtime-metrics/internal/models"

type MetricsRepository interface {
	UpdateGauge(models.Metric)
	UpdateCounter(models.Metric)
	Get(models.Metric) (models.Metric, error)
	GetAll() []models.Metric
}
