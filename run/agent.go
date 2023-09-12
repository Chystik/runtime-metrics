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
	triesCount, triesInterval := 3, 1

	updateTimer := time.NewTimer(p)
	reportTimer := time.NewTimer(r)

	for {
		select {
		case <-updateTimer.C:
			agentService.UpdateMetrics()
			updateTimer.Reset(p)
		case <-reportTimer.C:
			reportCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			err := agentService.ReportMetrics(reportCtx)
			if err != nil {
				if triesCount == 0 {
					fmt.Println("cannot connect to the server. exit")
					os.Exit(0)
				}
				var netOpErr *net.OpError
				if errors.As(err, &netOpErr) && errors.Is(netOpErr, syscall.ECONNREFUSED) {
					reportTimer.Reset(time.Duration(triesInterval * int(time.Second)))
					fmt.Printf("cannot connect to the server. next try in %d seconds\n", triesInterval)
					triesCount--
					triesInterval = triesInterval + 2
					cancel()
					continue
				}
				fmt.Println(err)
			}
			reportTimer.Reset(r)
			triesCount, triesInterval = 3, 1
			cancel()
		case <-quit:
			fmt.Println("Interrupt signal. Shutdown")
			os.Exit(0)
		}
	}
}
