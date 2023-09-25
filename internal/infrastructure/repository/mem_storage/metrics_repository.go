package memstorage

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

type MemStorage struct {
	Data map[string]models.Metric
	Mu   sync.RWMutex
}

func New(cfg *config.ServerConfig) *MemStorage {
	return &MemStorage{Data: make(map[string]models.Metric)}
}

func (ms *MemStorage) UpdateGauge(ctx context.Context, metric models.Metric) error {
	ms.Mu.Lock()
	defer ms.Mu.Unlock()

	m, ok := ms.Data[metric.ID]
	if !ok {
		ms.Data[metric.ID] = metric
	} else {
		m.Value = metric.Value
		ms.Data[metric.ID] = m
	}

	return nil
}

func (ms *MemStorage) UpdateCounter(ctx context.Context, metric models.Metric) error {
	ms.Mu.Lock()
	defer ms.Mu.Unlock()

	m, ok := ms.Data[metric.ID]
	if !ok {
		ms.Data[metric.ID] = metric
	} else {
		*m.Delta = *metric.Delta + *m.Delta
		ms.Data[metric.ID] = m
	}

	return nil
}

func (ms *MemStorage) Get(ctx context.Context, metric models.Metric) (models.Metric, error) {
	ms.Mu.RLock()
	defer ms.Mu.RUnlock()

	m, ok := ms.Data[metric.ID]
	if !ok {
		return models.Metric{ID: metric.ID, MType: "", Delta: nil, Value: nil}, fmt.Errorf("metric with ID %s %w", metric.ID, ErrNotFoundMetric)
	}

	return m, nil
}

func (ms *MemStorage) GetAll(ctx context.Context) ([]models.Metric, error) {
	var metrics []models.Metric

	for _, v := range ms.Data {
		ms.Mu.RLock()
		m := v
		ms.Mu.RUnlock()
		metrics = append(metrics, m)
	}

	return metrics, nil
}

func (ms *MemStorage) UpdateAll(ctx context.Context, metrics []models.Metric) error {
	for _, m := range metrics {
		ms.Mu.Lock()
		ms.Data[m.ID] = m
		ms.Mu.Unlock()
	}

	return nil
}
