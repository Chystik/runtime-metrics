package agenthttpclient

import (
	"fmt"
	"net/http"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/models"
)

type AgentHTTPClient interface {
	ReportMetrics(metrics map[string]interface{}) error
}

type agentHTTPClient struct {
	client  *http.Client
	address string
}

func New(c *http.Client, s *config.AgentConfig) *agentHTTPClient {
	return &agentHTTPClient{
		client:  c,
		address: s.Address,
	}
}

func (ac *agentHTTPClient) ReportMetrics(metrics map[string]interface{}) error {
	for name, value := range metrics {
		var mType string

		switch t := value.(type) {
		case models.Gauge:
			mType = "gauge"
		case models.Counter:
			mType = "counter"
		default:
			return fmt.Errorf("unknown metric type: %#v", t)
		}

		url := fmt.Sprintf("http://%s/update/%s/%s/%v", ac.address, mType, name, value)

		request, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			panic(err)
		}
		request.Header.Set("Content-Type", "text/plain")
		response, err := ac.client.Do(request)
		if err != nil {
			panic(err)
		}
		response.Body.Close()
	}
	return nil
}
