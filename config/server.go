package config

import (
	"errors"
	"strconv"
	"strings"
)

type (
	ServerConfig struct {
		Address string `env:"ADDRESS"`
	}
)

func NewServerCfg() *ServerConfig {
	cfg := &ServerConfig{
		Address: ":8080",
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
