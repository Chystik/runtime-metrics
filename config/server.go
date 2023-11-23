package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	ServerConfig struct {
		Address         string `env:"ADDRESS"`
		LogLevel        string `env:"LOG_LEVEL"`
		StoreInterval   `env:"STORE_INTERVAL"`
		FileStoragePath string `env:"FILE_STORAGE_PATH"`
		Restore         bool   `env:"RESTORE"`
		DBDsn           string `env:"DATABASE_DSN"`
		SHAkey          string `env:"KEY"`
		ProfileConfig   ProfileConfig
	}

	StoreInterval time.Duration
)

func NewServerCfg() *ServerConfig {
	cfg := &ServerConfig{
		Address:         ":8080",
		LogLevel:        "info",
		StoreInterval:   StoreInterval(300 * time.Second),
		FileStoragePath: "/tmp/metrics-db.json",
		Restore:         true,
		ProfileConfig:   ProfileConfig{},
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

func (cfg StoreInterval) String() string {
	return fmt.Sprintf("%d", cfg)
}

func (cfg *StoreInterval) Set(s string) error {
	t, err := strconv.Atoi(s)
	if err != nil {
		return errors.New("only digits allowed for store intervar")
	}
	*cfg = StoreInterval(time.Duration(t) * time.Second)
	return nil
}
