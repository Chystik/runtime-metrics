package syncer

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/infrastructure/storage"
	"github.com/Chystik/runtime-metrics/internal/models"
	metricsservice "github.com/Chystik/runtime-metrics/internal/service/server"
)

var once sync.Once

type syncer struct {
	src  metricsservice.MetricsRepository
	dst  storage.MetricsStorage
	tick chan struct{}
	i    time.Duration
}

func New(cfg *config.ServerConfig) *syncer {
	return &syncer{
		tick: make(chan struct{}, 1),
		i:    time.Duration(cfg.StoreInterval),
	}
}

func (s *syncer) Initialize(cfg *config.ServerConfig, src metricsservice.MetricsRepository, dst storage.MetricsStorage) error {
	s.src = src
	s.dst = dst

	if cfg.Restore {
		err := s.dst.Read()
		if err != nil && err != io.EOF {
			return err
		}
	}
	return nil
}

func (s *syncer) UpdateGauge(ctx context.Context, metric models.Metric) error {
	err := s.src.UpdateGauge(ctx, metric)
	if err != nil {
		return err
	}

	s.sync()
	return nil
}
func (s *syncer) UpdateCounter(ctx context.Context, metric models.Metric) error {
	err := s.src.UpdateCounter(ctx, metric)
	if err != nil {
		return err
	}

	s.sync()
	return nil
}

func (s *syncer) Get(ctx context.Context, metric models.Metric) (models.Metric, error) {
	return s.src.Get(ctx, metric)
}

func (s *syncer) GetAll(ctx context.Context) ([]models.Metric, error) {
	return s.src.GetAll(ctx)
}

func (s *syncer) Shutdown(ctx context.Context) error {
	once.Do(func() {
		close(s.tick)
	})

	if err := s.dst.Write(); err != nil {
		return err
	}

	return s.dst.CloseFile()
}

func (s *syncer) SyncData() error {
	if s.i != 0 {
		go s.ticker()
	}

	for range s.tick {
		err := s.dst.Write()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *syncer) sync() {
	if s.i == 0 {
		s.tick <- struct{}{}
	}
}

func (s *syncer) ticker() {
	for range s.tick {
		time.Sleep(s.i)
		s.tick <- struct{}{}
	}
}
