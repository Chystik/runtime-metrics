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
	Mu   sync.Mutex
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
		m.Delta = metric.Delta
		ms.Data[metric.ID] = m
	}

	return nil
}

func (ms *MemStorage) Get(ctx context.Context, metric models.Metric) (models.Metric, error) {
	m, ok := ms.Data[metric.ID]
	if !ok {
		return models.Metric{ID: metric.ID, MType: "", Delta: nil, Value: nil}, fmt.Errorf("metric with ID %s %w", metric.ID, ErrNotFoundMetric)
	}

	return m, nil
}

func (ms *MemStorage) GetAll(ctx context.Context) ([]models.Metric, error) {
	var metrics []models.Metric

	for _, v := range ms.Data {
		m := v
		metrics = append(metrics, m)
	}

	return metrics, nil
}
