package main

import (
	"flag"

	"github.com/Chystik/runtime-metrics/config"
)

func parseFlags(cfg *config.AgentConfig) {
	// checking interface implementation
	_ = flag.Value(cfg)
	_ = flag.Value(&cfg.PollInterval)
	_ = flag.Value(&cfg.ReportInterval)

	flag.Var(cfg, "a", "Net address host:port")
	flag.Var(&cfg.PollInterval, "p", "Poll interval in seconds")
	flag.Var(&cfg.ReportInterval, "r", "Report interval in seconds")
	flag.StringVar(&cfg.SHAkey, "k", "", "sha key")
	flag.IntVar(&cfg.RateLimit, "l", 1, "report metrics rate limiter")
	flag.Parse()
}
