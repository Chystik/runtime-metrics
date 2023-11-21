package service

import (
	"context"

	"github.com/Chystik/runtime-metrics/internal/models"
)

type AgentAPIClient interface {
	ReportMetrics(metrics map[string]interface{}) error
	ReportMetricsJSON(metrics map[string]models.Metric) error
	ReportMetricsJSONBatch(ctx context.Context, metrics map[string]models.Metric) error
}
