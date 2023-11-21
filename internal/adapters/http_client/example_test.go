package agentapiclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/Chystik/runtime-metrics/config"
)

func Example_agentAPIClient_ReportMetrics() {
	gaugeValue := 123.20
	gaugeName := "metric1"

	var counterValue int64 = 456
	counterName := "metric2"

	wrongCounterType := 456

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		_, _ = rw.Write([]byte(`OK`))
	}))
	defer server.Close()

	cfg := config.AgentConfig{Address: server.URL[7:]}
	client := New(server.Client(), &cfg)

	errGauge := client.ReportMetrics(map[string]interface{}{gaugeName: gaugeValue})
	errCounter := client.ReportMetrics(map[string]interface{}{counterName: counterValue})
	errWrongType := client.ReportMetrics(map[string]interface{}{counterName: wrongCounterType})

	fmt.Println(errGauge)
	fmt.Println(errCounter)
	fmt.Println(errWrongType)

	// Output:
	// <nil>
	// <nil>
	// unknown metric type: 456
}
