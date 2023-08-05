package restapi

import (
	"fmt"
	"net/http"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/adapters"
)

func NewServer(cfg config.HTTP, handlers adapters.ServerHandlers) *http.Server {
	router := http.NewServeMux()

	registerHandlers(router, handlers)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: router,
	}

	return server
}
