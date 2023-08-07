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
		HTTPServer
		CollectableMetrics
		PollInterval
		ReportInterval
	}

	PollInterval   time.Duration
	ReportInterval time.Duration

	HTTPServer struct {
		Host string
		Port int
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
		PollInterval:       PollInterval(2 * time.Second),
		ReportInterval:     ReportInterval(10 * time.Second),
	}
	return cfg
}

func (adr HTTPServer) String() string {
	return fmt.Sprintf("%s:%s", adr.Host, strconv.Itoa(adr.Port))
}

func (adr *HTTPServer) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("expect address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	adr.Host = hp[0]
	adr.Port = port
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
