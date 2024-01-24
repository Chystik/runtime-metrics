package service

import (
	"context"

	"github.com/Chystik/runtime-metrics/internal/models"
)

type AgentAPIClient interface {
	ReportMetrics(ctx context.Context, metrics map[string]interface{}) error
	ReportMetricsJSON(ctx context.Context, metrics map[string]models.Metric) error
	ReportMetricsBatch(ctx context.Context, metrics map[string]models.Metric) error
}

type AgentService interface {
	UpdateMetrics()
	UpdateGoPsUtilMetrics()
	ReportMetrics(context.Context) error
}
