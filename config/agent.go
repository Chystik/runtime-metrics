package config

import (
	"time"
)

type (
	AgentConfig struct {
		HTTPServer
		CollectableMetrics
		PollInterval   time.Duration
		ReportInterval time.Duration
	}

	HTTPServer struct {
		Host string
		Port uint
	}

	CollectableMetrics []string
)

func NewAgentCfg() *AgentConfig {
	cfg := &AgentConfig{
		HTTPServer: HTTPServer{
			Host: "localhost",
			Port: 8080,
		},
		CollectableMetrics: []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc"},
		PollInterval:       2 * time.Second,
		ReportInterval:     10 * time.Second,
	}
	return cfg
}
