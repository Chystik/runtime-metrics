package restapi

import (
	handlers "github.com/Chystik/runtime-metrics/internal/adapters/rest_api_handlers"

	"github.com/go-chi/chi/v5"
)

func registerHandlers(router *chi.Mux, handlers handlers.MetricsHandlers) {
	router.Post("/update/*", handlers.UpdateMetric)
	router.Get("/value/*", handlers.GetMetric)
	router.Get("/", handlers.AllMetrics)
}
