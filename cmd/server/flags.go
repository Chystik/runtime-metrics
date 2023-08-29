package main

import (
	"flag"
	"time"

	"github.com/Chystik/runtime-metrics/config"
)

func parseFlags(cfg *config.ServerConfig) {
	// checking interface implementation
	_ = flag.Value(cfg)

	flag.StringVar(&cfg.LogLevel, "l", "info", "log levels")
	flag.DurationVar(&cfg.StoreInterval, "i", 300*time.Second, "interval for saving data to a file, in seconds. 0 value means synchronous data writing")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "file storage path")
	flag.BoolVar(&cfg.Restore, "r", true, "restore data from file on startup")
	flag.Var(cfg, "a", "Net address host:port")
	flag.Parse()
}
