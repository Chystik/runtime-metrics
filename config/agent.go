package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	AgentConfig struct {
		Address        string `env:"ADDRESS"`
		PollInterval   `env:"POLL_INTERVAL"`
		ReportInterval `env:"REPORT_INTERVAL"`
		CollectableMetrics
	}

	PollInterval   time.Duration
	ReportInterval time.Duration

	CollectableMetrics []string
)

func NewAgentCfg() *AgentConfig {
	cfg := &AgentConfig{
		Address:            "localhost:8080",
		PollInterval:       PollInterval(2 * time.Second),
		ReportInterval:     ReportInterval(10 * time.Second),
		CollectableMetrics: []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc"},
	}

	return cfg
}

func (cfg AgentConfig) String() string {
	return cfg.Address
}

func (cfg *AgentConfig) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("expect address in a form host:port")
	}
	_, err := strconv.Atoi(hp[1])
	if err != nil {
		return errors.New("only digits allowed for port in a form host:port")
	}
	cfg.Address = s
	return nil
}

func (pi PollInterval) String() string {
	return fmt.Sprintf("%d", pi)
}

func (pi *PollInterval) Set(s string) error {
	t, err := strconv.Atoi(s)
	if err != nil {
		return errors.New("only digits allowed for Poll intervar")
	}
	*pi = PollInterval(time.Duration(t) * time.Second)
	return nil
}

func (ri ReportInterval) String() string {
	return fmt.Sprintf("%d", ri)
}

func (ri *ReportInterval) Set(s string) error {
	t, err := strconv.Atoi(s)
	if err != nil {
		return errors.New("only digits allowed for Report intervar")
	}
	*ri = ReportInterval(time.Duration(t) * time.Second)
	return nil
}
