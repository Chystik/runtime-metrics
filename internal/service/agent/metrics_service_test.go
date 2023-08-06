package agentservice

import (
	"runtime"
	"testing"

	"github.com/Chystik/runtime-metrics/internal/adapters"
	"github.com/Chystik/runtime-metrics/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	var c adapters.AgentClient
	var m []string

	agentService := New(c, m)

	assert.NotNil(t, agentService)
}

func Test_agentService_UpdateMetrics(t *testing.T) {
	cache := make(map[string]interface{})
	cache["PollCount"] = models.Counter(0)

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
