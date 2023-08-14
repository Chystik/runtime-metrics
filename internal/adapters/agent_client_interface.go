package adapters

type AgentHTTPClient interface {
	ReportMetrics(metrics map[string]interface{}) error
}
