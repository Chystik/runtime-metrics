package run

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/adapters"
	memstorage "github.com/Chystik/runtime-metrics/internal/infrastructure/repository/mem_storage"
	metricsservice "github.com/Chystik/runtime-metrics/internal/service/server"
	"github.com/Chystik/runtime-metrics/internal/transport/restapi"
)

const (
	logHTTPServerStart            = "HTTP server started on: %s\n"
	logHTTPServerStop             = "Stopped serving new connections"
	logSignalInterrupt            = "Interrupt signal. Shutdown"
	logGracefulHTTPServerShutdown = "Graceful shutdown of HTTP Server complete."
)

func Server(cfg *config.ServerConfig, quit chan os.Signal) {
	// repository
	metricsRepository := memstorage.New()

	// services
	metricsService := metricsservice.New(metricsRepository)

	// handlers
	serverHandlers := adapters.NewServerHandlers(metricsService)

	// http server
	server := restapi.NewServer(cfg, serverHandlers)
	go func() {
		fmt.Printf(logHTTPServerStart, cfg.Address)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
		fmt.Println(logHTTPServerStop)
	}()

	<-quit
	fmt.Println(logSignalInterrupt)
	ctxShutdown, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	// Graceful shutdown HTTP Server
	if err := server.Shutdown(ctxShutdown); err != nil {
		panic(err)
	}
	fmt.Println(logGracefulHTTPServerShutdown)
}
