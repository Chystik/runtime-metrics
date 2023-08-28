package run

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Chystik/runtime-metrics/config"
	handlers "github.com/Chystik/runtime-metrics/internal/adapters/rest_api_handlers"
	memstorage "github.com/Chystik/runtime-metrics/internal/infrastructure/repository/mem_storage"
	"github.com/Chystik/runtime-metrics/internal/logger"
	metricsservice "github.com/Chystik/runtime-metrics/internal/service/server"
	"github.com/Chystik/runtime-metrics/internal/transport/restapi"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const (
	logHTTPServerStart            = "HTTP server started on port: %s"
	logHTTPServerStop             = "Stopped serving new connections"
	logSignalInterrupt            = "Interrupt signal. Shutdown"
	logGracefulHTTPServerShutdown = "Graceful shutdown of HTTP Server complete."
)

func Server(cfg *config.ServerConfig, quit chan os.Signal) {
	// logger
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		logger.Log.Fatal(err.Error())
	}

	// repository
	metricsRepository := memstorage.New()

	// services
	metricsService := metricsservice.New(metricsRepository)

	// router
	router := chi.NewRouter()
	//router.Use(logger.WithLogging)
	router.Use(middleware.Recoverer)

	// handlers
	metricHandlers := handlers.NewMetricsHandlers(metricsService)
	handlers.RegisterHandlers(router, metricHandlers)

	// http server
	server := restapi.NewServer(cfg, router)
	go func() {
		logger.Log.Info(fmt.Sprintf(logHTTPServerStart, cfg.Address))
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal(err.Error())
		}
		logger.Log.Info(logHTTPServerStop)
	}()

	<-quit
	logger.Log.Info(logSignalInterrupt)
	ctxShutdown, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	// Graceful shutdown HTTP Server
	if err := server.Shutdown(ctxShutdown); err != nil {
		logger.Log.Fatal(err.Error())
	}
	logger.Log.Info(logGracefulHTTPServerShutdown)
}
