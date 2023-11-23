package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/run"
)

func main() {
	cfg := config.NewServerCfg()

	parseFlags(cfg)

	err := parseEnv(cfg)
	if err != nil {
		panic(err)
	}

	// Graceful shutdown setup
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	var wg sync.WaitGroup

	if cfg.ProfileConfig.CPUFilePath != "" && cfg.ProfileConfig.MemFilePath != "" {
		prof, err := run.NewProfile(cfg.ProfileConfig)
		if err != nil {
			panic(err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err = prof.Run(ctx)
			if err != nil {
				panic(err)
			}
			stop()
		}()
	}

	run.Server(ctx, cfg)

	wg.Wait()
}
