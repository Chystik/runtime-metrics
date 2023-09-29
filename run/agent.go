package run

import (
	"context"
	"os"
	"time"

	"github.com/Chystik/runtime-metrics/config"
	agenthttpclient "github.com/Chystik/runtime-metrics/internal/adapters/http_client"
	"github.com/Chystik/runtime-metrics/internal/logger"
	"github.com/Chystik/runtime-metrics/internal/retryer"
	agentservice "github.com/Chystik/runtime-metrics/internal/service/agent"
	"github.com/Chystik/runtime-metrics/internal/transport/httpclient"
	"go.uber.org/zap"
)

type connRetryerFn interface {
	DoWithRetryFn() error
}

type jobResult struct {
	err error
}

func Agent(cfg *config.AgentConfig, quit chan os.Signal) {
	client := httpclient.NewHTTPClient(cfg)
	agentClient := agenthttpclient.New(client, cfg)
	agentService := agentservice.New(agentClient, cfg.CollectableMetrics)

	logger, err := logger.Initialize("info", "./agent.log")
	if err != nil {
		panic(err)
	}

	p, r := time.Duration(cfg.PollInterval), time.Duration(cfg.ReportInterval)

	reportMetrics := retryer.NewConnRetryerFn(
		3,
		time.Duration(time.Second),
		time.Duration(2*time.Second),
		logger.Logger,
		func() error {
			reportCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			return agentService.ReportMetrics(reportCtx)
		},
	)

	updateTicker := time.NewTicker(p)
	reportTicker := time.NewTicker(r)

	numJobs := cfg.RateLimit

	jobs := make(chan struct{}, numJobs)
	results := make(chan jobResult, numJobs)

	logger.Info(
		"agent started",
		zap.String("Address", cfg.Address),
		zap.Duration("Poll interval", time.Duration(cfg.PollInterval)),
		zap.Duration("Report interval", time.Duration(cfg.ReportInterval)),
		zap.Int("Rate limit", cfg.RateLimit),
	)

	// init and run N workers, where N = RATE_LIMIT
	for w := 0; w < numJobs; w++ {
		go worker(reportMetrics, jobs, results)
	}

	go func() {
		for {
			select {
			case <-updateTicker.C:
				go agentService.UpdateMetrics()
				go agentService.UpdateGoPsUtilMetrics()
			case <-reportTicker.C:
				go func() {
					jobs <- struct{}{}
				}()
			case <-quit:
				logger.Info("Interrupt signal. Shutdown")
				close(jobs)
				close(results)
				os.Exit(0)
			}
		}
	}()

	for r := range results {
		if r.err != nil {
			logger.Error(r.err.Error())
		}
	}
}

func worker(fn connRetryerFn, jobs chan struct{}, results chan jobResult) {
	for range jobs {
		err := fn.DoWithRetryFn()
		results <- jobResult{err: err}
	}
}
