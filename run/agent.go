package run

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/Chystik/runtime-metrics/config"
	agenthttpclient "github.com/Chystik/runtime-metrics/internal/adapters/http_client"
	agentservice "github.com/Chystik/runtime-metrics/internal/service/agent"
	"github.com/Chystik/runtime-metrics/internal/transport/httpclient"
)

func Agent(cfg *config.AgentConfig, quit chan os.Signal) {
	client := httpclient.NewHTTPClient(cfg)
	agentClient := agenthttpclient.New(client, cfg)
	agentService := agentservice.New(agentClient, cfg.CollectableMetrics)

	p, r := time.Duration(cfg.PollInterval), time.Duration(cfg.ReportInterval)

	attempts := 3

	// intervals in seconds
	triesInterval := 1
	deltaInterval := 2

	updateTimer := time.NewTimer(p)
	reportTimer := time.NewTimer(r)

	for {
		select {
		case <-updateTimer.C:
			agentService.UpdateMetrics()
			updateTimer.Reset(p)
		case <-reportTimer.C:
			reportMetricsRetrier(agentService, attempts, triesInterval, deltaInterval)
			reportTimer.Reset(r)
		case <-quit:
			fmt.Println("Interrupt signal. Shutdown")
			os.Exit(0)
		}
	}
}

func reportMetricsRetrier(as agentservice.AgentService, attempts, triesInterval, deltaInterval int) {
	reportCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := as.ReportMetrics(reportCtx)
	if err != nil {
		if attempts <= 0 {
			fmt.Println("cannot connect to the server. exit")
			os.Exit(0)
		}
		var netOpErr *net.OpError
		if errors.As(err, &netOpErr) && errors.Is(netOpErr, syscall.ECONNREFUSED) {
			fmt.Printf("cannot connect to the server. next try in %d seconds\n", triesInterval)
			time.Sleep(time.Duration(triesInterval) * time.Second)
			reportMetricsRetrier(as, attempts-1, triesInterval+deltaInterval, deltaInterval)
			return
		}
		fmt.Println(err)
	}
}
