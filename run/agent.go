package run

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Chystik/runtime-metrics/config"
	agentapiclient "github.com/Chystik/runtime-metrics/internal/adapters/http_client"
	"github.com/Chystik/runtime-metrics/internal/service"
	agentservice "github.com/Chystik/runtime-metrics/internal/service/agent"
	"github.com/Chystik/runtime-metrics/pkg/httpclient"
	"github.com/Chystik/runtime-metrics/pkg/logger"
	"github.com/Chystik/runtime-metrics/pkg/retryer"

	"go.uber.org/zap"
)

const (
	httpClientTimeout    = 20 * time.Second
	reportMetricsTimeout = 10 * time.Second
	loggerLevel          = "info"
)

func Agent(ctx context.Context, cfg *config.AgentConfig) {
	logger, err := logger.Initialize(loggerLevel, "./agent.log")
	if err != nil {
		panic(err)
	}

	var client service.HTTPClient
	var pemKey []byte

	if cfg.CryptoKey != "" {
		pemKey, err = os.ReadFile(cfg.CryptoKey)
		if err != nil {
			logger.Fatal(err.Error())
		}
		client, err = httpclient.NewClient(
			httpclient.Timeout(httpClientTimeout),
			httpclient.WithEncryption(pemKey),
			httpclient.ExtractOutboundIP("X-Real-IP"),
		)
	} else {
		client, err = httpclient.NewClient(
			httpclient.Timeout(httpClientTimeout),
			httpclient.ExtractOutboundIP("X-Real-IP"),
		)
	}
	if err != nil {
		logger.Fatal(err.Error())
	}

	agentClient := agentapiclient.New(client, cfg)
	agentService := agentservice.New(agentClient, cfg.CollectableMetrics)

	p, r := cfg.PollInterval.Duration, cfg.ReportInterval.Duration

	reportMetrics := retryer.NewConnRetryerFn(
		3,
		time.Duration(time.Second),
		time.Duration(2*time.Second),
		logger,
		func() error {
			reportCtx, cancel := context.WithTimeout(context.Background(), reportMetricsTimeout)
			defer cancel()
			return agentService.ReportMetrics(reportCtx)
		},
	)

	updateTicker := time.NewTicker(p)
	reportTicker := time.NewTicker(r)

	numJobs := cfg.RateLimit

	jobs := make(chan struct{}, 1)

	logger.Info(
		"agent started",
		zap.String("Address", cfg.Address),
		zap.Duration("Poll interval", cfg.PollInterval.Duration),
		zap.Duration("Report interval", cfg.ReportInterval.Duration),
		zap.Int("Rate limit", cfg.RateLimit),
	)

	var wg sync.WaitGroup

	// init and run N workers, where N = RATE_LIMIT
	for w := 1; w < numJobs+1; w++ {
		wg.Add(1)
		go func(i int) {
			worker(i, reportMetrics, jobs, logger)
			wg.Done()
		}(w)
	}

loop:
	for {
		select {
		case <-updateTicker.C:
			go agentService.UpdateMetrics()
			go agentService.UpdateGoPsUtilMetrics()
		case <-reportTicker.C:
			if len(jobs) < cap(jobs) {
				jobs <- struct{}{}
			}
		case <-ctx.Done():
			logger.Info("Interrupt signal. Shutting down.")
			updateTicker.Stop()
			reportTicker.Stop()
			close(jobs)
			break loop
		}
	}

	logger.Info("Waiting for requests to be completed by all workers")
	wg.Wait()
}

func worker(w int, fn service.ConnectionRetrierFn, jobs chan struct{}, logger service.AppLogger) {
	for range jobs {
		logger.Debug(fmt.Sprintf("Worker %d started job", w))
		err := fn.DoWithRetryFn()
		if err != nil {
			logger.Error(err.Error())
		}
		logger.Debug(fmt.Sprintf("Worker %d finished job", w))
	}
}
