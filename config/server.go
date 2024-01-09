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
		Address         string        `env:"ADDRESS" json:"address"`
		AddressGRPC     string        `env:"ADDRESS_GRPC" json:"address_grpc"`
		LogLevel        string        `env:"LOG_LEVEL"`
		StoreInterval   StoreInterval `json:"store_interval"`
		FileStoragePath string        `env:"FILE_STORAGE_PATH" json:"store_file"`
		Restore         bool          `env:"RESTORE" json:"restore"`
		DBDsn           string        `env:"DATABASE_DSN" json:"database_dsn"`
		SHAkey          string        `env:"KEY"`
		CryptoKey       string        `env:"CRYPTO_KEY" json:"crypto_key"`
		TrustedSubnet   string        `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
		ProfileConfig   ProfileConfig
	}

	StoreInterval struct {
		time.Duration `env:"STORE_INTERVAL"`
	}
)

func NewServerCfg() *ServerConfig {
	cfg := &ServerConfig{
		Address:         ":8080",
		AddressGRPC:     ":8081",
		LogLevel:        "info",
		StoreInterval:   StoreInterval{Duration: 300 * time.Second},
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

func (si StoreInterval) String() string {
	return fmt.Sprintf("%d", si)
}

func (si *StoreInterval) Set(s string) error {
	t, err := strconv.Atoi(s)
	if err != nil {
		return errors.New("only digits allowed for store intervar")
	}
	*si = StoreInterval{Duration: time.Duration(t) * time.Second}
	return nil
}

func (si *StoreInterval) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`) // remove quotes
	t, err := time.ParseDuration(s)
	si.Duration = t
	return
}
