package agentapiclient

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/Chystik/runtime-metrics/config"
)

func Example_reportMetrics() {
	var (
		gaugeValue = 123.20
		gaugeName  = "metric1"

		counterValue int64 = 456
		counterName        = "metric2"

		wrongCounterType = 456
	)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		_, err := rw.Write([]byte("ok"))
		if err != nil {
			panic(err)
		}
	}))
	defer server.Close()

	cfg := config.AgentConfig{Address: server.URL[7:]}
	client := New(server.Client(), &cfg)

	errGauge := client.ReportMetrics(context.Background(), map[string]interface{}{gaugeName: gaugeValue})
	errCounter := client.ReportMetrics(context.Background(), map[string]interface{}{counterName: counterValue})
	errWrongType := client.ReportMetrics(context.Background(), map[string]interface{}{counterName: wrongCounterType})

	fmt.Println(errGauge)
	fmt.Println(errCounter)
	fmt.Println(errWrongType)

	// Output:
	// <nil>
	// <nil>
	// unknown metric type: 456
}
