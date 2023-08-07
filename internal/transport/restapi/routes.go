package restapi

import (
	"github.com/Chystik/runtime-metrics/internal/adapters"

	"github.com/go-chi/chi/v5"
)

func registerHandlers(router *chi.Mux, handlers adapters.ServerHandlers) {
	router.Post("/update/*", handlers.UpdateMetric)
	router.Get("/value/*", handlers.GetMetric)
	router.Get("/", handlers.AllMetrics)
}
