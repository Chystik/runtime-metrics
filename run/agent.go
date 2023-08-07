package run

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/adapters"
	agentservice "github.com/Chystik/runtime-metrics/internal/service/agent"
	httpclient "github.com/Chystik/runtime-metrics/internal/transport/agent"
)

func Agent(cfg *config.AgentConfig) {
	client := httpclient.NewHTTPClient(cfg.HTTPServer)
	agentClient := adapters.NewAgentClient(client, cfg.HTTPServer)
	agentService := agentservice.New(agentClient, cfg.CollectableMetrics)

	updateTicker := time.NewTicker(time.Duration(cfg.PollInterval))
	reportTicker := time.NewTicker(time.Duration(cfg.ReportInterval))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case <-updateTicker.C:
			agentService.UpdateMetrics()
		case <-reportTicker.C:
			agentService.ReportMetrics()
		case <-quit:
			fmt.Println("Interrupt signal. Shutdown")
			os.Exit(0)
		}
	}
}
