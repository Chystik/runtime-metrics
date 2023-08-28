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
	m, ok := ms.data[metric.Id]
	if !ok {
		ms.data[metric.Id] = metric
	} else {
		m.Value = metric.Value
	}

}

func (ms *memStorage) UpdateCounter(metric models.Metric) {
	m, ok := ms.data[metric.Id]
	if !ok {
		ms.data[metric.Id] = metric
	} else {
		m.Delta = metric.Delta
	}
}

func (ms *memStorage) Get(Id string) (models.Metric, error) {
	m, ok := ms.data[Id]
	if !ok {
		return models.Metric{Id: Id, MType: "", Delta: nil, Value: nil}, fmt.Errorf("metric with Id %s %w", Id, ErrNotFoundMetric)
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
