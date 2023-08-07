package restapi

import (
	"fmt"
	"net/http"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/adapters"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func NewServer(cfg config.HTTP, handlers adapters.ServerHandlers) *http.Server {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	registerHandlers(router, handlers)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: router,
	}

	return server
}
