package handlers

import (
	"github.com/Chystik/runtime-metrics/internal/adapters"

	"github.com/go-chi/chi/v5"
)

func RegisterHandlers(router *chi.Mux, h adapters.MetricsHandlers, pg adapters.PgClient) {
	router.Route("/update/", func(r chi.Router) {
		r.Post("/", h.UpdateMetricJSON)
		r.Post("/*", h.UpdateMetric)
	})
	router.Route("/value/", func(r chi.Router) {
		r.Post("/", h.GetMetricJSON)
		r.Post("/*", h.GetMetric)
	})
	router.Get("/value/*", h.GetMetric)
	if pg != nil {
		router.Get("/ping", pg.PingHandler)
	}
	router.Get("/", h.AllMetrics)
	router.Post("/updates/", h.UpdateMetricsJSON)
}
