package handlers

import (
	"testing"
	"time"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/service"
	"github.com/Chystik/runtime-metrics/internal/service/mocks"
	"github.com/go-chi/chi/v5"
)

func TestNewRouter(t *testing.T) {
	type args struct {
		cfg         *config.ServerConfig
		router      *chi.Mux
		ms          service.MetricsService
		db          service.DBClient
		pingTimeout time.Duration
		logger      service.AppLogger
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Init router",
			args: args{
				cfg: &config.ServerConfig{
					SHAkey: "some key",
				},
				router:      chi.NewRouter(),
				ms:          &mocks.MetricsService{},
				db:          &mocks.DBClient{},
				pingTimeout: time.Second,
				logger:      &mocks.AppLogger{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewRouter(tt.args.cfg, tt.args.router, tt.args.ms, tt.args.db, tt.args.pingTimeout, tt.args.logger)
		})
	}
}
