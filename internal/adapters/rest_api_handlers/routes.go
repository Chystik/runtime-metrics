package handlers

import (
	"github.com/go-chi/chi/v5"
)

func RegisterHandlers(router *chi.Mux, h MetricsHandlers) {
	router.Route("/update", func(r chi.Router) {
		r.Post("/", h.UpdateMetric)
		r.Post("/*", h.UpdateMetric)
	})
	/* router.Post("/update", h.UpdateMetric)
	router.Post("/update/*", h.UpdateMetric) */
	router.Get("/value/*", h.GetMetric)
	router.Get("/", h.AllMetrics)
}
