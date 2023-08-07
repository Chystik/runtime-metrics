package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

func NewServerCfg() *ServerConfig {
	cfg := &ServerConfig{
		HTTP: HTTP{
			Host: "localhost",
			Port: 8080,
		},
	}
	return cfg
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
