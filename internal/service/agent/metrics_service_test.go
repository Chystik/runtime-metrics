package agentservice

import (
	"context"
	"testing"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/service"
	"github.com/Chystik/runtime-metrics/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_New(t *testing.T) {
	var c service.AgentAPIClient
	var m []string

	agentService := New(c, m)

	assert.NotNil(t, agentService)
}

func Test_agentService_UpdateMetrics(t *testing.T) {
	collectableMetrics := []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "NumGC"}
	as := New(nil, collectableMetrics)

	tests := []struct {
		name string
		as   *agentService
	}{
		{
			name: "get metrics",
			as:   as,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.as.UpdateMetrics()
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.as.UpdateGoPsUtilMetrics()
		})
	}
}

func TestReportMetrics_WhenClientRetunNoError(t *testing.T) {
	c, mks := getAgentServiceMocks()

	mks.client.On("ReportMetricsJSONBatch", mock.Anything, mock.Anything).Return(nil)
	err := c.ReportMetrics(context.Background())

	assert.NoError(t, err)
}

type agentServiceMocks struct {
	client *mocks.AgentAPIClient
}

func getAgentServiceMocks() (AgentService, *agentServiceMocks) {
	mks := &agentServiceMocks{
		client: &mocks.AgentAPIClient{},
	}

	as := New(mks.client, config.CollectableMetrics{})
	return as, mks
}
