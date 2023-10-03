package run

import (
	"context"
	"fmt"
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

	jobs := make(chan int, 1)

	logger.Info(
		"agent started",
		zap.String("Address", cfg.Address),
		zap.Duration("Poll interval", time.Duration(cfg.PollInterval)),
		zap.Duration("Report interval", time.Duration(cfg.ReportInterval)),
		zap.Int("Rate limit", cfg.RateLimit),
	)

	// init and run N workers, where N = RATE_LIMIT
	for w := 0; w < numJobs; w++ {
		go worker(w, reportMetrics, jobs, logger.Logger)
	}

loop:
	for j := 0; ; {
		select {
		case <-updateTicker.C:
			go agentService.UpdateMetrics()
			go agentService.UpdateGoPsUtilMetrics()
		case <-reportTicker.C:
			if len(jobs) < cap(jobs) {
				jobs <- j
				j++
			}
		case <-quit:
			logger.Info("Interrupt signal. Shutdown")
			reportTicker.Stop()
			close(jobs)
			break loop
		}
	}
}

func worker(w int, fn connRetryerFn, jobs chan int, logger *zap.Logger) {
	for j := range jobs {
		logger.Debug(fmt.Sprintf("Worker %d started job %d", w, j))
		err := fn.DoWithRetryFn()
		if err != nil {
			logger.Error(err.Error())
		}
		logger.Debug(fmt.Sprintf("Worker %d finished job %d", w, j))
	}
}
