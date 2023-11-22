package agentapiclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/models"
)

type metric struct {
	metricName    string
	metricTypeStr string
	metricValue   float64
	metricDelat   int64
}

func Example_reportMetricsJSON() {
	m := metric{
		metricName:    "test1",
		metricTypeStr: "gauge",
		metricValue:   float64(123),
	}

	metrics := make(map[string]models.Metric)

	metrics[m.metricName] = models.Metric{
		ID:    m.metricName,
		MType: m.metricTypeStr,
		Value: &m.metricValue,
		Delta: &m.metricDelat,
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		b, err := io.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}

		_, err = rw.Write(b)
		if err != nil {
			panic(err)
		}
	}))
	defer server.Close()

	cfg := config.AgentConfig{Address: server.URL[7:]}
	client := New(server.Client(), &cfg)

	err := client.ReportMetricsJSON(context.Background(), metrics)

	fmt.Println(err) // Output: <nil>
}
