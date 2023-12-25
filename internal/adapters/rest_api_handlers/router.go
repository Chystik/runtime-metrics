package handlers

import (
	"time"

	"github.com/Chystik/runtime-metrics/config"
	md "github.com/Chystik/runtime-metrics/internal/middleware"
	"github.com/Chystik/runtime-metrics/internal/service"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(
	cfg *config.ServerConfig,
	router *chi.Mux,
	ms service.MetricsService,
	db service.DBClient,
	pingTimeout time.Duration,
	logger service.AppLogger,
) error {
	// middleware
	router.Use(md.MidLogger(logger).WithLogging)
	if cfg.CryptoKey != "" {
		d, err := md.NewDecryptor(cfg.CryptoKey)
		if err != nil {
			return err
		}
		router.Use(d.WithDecryptor)
	}
	router.Use(md.NewHasher(cfg.SHAkey, "HashSHA256").WithHasher)
	router.Use(md.GzipPoolMiddleware())
	router.Use(middleware.Recoverer)

	// routes
	mh := NewMetricsHandlers(ms)

	router.Route("/update/", func(r chi.Router) {
		r.Post("/", mh.UpdateMetricJSON)
		r.Post("/*", mh.UpdateMetric)
	})
	router.Route("/value/", func(r chi.Router) {
		r.Post("/", mh.GetMetricJSON)
		r.Post("/*", mh.GetMetric)
	})
	router.Get("/value/*", mh.GetMetric)
	if db != nil {
		dh := NewDBHandlers(db, logger, pingTimeout)
		router.Get("/ping", dh.PingDB)
	}
	router.Get("/", mh.AllMetrics)
	router.Post("/updates/", mh.UpdateMetricsJSON)

	return nil
}
