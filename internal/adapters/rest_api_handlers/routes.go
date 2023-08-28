package handlers

import (
	"github.com/go-chi/chi/v5"
)

func RegisterHandlers(router *chi.Mux, h MetricsHandlers) {
	router.Route("/update", func(r chi.Router) {
		r.Post("/", h.UpdateMetricJSON)
		r.Post("/*", h.UpdateMetric)
	})
	router.Route("/value", func(r chi.Router) {
		r.Post("/", h.GetMetricJSON)
		r.Post("/*", h.UpdateMetric)
	})
	router.Get("/", h.AllMetrics)
}
