package memstorage

import (
	"errors"
	"fmt"

	"github.com/Chystik/runtime-metrics/internal/models"
)

var ErrNotFoundMetric = errors.New("not found in repository")

type memStorage struct {
	data map[string]models.Metric
}

func New() *memStorage {
	return &memStorage{
		data: make(map[string]models.Metric),
	}
}

func (ms *memStorage) UpdateGauge(metric models.Metric) {
	m, ok := ms.data[metric.ID]
	if !ok {
		ms.data[metric.ID] = metric
	} else {
		m.Value = metric.Value
	}

}

func (ms *memStorage) UpdateCounter(metric models.Metric) {
	m, ok := ms.data[metric.ID]
	if !ok {
		ms.data[metric.ID] = metric
	} else {
		m.Delta = metric.Delta
	}
}

func (ms *memStorage) Get(ID string) (models.Metric, error) {
	m, ok := ms.data[ID]
	if !ok {
		return models.Metric{ID: ID, MType: "", Delta: nil, Value: nil}, fmt.Errorf("metric with ID %s %w", ID, ErrNotFoundMetric)
	}

	return m, nil
}

func (ms *memStorage) GetAll() []models.Metric {
	var metrics []models.Metric

	for _, v := range ms.data {
		m := v
		metrics = append(metrics, m)
	}

	return metrics
}
