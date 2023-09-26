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

var (
	errConnNextTry = "cannot connect to the server. next try in %d seconds\n"
	errConnExit    = "cannot connect to the server: %v. exit"
)

type jobResult struct {
	err error
}

func Agent(cfg *config.AgentConfig, quit chan os.Signal) {
	client := httpclient.NewHTTPClient(cfg)
	agentClient := agenthttpclient.New(client, cfg)
	agentService := agentservice.New(agentClient, cfg.CollectableMetrics)

	p, r := time.Duration(cfg.PollInterval), time.Duration(cfg.ReportInterval)

	attempts := 3

	// intervals in seconds
	triesInterval := 1
	deltaInterval := 2

	updateTicker := time.NewTicker(p)
	reportTicker := time.NewTicker(r)

	numJobs := cfg.RateLimit

	jobs := make(chan struct{}, numJobs)
	results := make(chan jobResult, numJobs)

	// init and run N workers, where N = RATE_LIMIT
	for w := 0; w < numJobs; w++ {
		go worker(agentService, attempts, triesInterval, deltaInterval, jobs, results)
	}

	go func() {
		for {
			select {
			case <-updateTicker.C:
				go agentService.UpdateMetrics()
				go agentService.UpdateGoPsUtilMetrics()
			case <-reportTicker.C:
				jobs <- struct{}{}
			case <-quit:
				fmt.Println("Interrupt signal. Shutdown")
				close(jobs)
				close(results)
				os.Exit(0)
			}
		}
	}()

	for r := range results {
		if r.err != nil {
			fmt.Println(r.err)
			os.Exit(0)
		}
	}
}

func worker(as agentservice.AgentService, attempts, triesInterval, deltaInterval int, jobs chan struct{}, results chan jobResult) {
	for range jobs {
		err := reportMetricsRetryer(as, attempts, triesInterval, deltaInterval)
		results <- jobResult{err: err}
	}
}

func reportMetricsRetryer(as agentservice.AgentService, attempts, triesInterval, deltaInterval int) error {
	reportCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := as.ReportMetrics(reportCtx)
	if err != nil {
		var netOpErr *net.OpError
		if errors.As(err, &netOpErr) && errors.Is(netOpErr, syscall.ECONNREFUSED) {
			a := attempts
			ti := triesInterval
			for a > 0 && err != nil {
				fmt.Printf(errConnNextTry, ti)
				time.Sleep(time.Duration(ti) * time.Second)
				err = as.ReportMetrics(reportCtx)
				a--
				ti = ti + deltaInterval
			}
			if err != nil {
				return fmt.Errorf(errConnExit, err)
			}
		}
	}
	return err
}
