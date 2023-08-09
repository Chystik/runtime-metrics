package restapi

import (
	"net/http"

	"github.com/Chystik/runtime-metrics/config"
)

func NewServer(cfg *config.ServerConfig, handler http.Handler) *http.Server {
	server := &http.Server{
		Addr:    cfg.Address,
		Handler: handler,
	}

	return server
}
