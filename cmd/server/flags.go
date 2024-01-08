package main

import (
	"flag"
	"os"

	"github.com/Chystik/runtime-metrics/config"
)

const (
	confFileEnv = "CONFIG"
)

func parseFlags(cfg *config.ServerConfig) error {
	// checking interface implementation
	_ = flag.Value(cfg)

	var configFileShort, conigFile string

	flag.StringVar(&configFileShort, "c", "", "path to config file")
	flag.StringVar(&conigFile, "config", "", "path to config file")

	flag.StringVar(&cfg.LogLevel, "l", "info", "log levels")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "file storage path")
	flag.BoolVar(&cfg.Restore, "r", true, "restore data from file on startup")
	flag.StringVar(&cfg.DBDsn, "d", "", "postgres dsn")
	flag.StringVar(&cfg.SHAkey, "k", "", "sha key")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "private key (PEM) file path")
	flag.Var(&cfg.StoreInterval, "i", "interval for saving data to a file, in seconds. 0 value means synchronous data writing")
	flag.Var(cfg, "a", "Net address host:port")
	flag.StringVar(&cfg.TrustedSubnet, "t", "", "trusted subnet in CIDR format")
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
