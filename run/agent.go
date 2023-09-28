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
	"github.com/Chystik/runtime-metrics/internal/logger"
	agentservice "github.com/Chystik/runtime-metrics/internal/service/agent"
	"github.com/Chystik/runtime-metrics/internal/transport/httpclient"
	"go.uber.org/zap"
)

var (
	errConnNextTry = "cannot connect to the server. next try in %d seconds"
	errConnExit    = "cannot connect to the server: %v. exit"
)

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

	attempts := 3

	// intervals in seconds
	triesInterval := 1
	deltaInterval := 2

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
		go worker(agentService, attempts, triesInterval, deltaInterval, jobs, results, logger.Logger)
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

func worker(as agentservice.AgentService, attempts, triesInterval, deltaInterval int, jobs chan struct{}, results chan jobResult, l *zap.Logger) {
	for range jobs {
		err := reportMetricsRetryer(as, attempts, triesInterval, deltaInterval, l)
		results <- jobResult{err: err}
	}
}

func reportMetricsRetryer(as agentservice.AgentService, attempts, triesInterval, deltaInterval int, l *zap.Logger) error {
	reportCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := as.ReportMetrics(reportCtx)
	if err != nil {
		var netOpErr *net.OpError
		if errors.As(err, &netOpErr) && errors.Is(netOpErr, syscall.ECONNREFUSED) {
			a := attempts
			ti := triesInterval
			for a > 0 && err != nil {
				l.Info(fmt.Sprintf(errConnNextTry, ti))
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
