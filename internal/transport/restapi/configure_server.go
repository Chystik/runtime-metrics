package restapi

import (
	"net/http"

	"github.com/Chystik/runtime-metrics/config"
	handlers "github.com/Chystik/runtime-metrics/internal/adapters/rest_api_handlers"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func NewServer(cfg *config.ServerConfig, handlers handlers.MetricsHandlers) *http.Server {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	registerHandlers(router, handlers)

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	return server
}
