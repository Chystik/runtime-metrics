package handlers

import (
	"github.com/go-chi/chi/v5"
)

func RegisterHandlers(router *chi.Mux, h MetricsHandlers) {
	router.Post("/update/*", h.UpdateMetric)
	router.Get("/value/*", h.GetMetric)
	router.Get("/", h.AllMetrics)
}
