package adapters

import "github.com/Chystik/runtime-metrics/internal/models"

type AgentHTTPClient interface {
	ReportMetrics(metrics map[string]interface{}) error
	ReportMetricsJSON(metrics map[string]models.Metric) error
}
