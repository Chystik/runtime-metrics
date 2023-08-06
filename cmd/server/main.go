package main

import (
	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/run"
)

func main() {
	cfg := config.NewServerCfg()
	run.Server(cfg)
}
