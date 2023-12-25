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
		Address        string         `env:"ADDRESS" json:"address"`
		PollInterval   PollInterval   `json:"poll_interval"`
		ReportInterval ReportInterval `json:"report_interval"`
		SHAkey         string         `env:"KEY"`
		CryptoKey      string         `env:"CRYPTO_KEY" json:"crypto_key"`
		RateLimit      int            `env:"RATE_LIMIT"`
		CollectableMetrics
		ProfileConfig ProfileConfig
	}

	PollInterval struct {
		time.Duration `env:"POLL_INTERVAL" `
	}
	ReportInterval struct {
		time.Duration `env:"REPORT_INTERVAL"`
	}

	CollectableMetrics []string
)

func NewAgentCfg() *AgentConfig {
	cfg := &AgentConfig{
		Address:            ":8080",
		PollInterval:       PollInterval{Duration: 2 * time.Second},
		ReportInterval:     ReportInterval{Duration: 10 * time.Second},
		CollectableMetrics: []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc"},
		ProfileConfig:      ProfileConfig{},
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
	t, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return errors.New("only digits allowed for Poll intervar")
	}
	*pi = PollInterval{Duration: time.Duration(t * 1e9)} // second
	return nil
}

func (pi *PollInterval) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`) // remove quotes
	t, err := time.ParseDuration(s)
	pi.Duration = t
	return
}

func (ri ReportInterval) String() string {
	return fmt.Sprintf("%d", ri)
}

func (ri *ReportInterval) Set(s string) error {
	t, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return errors.New("only digits allowed for Report intervar")
	}
	*ri = ReportInterval{Duration: time.Duration(t * 1e9)} // second
	return nil
}

func (ri *ReportInterval) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`) // remove quotes
	t, err := time.ParseDuration(s)
	ri.Duration = t
	return
}
