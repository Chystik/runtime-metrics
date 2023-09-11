package adapters

import (
	"context"
	"net/http"

	"github.com/Chystik/runtime-metrics/internal/models"
	"github.com/jmoiron/sqlx"
)

type AgentHTTPClient interface {
	ReportMetrics(metrics map[string]interface{}) error
	ReportMetricsJSON(metrics map[string]models.Metric) error
	ReportMetricsJSONBatch(metrics map[string]models.Metric) error
}

type MetricsHandlers interface {
	GetMetric(w http.ResponseWriter, r *http.Request)
	GetMetricJSON(w http.ResponseWriter, r *http.Request)
	UpdateMetric(w http.ResponseWriter, r *http.Request)
	UpdateMetricJSON(w http.ResponseWriter, r *http.Request)
	UpdateMetricsJSON(w http.ResponseWriter, r *http.Request)
	AllMetrics(w http.ResponseWriter, r *http.Request)
}

type PgClient interface {
	Connect(ctx context.Context) (*sqlx.DB, error)
	Disconnect(ctx context.Context) error
	Migrate() error
	Ping(ctx context.Context) error
	PingHandler(w http.ResponseWriter, r *http.Request)
}
