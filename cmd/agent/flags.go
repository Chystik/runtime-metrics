package main

import (
	"flag"
	"os"

	"github.com/Chystik/runtime-metrics/config"
)

const (
	confFileEnv = "CONFIG"
)

func parseFlags(cfg *config.AgentConfig) error {
	// checking interface implementation
	_ = flag.Value(cfg)
	_ = flag.Value(&cfg.PollInterval)
	_ = flag.Value(&cfg.ReportInterval)

	var configFileShort, conigFile string

	flag.StringVar(&configFileShort, "c", "", "path to config file")
	flag.StringVar(&conigFile, "config", "", "path to config file")

	flag.Var(cfg, "a", "Net address host:port")
	flag.StringVar((*string)(&cfg.TransportType), "t", "http", "Transport type: grpc or http")
	flag.Var(&cfg.PollInterval, "p", "Poll interval in seconds, min 0.000000001 sec")
	flag.Var(&cfg.ReportInterval, "r", "Report interval in seconds, min 0.000000001 sec")
	flag.StringVar(&cfg.SHAkey, "k", "", "sha key")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "public key (CRT) file path")
	flag.IntVar(&cfg.RateLimit, "l", 1, "report metrics rate limiter")
	flag.StringVar(&cfg.ProfileConfig.CPUFilePath, "cpu", "", "pprof CPU out profile")
	flag.StringVar(&cfg.ProfileConfig.MemFilePath, "mem", "", "pprof Memory out profile")

	flag.Parse()

	if configFileShort != "" {
		return config.ParseFile(cfg, configFileShort)
	} else if conigFile != "" {
		return config.ParseFile(cfg, conigFile)
	} else if ce := os.Getenv(confFileEnv); ce != "" {
		return config.ParseFile(cfg, ce)
	}

	return nil
}
