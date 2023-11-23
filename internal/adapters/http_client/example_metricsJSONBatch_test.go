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

type m struct {
	metricName    string
	metricTypeStr string
	metricValue   float64
	metricDelat   int64
}

type metrics []m

func Example_reportMetricsJSONBatch() {
	ms := metrics{
		m{
			metricName:    "test1",
			metricTypeStr: "gauge",
			metricValue:   float64(123),
		},
		m{
			metricName:    "test2",
			metricTypeStr: "counter",
			metricDelat:   int64(456),
		},
	}

	metrics := make(map[string]models.Metric)

	for i := range ms {
		metrics[ms[i].metricName] = models.Metric{
			ID:    ms[i].metricName,
			MType: ms[i].metricTypeStr,
			Value: &ms[i].metricValue,
			Delta: &ms[i].metricDelat,
		}

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
