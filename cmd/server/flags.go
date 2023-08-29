package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/Chystik/runtime-metrics/config"
)

func parseFlags(cfg *config.ServerConfig) {
	// checking interface implementation
	_ = flag.Value(cfg)

	flag.StringVar(&cfg.LogLevel, "l", "info", "log levels")
	tmpStrInterval := new(string)
	flag.StringVar(tmpStrInterval, "i", "300", "interval for saving data to a file, in seconds. 0 value means synchronous data writing")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "file storage path")
	flag.BoolVar(&cfg.Restore, "r", true, "restore data from file on startup")
	flag.Var(cfg, "a", "Net address host:port")
	flag.Parse()

	if *tmpStrInterval != "" {
		tmpInterval, err := strconv.Atoi(*tmpStrInterval)
		if err != nil {
			panic(err)
		}
		cfg.StoreInterval = time.Duration(tmpInterval * int(time.Second))
	}
	fmt.Println(cfg.StoreInterval, cfg.FileStoragePath)
}
