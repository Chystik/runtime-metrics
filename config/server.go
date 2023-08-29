package config

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type (
	ServerConfig struct {
		Address         string        `env:"ADDRESS"`
		LogLevel        string        `env:"LOG_LEVEL"`
		StoreInterval   time.Duration `env:"STORE_INTERVAL"`
		FileStoragePath string        `env:"FILE_STORAGE_PATH"`
		Restore         bool          `env:"RESTORE"`
	}
)

func NewServerCfg() *ServerConfig {
	cfg := &ServerConfig{
		Address:         ":8080",
		LogLevel:        "info",
		StoreInterval:   300 * time.Second,
		FileStoragePath: "/tmp/metrics-db.json",
		Restore:         true,
	}

	return cfg
}

func (cfg ServerConfig) String() string {
	return cfg.Address
}

func (cfg *ServerConfig) Set(s string) error {
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
