package main

import (
	"os"
	"os/signal"
	"syscall"

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

	// Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	run.Agent(cfg, quit)
}
