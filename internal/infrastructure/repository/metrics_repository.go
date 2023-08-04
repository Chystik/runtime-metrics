package repository

import "github.com/Chystik/runtime-metrics/internal/models"

type MetricsRepository interface {
	UpdateGauge(models.Metric)
	UpdateCounter(models.Metric)
	Get(name string) (models.Metric, error)
}
