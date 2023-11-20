package run

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/adapters"
	"github.com/Chystik/runtime-metrics/internal/adapters/db"
	handlers "github.com/Chystik/runtime-metrics/internal/adapters/rest_api_handlers"
	"github.com/Chystik/runtime-metrics/internal/compressor"
	"github.com/Chystik/runtime-metrics/internal/hasher"
	memstorage "github.com/Chystik/runtime-metrics/internal/infrastructure/repository/mem_storage"
	"github.com/Chystik/runtime-metrics/internal/infrastructure/repository/postgres"
	localfs "github.com/Chystik/runtime-metrics/internal/infrastructure/storage/local"
	"github.com/Chystik/runtime-metrics/internal/logger"
	"github.com/Chystik/runtime-metrics/internal/retryer"
	metricsservice "github.com/Chystik/runtime-metrics/internal/service/server"
	"github.com/Chystik/runtime-metrics/internal/syncer"
	"github.com/Chystik/runtime-metrics/internal/transport/restapi"

	"github.com/go-chi/chi/middleware"
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

func Server(ctx context.Context, cfg *config.ServerConfig) {
	// logger
	logger, err := logger.Initialize(cfg.LogLevel, "./server.log")
	if err != nil {
		logger.Fatal(err.Error())
	}

	// repository
	var meticsRepository metricsservice.MetricsRepository
	var pgClient adapters.PgClient
	repoWithSyncer := syncer.New(cfg)

	inMemRepo := memstorage.New(cfg)

	if cfg.DBDsn != "" {
		// postgres
		pgClient, err = db.NewPgClient(cfg, logger.Logger)
		if err != nil {
			logger.Fatal(err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		db, err := pgClient.Connect(ctx)
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
			logger.Logger,
		)

		meticsRepository = postgres.NewMetricsRepo(db, r, logger.Logger)
	} else if cfg.FileStoragePath != "" {
		// fs storage
		localStorage, err := localfs.New(cfg, inMemRepo)
		if err != nil {
			logger.Fatal(err.Error())
		}
		err = repoWithSyncer.Initialize(cfg, inMemRepo, localStorage)
		if err != nil {
			logger.Fatal(err.Error())
		}
		// periodically or permanently writes repo data to a file
		go func() {
			logger.Info(fmt.Sprintf(logStorageSyncStart, cfg.FileStoragePath, time.Duration(cfg.StoreInterval)))
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

	// hasher
	h := hasher.NewHasher(cfg.SHAkey, "HashSHA256")

	// router
	router := chi.NewRouter()
	router.Use(logger.WithLogging)
	router.Use(h.WithHasher)
	router.Use(compressor.Compress(1))
	router.Use(middleware.Recoverer)

	// handlers
	metricHandlers := handlers.NewMetricsHandlers(metricsService)
	handlers.RegisterHandlers(router, metricHandlers, pgClient)

	// http server
	server := restapi.NewServer(cfg, router)
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
