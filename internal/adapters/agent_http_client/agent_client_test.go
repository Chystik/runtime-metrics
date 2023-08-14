package agenthttpclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_agentHTTPClient_ReportMetrics(t *testing.T) {
	type args struct {
		metricName    string
		metricTypeStr string
		metricValue   interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test gauge",
			args: args{
				metricName:    "test1",
				metricTypeStr: "gauge",
				metricValue:   models.Gauge(123),
			},
			wantErr: false,
		},
		{
			name: "test counter",
			args: args{
				metricName:    "some",
				metricTypeStr: "counter",
				metricValue:   models.Counter(123),
			},
			wantErr: false,
		},
		{
			name: "test wrong metric type",
			args: args{
				metricName:    "some",
				metricTypeStr: "wrong",
				metricValue:   123,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expURL := fmt.Sprintf("/update/%s/%s/%v", tt.args.metricTypeStr, tt.args.metricName, tt.args.metricValue)

			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Test request parameters
				assert.Equal(t, req.URL.Path, expURL)
				// Send response to be tested
				_, err := rw.Write([]byte(`OK`))
				assert.NoError(t, err)
			}))
			defer server.Close()

			cfg := config.AgentConfig{Address: server.URL[7:]}
			client := New(server.Client(), &cfg)

			err := client.ReportMetrics(map[string]interface{}{tt.args.metricName: tt.args.metricValue})

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
