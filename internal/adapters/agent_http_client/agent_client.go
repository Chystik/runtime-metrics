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
	client *http.Client
	server config.HTTPServer
}

func New(c *http.Client, s config.HTTPServer) AgentHTTPClient {
	return &agentHTTPClient{
		client: c,
		server: s,
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

		url := fmt.Sprintf("http://%s:%d/update/%s/%s/%v", ac.server.Host, ac.server.Port, mType, name, value)

		request, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			panic(err)
		}
		request.Header.Set("Content-Type", "text/plain")
		response, err := ac.client.Do(request)
		if err != nil {
			panic(err)
		}
		//io.Copy(os.Stdout, response.Body)
		response.Body.Close()
	}
	return nil
}
