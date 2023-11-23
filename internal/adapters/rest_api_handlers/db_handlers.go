package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Chystik/runtime-metrics/internal/service"
)

type dbHandlers struct {
	db          service.DBClient
	logger      service.AppLogger
	pingTimeout time.Duration
}

func NewDBHandlers(db service.DBClient, logger service.AppLogger, pingTimeout time.Duration) *dbHandlers {
	return &dbHandlers{
		db:          db,
		logger:      logger,
		pingTimeout: pingTimeout,
	}
}

func (dh *dbHandlers) PingDB(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), dh.pingTimeout)
	defer cancel()

	w.Header().Set("Content-Type", "text/plain")

	err := dh.db.Ping(ctx)
	if err != nil {
		dh.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
