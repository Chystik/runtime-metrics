package main

import (
	"flag"

	"github.com/Chystik/runtime-metrics/config"
)

func parseFlags(cfg *config.ServerConfig) {
	// checking interface implementation
	_ = flag.Value(cfg)

	flag.StringVar(&cfg.LogLevel, "l", "info", "log levels")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "file storage path")
	flag.BoolVar(&cfg.Restore, "r", true, "restore data from file on startup")
	flag.StringVar(&cfg.DBDsn, "d", "", "postgres dsn")
	flag.StringVar(&cfg.SHAkey, "k", "", "sha key")
	flag.Var(&cfg.StoreInterval, "i", "interval for saving data to a file, in seconds. 0 value means synchronous data writing")
	flag.Var(cfg, "a", "Net address host:port")
	flag.StringVar(&cfg.ProfileConfig.CPUFilePath, "cpu", "", "pprof CPU out profile")
	flag.StringVar(&cfg.ProfileConfig.MemFilePath, "mem", "", "pprof Memory out profile")

	flag.Parse()
}
