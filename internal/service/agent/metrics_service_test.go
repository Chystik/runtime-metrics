package agentservice

import (
	"runtime"
	"testing"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/adapters"
	"github.com/Chystik/runtime-metrics/internal/adapters/agent_http_client/mocks"
	"github.com/Chystik/runtime-metrics/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_New(t *testing.T) {
	var c adapters.AgentHTTPClient
	var m []string

	agentService := New(c, m)

	assert.NotNil(t, agentService)
}

func Test_agentService_UpdateMetrics(t *testing.T) {
	cache := make(map[string]models.Metric)
	cache["PollCount"] = models.Metric{ID: "PollCount", MType: "counter", Delta: new(int64)}

	tests := []struct {
		name string
		as   *agentService
	}{
		{
			name: "get metrics",
			as: &agentService{
				collectableMetrics: []string{"Alloc", "BuckHashSys", "Frees"},
				runtimeMetrics:     runtime.MemStats{},
				cache:              cache,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.as.UpdateMetrics()
		})
	}
}

func TestReportMetrics_WhenClientRetunNoError(t *testing.T) {
	c, mks := getAgentServiceMocks()

	mks.client.On("ReportMetricsJSON", mock.Anything).Return(nil)
	c.ReportMetrics()
}

type agentServiceMocks struct {
	client *mocks.AgentHTTPClient
}

func getAgentServiceMocks() (AgentService, *agentServiceMocks) {
	mks := &agentServiceMocks{
		client: &mocks.AgentHTTPClient{},
	}

	as := New(mks.client, config.CollectableMetrics{})
	return as, mks
}
