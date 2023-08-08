package main

import (
	"flag"

	"github.com/Chystik/runtime-metrics/config"
)

func parseFlags(cfg *config.ServerConfig) {
	// если интерфейс не реализован,
	// здесь будет ошибка компиляции
	_ = flag.Value(cfg)

	flag.Var(cfg, "a", "Net address host:port")
	flag.Parse()
}
