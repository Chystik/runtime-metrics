package main

import (
	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/run"
)

func main() {
	cfg := config.NewAgentCfg()

	parseFlags(cfg)

	err := parseEnv(cfg)
	if err != nil {
		panic(err)
	}

	run.Agent(cfg)
}
