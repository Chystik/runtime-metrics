package main

import (
	"flag"

	"github.com/Chystik/runtime-metrics/config"
)

func parseFlags(cfg *config.ServerConfig) {
	// checking interface implementation
	_ = flag.Value(cfg)

	flag.StringVar(&cfg.LogLevel, "l", "info", "log levels")
	flag.Var(cfg, "a", "Net address host:port")
	flag.Parse()
}
