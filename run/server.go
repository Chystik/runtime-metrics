package run

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Chystik/runtime-metrics/config"
	handlers "github.com/Chystik/runtime-metrics/internal/adapters/rest_api_handlers"

	"github.com/Chystik/runtime-metrics/internal/infrastructure/repository/inmemory"
	postgresrepo "github.com/Chystik/runtime-metrics/internal/infrastructure/repository/postgres"
	localfs "github.com/Chystik/runtime-metrics/internal/infrastructure/storage/local"
	"github.com/Chystik/runtime-metrics/internal/service"
	metricsservice "github.com/Chystik/runtime-metrics/internal/service/server"
	"github.com/Chystik/runtime-metrics/internal/syncer"
	"github.com/Chystik/runtime-metrics/pkg/httpserver"
	"github.com/Chystik/runtime-metrics/pkg/logger"
	"github.com/Chystik/runtime-metrics/pkg/postgres"
	"github.com/Chystik/runtime-metrics/pkg/retryer"

	"github.com/go-chi/chi/v5"
)

const (
	logHTTPServerStart             = "HTTP server started on port: %s"
	logHTTPServerStop              = "Stopped serving new connections"
	logSignalInterrupt             = "Interrupt signal. Shutdown"
	logGracefulHTTPServerShutdown  = "Graceful shutdown of HTTP Server complete."
	logStorageSyncStart            = "data syncronization to file %s with interval %v started"
	logStorageSyncStop             = "Stopped saving storage data to a file"
	logGracefulStorageSyncShutdown = "Graceful shutdown of storage sync complete."
	logDBDisconnect                = "Graceful close connection for DB client complete."
)

const (
	defaultDBPingTimeout = 3 * time.Second
)

func Server(ctx context.Context, cfg *config.ServerConfig) {
	// logger
	logger, err := logger.Initialize(cfg.LogLevel, "./server.log")
	if err != nil {
		logger.Fatal(err.Error())
	}

	// repository
	var meticsRepository service.MetricsRepository
	var pgClient *postgres.Postgres
	repoWithSyncer := syncer.New(cfg)

	inMemRepo := inmemory.NewMetricsRepo(cfg)

	if cfg.DBDsn != "" {
		// postgres
		pgClient, err = postgres.New(cfg.DBDsn, logger)
		if err != nil {
			logger.Fatal(err.Error())
		}

		ctxPing, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err = pgClient.Connect(ctxPing)
		if err != nil {
			logger.Fatal(err.Error())
		}

		err = pgClient.Migrate()
		if err != nil {
			logger.Fatal(err.Error())
		}

		// retryer
		r := retryer.NewConnRetryer(
			3,
			time.Duration(time.Second),
			time.Duration(2*time.Second),
			logger,
		)

		meticsRepository = postgresrepo.NewMetricsRepo(pgClient.DB, r, logger)
	} else if cfg.FileStoragePath != "" {
		// fs storage
		localStorage, err := localfs.NewMetricsStorage(cfg, inMemRepo)
		if err != nil {
			logger.Fatal(err.Error())
		}
		err = repoWithSyncer.Initialize(cfg, inMemRepo, localStorage)
		if err != nil {
			logger.Fatal(err.Error())
		}
		// periodically or permanently writes repo data to a file
		go func() {
			logger.Info(fmt.Sprintf(logStorageSyncStart, cfg.FileStoragePath, cfg.StoreInterval.Duration))
			if err := repoWithSyncer.SyncData(); err != nil {
				logger.Fatal(err.Error())
			}
			logger.Info(logStorageSyncStop)
		}()
		meticsRepository = repoWithSyncer
	} else {
		// inmemory storage
		meticsRepository = inMemRepo
	}

	// services
	metricsService := metricsservice.New(meticsRepository)

	// router
	handler := chi.NewRouter()
	handlers.NewRouter(
		cfg,
		handler,
		metricsService,
		pgClient,
		defaultDBPingTimeout,
		logger)

	// http server
	server := httpserver.NewServer(handler, httpserver.Address(cfg.Address))
	go func() {
		logger.Info(fmt.Sprintf(logHTTPServerStart, cfg.Address))
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err.Error())
		}
		logger.Info(logHTTPServerStop)
	}()

	// interrupt signal
	<-ctx.Done()

	logger.Info(logSignalInterrupt)
	ctxShutdown, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	// Graceful shutdown syncer
	if cfg.DBDsn == "" {
		if err := repoWithSyncer.Shutdown(ctxShutdown); err != nil {
			logger.Fatal(err.Error())
		}
		logger.Info(logGracefulStorageSyncShutdown)
	}

	// Graceful shutdown HTTP Server
	if err := server.Shutdown(ctxShutdown); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info(logGracefulHTTPServerShutdown)

	// Graceful disconnect db client
	if cfg.DBDsn != "" {
		if err := pgClient.Disconnect(ctxShutdown); err != nil {
			logger.Fatal(err.Error())
		}
		logger.Info(logDBDisconnect)
	}
}
