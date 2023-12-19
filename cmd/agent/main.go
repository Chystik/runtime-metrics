package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/run"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	fmt.Println("Build version", "\t", buildVersion)
	fmt.Println("Build date", "\t", buildDate)
	fmt.Println("Build commit", "\t", buildCommit)

	cfg := config.NewAgentCfg()

	err := parseFlags(cfg)
	if err != nil {
		panic(err)
	}

	err = parseEnv(cfg)
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
		}()
	}

	run.Agent(ctx, cfg)

	wg.Wait()
}
