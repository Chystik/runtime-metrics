package inmemory

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/models"
)

var (
	ErrNotFoundMetric = errors.New("not found in repository")
)

type memStorage struct {
	data map[string]models.Metric
	Mu   sync.RWMutex
}

func NewMetricsRepo(cfg *config.ServerConfig) *memStorage {
	return &memStorage{data: make(map[string]models.Metric)}
}

func (ms *memStorage) UpdateGauge(ctx context.Context, metric models.Metric) error {
	ms.Mu.Lock()
	defer ms.Mu.Unlock()

	m, ok := ms.data[metric.ID]
	if !ok {
		ms.data[metric.ID] = metric
	} else {
		m.Value = metric.Value
		ms.data[metric.ID] = m
	}

	return nil
}

func (ms *memStorage) UpdateCounter(ctx context.Context, metric models.Metric) error {
	ms.Mu.Lock()
	defer ms.Mu.Unlock()

	m, ok := ms.data[metric.ID]
	if !ok {
		ms.data[metric.ID] = metric
	} else {
		*m.Delta = *metric.Delta + *m.Delta
		ms.data[metric.ID] = m
	}

	return nil
}

func (ms *memStorage) Get(ctx context.Context, metric models.Metric) (models.Metric, error) {
	ms.Mu.RLock()
	defer ms.Mu.RUnlock()

	m, ok := ms.data[metric.ID]
	if !ok {
		return models.Metric{ID: metric.ID, MType: "", Delta: nil, Value: nil}, fmt.Errorf("metric with ID %s %w", metric.ID, ErrNotFoundMetric)
	}

	return m, nil
}

func (ms *memStorage) GetAll(ctx context.Context) ([]models.Metric, error) {
	var metrics []models.Metric

	for _, v := range ms.data {
		ms.Mu.RLock()
		m := v
		ms.Mu.RUnlock()
		metrics = append(metrics, m)
	}

	return metrics, nil
}

func (ms *memStorage) UpdateAll(ctx context.Context, metrics []models.Metric) error {
	for _, m := range metrics {
		ms.Mu.Lock()
		ms.data[m.ID] = m
		ms.Mu.Unlock()
	}

	return nil
}
