package run

import (
	"fmt"
	"os"
	"time"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/adapters"
	agentservice "github.com/Chystik/runtime-metrics/internal/service/agent"
	httpclient "github.com/Chystik/runtime-metrics/internal/transport/agent"
)

func Agent(cfg *config.AgentConfig, quit chan os.Signal) {
	client := httpclient.NewHTTPClient(cfg)
	agentClient := adapters.NewAgentClient(client, cfg)
	agentService := agentservice.New(agentClient, cfg.CollectableMetrics)

	updateTicker := time.NewTicker(time.Duration(cfg.PollInterval))
	reportTicker := time.NewTicker(time.Duration(cfg.ReportInterval))

	// waiting for the server to start
	time.Sleep(2 * time.Second)

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
