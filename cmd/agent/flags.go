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
	flag.Var(&cfg.PollInterval, "p", "Poll interval in seconds, min 0.000000001 sec")
	flag.Var(&cfg.ReportInterval, "r", "Report interval in seconds, min 0.000000001 sec")
	flag.StringVar(&cfg.SHAkey, "k", "", "sha key")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "public key (CRT) file path")
	flag.IntVar(&cfg.RateLimit, "l", 1, "report metrics rate limiter")
	flag.StringVar(&cfg.ProfileConfig.CPUFilePath, "cpu", "", "pprof CPU out profile")
	flag.StringVar(&cfg.ProfileConfig.MemFilePath, "mem", "", "pprof Memory out profile")

	flag.Parse()
}
