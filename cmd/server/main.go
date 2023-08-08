package main

import (
	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/run"
)

func main() {
	cfg, err := config.NewServerCfg()
	if err != nil {
		panic(err)
	}
	parseFlags(cfg)
	run.Server(cfg)
}
