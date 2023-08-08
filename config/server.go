package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type (
	ServerConfig struct {
		HTTP
	}

	HTTP struct {
		Host string
		Port int
	}
)

func NewServerCfg() (*ServerConfig, error) {
	var err error

	cfg := &ServerConfig{
		HTTP: HTTP{
			Host: "localhost",
			Port: 8080,
		},
	}

	// loads ENV from file if ENVIRONMENT val is set and parses it values
	err = cfg.loadEnv("ENVIRONMENT")
	if err != nil {
		return nil, err
	}

	err = cfg.parseEnv()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (adr HTTP) String() string {
	return fmt.Sprintf("%s:%s", adr.Host, strconv.Itoa(adr.Port))
}

func (adr *HTTP) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("expect address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	adr.Host = hp[0]
	adr.Port = port
	return nil
}

// load loads ENV from .env file, defined in envVal ENV
func (cfg *ServerConfig) loadEnv(envVal string) error {
	var envFile string
	var err error

	if osEnv := os.Getenv(envVal); osEnv != "" {
		switch osEnv {
		case "dev":
			envFile = "server_dev.env"
		case "prod":
			envFile = "server_prod.env"
		case "stage":
			envFile = "server_stage.env"
		}

		// load ENV from *.env file
		err = godotenv.Load(envFile)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

// pasrse parses ENV
func (cfg *ServerConfig) parseEnv() error {
	return env.Parse(cfg)
}
