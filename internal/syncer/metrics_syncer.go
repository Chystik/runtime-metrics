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

func New(cfg *config.ServerConfig, src metricsservice.MetricsRepository, dst storage.MetricsStorage) (*syncer, error) {
	s := &syncer{}
	s.src = src
	s.dst = dst
	s.i = time.Duration(cfg.StoreInterval)
	s.tick = make(chan struct{}, 1)

	if cfg.Restore {
		err := s.dst.Read()
		if err != nil && err != io.EOF {
			return nil, err
		}
	}

	return s, nil
}

func (s *syncer) UpdateGauge(metric models.Metric) {
	s.src.UpdateGauge(metric)
	s.sync()
}
func (s *syncer) UpdateCounter(metric models.Metric) {
	s.src.UpdateCounter(metric)
	s.sync()
}

func (s *syncer) Get(metric models.Metric) (models.Metric, error) {
	return s.src.Get(metric)
}

func (s *syncer) GetAll() []models.Metric {
	return s.src.GetAll()
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
