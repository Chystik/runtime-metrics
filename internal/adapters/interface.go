package adapters

import (
	"context"
	"net/http"

	"github.com/Chystik/runtime-metrics/internal/models"
)

type AgentHTTPClient interface {
	ReportMetrics(metrics map[string]interface{}) error
	ReportMetricsJSON(metrics map[string]models.Metric) error
}

type MetricsHandlers interface {
	UpdateMetric(w http.ResponseWriter, r *http.Request)
	GetMetric(w http.ResponseWriter, r *http.Request)
	UpdateMetricJSON(w http.ResponseWriter, r *http.Request)
	GetMetricJSON(w http.ResponseWriter, r *http.Request)
	AllMetrics(w http.ResponseWriter, r *http.Request)
}

type PgClient interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Ping(ctx context.Context) error
	PingHandler(w http.ResponseWriter, r *http.Request)
}
