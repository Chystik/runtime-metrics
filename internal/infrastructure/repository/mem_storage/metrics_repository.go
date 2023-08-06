package memstorage

import (
	"fmt"

	"github.com/Chystik/runtime-metrics/internal/models"
	metricsservice "github.com/Chystik/runtime-metrics/internal/service/server"
)

type memStorage struct {
	data map[string]models.MetricValue
}

func New() metricsservice.MetricsRepository {
	return &memStorage{
		data: make(map[string]models.MetricValue),
	}
}

func (ms *memStorage) UpdateGauge(metric models.Metric) {
	var val models.MetricValue
	var ok bool

	if val, ok = ms.data[metric.Name]; ok {
		val.Gauge = metric.Gauge
	} else {
		val = metric.MetricValue
	}

	ms.data[metric.Name] = val
	//fmt.Printf("%s \t %#v\n", metric.Name, ms.data[metric.Name])
}

func (ms *memStorage) UpdateCounter(metric models.Metric) {
	var val models.MetricValue
	var ok bool

	if val, ok = ms.data[metric.Name]; ok {
		val.Counter = metric.Counter
	} else {
		val = metric.MetricValue
	}

	ms.data[metric.Name] = val
	//fmt.Printf("%s \t %#v\n", metric.Name, ms.data[metric.Name])
}

func (ms *memStorage) Get(name string) (models.Metric, error) {
	var metric models.Metric

	val, ok := ms.data[name]
	if !ok {
		return models.Metric{}, fmt.Errorf("not found metric with name %s", name)
	}

	metric.Name = name
	metric.MetricValue = val

	return metric, nil
}
