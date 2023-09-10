package run

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/adapters/db"
	handlers "github.com/Chystik/runtime-metrics/internal/adapters/rest_api_handlers"
	"github.com/Chystik/runtime-metrics/internal/compressor"
	memstorage "github.com/Chystik/runtime-metrics/internal/infrastructure/repository/mem_storage"
	localfs "github.com/Chystik/runtime-metrics/internal/infrastructure/storage/local"
	"github.com/Chystik/runtime-metrics/internal/logger"
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
	logDbDisconnect                = "Graceful close connection for DB client complete."
)

func Server(cfg *config.ServerConfig, quit chan os.Signal) {
	// logger
	logger, err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		logger.Fatal(err.Error())
	}

	// repository
	meticsRepository := memstorage.New(cfg)

	localStorage, err := localfs.New(cfg, meticsRepository)
	if err != nil {
		logger.Fatal(err.Error())
	}

	repoWithSyncer, err := syncer.New(cfg, meticsRepository, localStorage)
	if err != nil {
		logger.Fatal(err.Error())
	}

	// postgres client
	pgClient, err := db.NewPgClient(cfg)
	if err != nil {
		logger.Fatal(err.Error())
	}

	// services
	metricsService := metricsservice.New(repoWithSyncer)

	// router
	router := chi.NewRouter()
	router.Use(logger.WithLogging)
	router.Use(compressor.GzipMiddleware)
	router.Use(middleware.Recoverer)

	// handlers
	metricHandlers := handlers.NewMetricsHandlers(metricsService)
	handlers.RegisterHandlers(router, metricHandlers, pgClient)

	// periodically or permanently writes repo data to a file
	go func() {
		logger.Info(fmt.Sprintf(logStorageSyncStart, cfg.FileStoragePath, time.Duration(cfg.StoreInterval)))
		if err := repoWithSyncer.SyncData(); err != nil {
			logger.Fatal(err.Error())
		}
		logger.Info(logStorageSyncStop)
	}()

	// http server
	server := restapi.NewServer(cfg, router)
	go func() {
		logger.Info(fmt.Sprintf(logHTTPServerStart, cfg.Address))
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err.Error())
		}
		logger.Info(logHTTPServerStop)
	}()

	<-quit
	logger.Info(logSignalInterrupt)
	ctxShutdown, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	// Graceful shutdown syncer
	if err := repoWithSyncer.Shutdown(ctxShutdown); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info(logGracefulStorageSyncShutdown)

	// Graceful shutdown HTTP Server
	if err := server.Shutdown(ctxShutdown); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info(logGracefulHTTPServerShutdown)

	// Graceful disconnect db client
	if err := pgClient.Disconnect(ctxShutdown); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info(logDbDisconnect)
}
