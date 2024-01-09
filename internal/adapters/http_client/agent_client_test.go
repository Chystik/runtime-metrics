package agentapiclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/models"

	"github.com/stretchr/testify/assert"
)

func Test_agentAPIClient_ReportMetrics(t *testing.T) {
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
				metricValue:   float64(123),
			},
			wantErr: false,
		},
		{
			name: "test counter",
			args: args{
				metricName:    "some",
				metricTypeStr: "counter",
				metricValue:   int64(123),
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

			err := client.ReportMetrics(context.Background(), map[string]interface{}{tt.args.metricName: tt.args.metricValue})

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func Test_agentAPIClient_ReportMetricsJSON(t *testing.T) {
	type args struct {
		metricName    string
		metricTypeStr string
		metricValue   float64
		metricDelat   int64
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
				metricValue:   float64(123),
			},
			wantErr: false,
		},
		{
			name: "test counter",
			args: args{
				metricName:    "test2",
				metricTypeStr: "counter",
				metricDelat:   int64(123),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := make(map[string]models.Metric)
			metrics[tt.args.metricName] = models.Metric{
				ID:    tt.args.metricName,
				MType: tt.args.metricTypeStr,
				Delta: &tt.args.metricDelat,
				Value: &tt.args.metricValue,
			}

			expURL := "/update/"

			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Test request parameters
				assert.Equal(t, req.URL.Path, expURL)
				// Send response to be tested
				b, err := io.ReadAll(req.Body)
				assert.NoError(t, err)

				_, err = rw.Write(b)
				assert.NoError(t, err)
			}))
			defer server.Close()

			cfg := config.AgentConfig{Address: server.URL[7:]}
			client := New(server.Client(), &cfg)

			err := client.ReportMetricsJSON(context.Background(), metrics)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func Test_agentAPIClient_ReportMetricsJSONBatch(t *testing.T) {
	type metric struct {
		metricName    string
		metricTypeStr string
		metricValue   float64
		metricDelat   int64
	}
	type args []metric

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test gauge",
			args: args{
				metric{
					metricName:    "test1",
					metricTypeStr: "gauge",
					metricValue:   float64(123),
				},
				metric{
					metricName:    "test2",
					metricTypeStr: "counter",
					metricDelat:   int64(456),
				},
			},
			wantErr: false,
		},
		{
			name: "test counter",
			args: args{
				metric{
					metricName:    "test3",
					metricTypeStr: "counter",
					metricDelat:   int64(55),
				},
				metric{
					metricName:    "test4",
					metricTypeStr: "gauge",
					metricValue:   float64(123),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := make(map[string]models.Metric)

			for i := range tt.args {
				metrics[tt.args[i].metricName] = models.Metric{
					ID:    tt.args[i].metricName,
					MType: tt.args[i].metricTypeStr,
					Delta: &tt.args[i].metricDelat,
					Value: &tt.args[i].metricValue,
				}
			}

			expURL := "/updates/"

			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Test request parameters
				assert.Equal(t, req.URL.Path, expURL)
				// Send response to be tested
				b, err := io.ReadAll(req.Body)
				assert.NoError(t, err)

				_, err = rw.Write(b)
				assert.NoError(t, err)
			}))
			defer server.Close()

			cfg := config.AgentConfig{Address: server.URL[7:]}
			client := New(server.Client(), &cfg)

			err := client.ReportMetricsBatch(context.Background(), metrics)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
