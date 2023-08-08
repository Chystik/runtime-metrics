package run

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/adapters"
	memstorage "github.com/Chystik/runtime-metrics/internal/infrastructure/repository/mem_storage"
	metricsservice "github.com/Chystik/runtime-metrics/internal/service/server"
	"github.com/Chystik/runtime-metrics/internal/transport/restapi"
)

func Server(cfg *config.ServerConfig) {
	// repository
	metricsRepository := memstorage.New()

	// services
	metricsService := metricsservice.New(metricsRepository)

	// handlers
	serverHandlers := adapters.NewServerHandlers(metricsService)

	// server
	server := restapi.NewServer(cfg, serverHandlers)
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
		fmt.Println("stopped serving new connections")
	}()

	// Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	fmt.Println("Interrupt signal. Shutdown")
	ctxShutdown, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	// Graceful shutdown HTTP Server
	if err := server.Shutdown(ctxShutdown); err != nil {
		panic(err)
	}
	fmt.Println("Graceful shutdown of HTTP Server complete.")
}
